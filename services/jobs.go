package services

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
	"unicode"

	_ "github.com/glebarez/go-sqlite"

	"gitlab.unjx.de/flohoss/gobackup/config"
	"gitlab.unjx.de/flohoss/gobackup/internal/commands"
	"gitlab.unjx.de/flohoss/gobackup/internal/cron"
	"gitlab.unjx.de/flohoss/gobackup/services/jobs"
)

//go:embed jobs.sql
var ddl string

type TemplateJob struct {
	Job      jobs.Job
	Name     string
	Cron     string
	Commands []jobs.Command
	Envs     []jobs.Env
	Runs     []jobs.RunsView
}

func generateID(input string) string {
	var result strings.Builder

	// Iterate over each character in the input string
	for _, ch := range input {
		// Convert to lowercase and check if the character is alphanumeric or a space
		if unicode.IsLetter(ch) || unicode.IsDigit(ch) {
			result.WriteRune(unicode.ToLower(ch)) // Convert to lowercase and add it
		} else if ch == ' ' {
			result.WriteRune('_') // Replace space with underscore
		}
	}

	return result.String()
}

func NewJobService(dbName string, config *config.Config, cron *cron.Cron) (*JobService, error) {
	ctx := context.Background()

	db, err := sql.Open("sqlite", dbName+"?_pragma=foreign_keys(1)")
	if err != nil {
		return nil, err
	}

	if _, err := db.ExecContext(ctx, ddl); err != nil {
		return nil, err
	}

	queries := jobs.New(db)
	initEnums(queries, ctx)

	err = createUpdateOrDeleteJob(ctx, queries, config)
	if err != nil {
		return nil, err
	}

	err = createUpdateOrDeleteEnvs(ctx, queries, config)
	if err != nil {
		return nil, err
	}

	err = createUpdateOrDeleteCommands(ctx, queries, config)
	if err != nil {
		return nil, err
	}

	// no need for config any longer as all information is in db
	config = nil

	return &JobService{Queries: queries}, nil
}

type JobService struct {
	Queries *jobs.Queries
}

func initEnums(queries *jobs.Queries, ctx context.Context) {
	severities, _ := queries.ListSeverities(ctx)
	if len(severities) == 0 {
		queries.CreateSeverity(ctx, Debug.String())
		queries.CreateSeverity(ctx, Info.String())
		queries.CreateSeverity(ctx, Warning.String())
		queries.CreateSeverity(ctx, Error.String())
	}
	status, _ := queries.ListStatus(ctx)
	if len(status) == 0 {
		queries.CreateStatus(ctx, Running.String())
		queries.CreateStatus(ctx, Stopped.String())
		queries.CreateStatus(ctx, Finished.String())
	}
}

func createUpdateOrDeleteJob(ctx context.Context, queries *jobs.Queries, config *config.Config) error {
	dbJobs, _ := queries.ListJobs(ctx)

	existingJobs := make(map[string]bool)
	for _, j := range dbJobs {
		existingJobs[j.ID] = true
	}

	for _, j := range config.Jobs {
		jobID := generateID(j.Name)
		if _, exists := existingJobs[jobID]; exists {
			queries.UpdateJob(ctx, jobs.UpdateJobParams{
				ID:   jobID,
				Name: j.Name,
				Cron: j.Cron,
			})
		} else {
			queries.CreateJob(ctx, jobs.CreateJobParams{
				ID:   jobID,
				Name: j.Name,
				Cron: j.Cron,
			})
		}
		delete(existingJobs, jobID)
	}

	for id := range existingJobs {
		queries.DeleteJob(ctx, id)
	}

	return nil
}

func createUpdateOrDeleteEnvs(ctx context.Context, queries *jobs.Queries, config *config.Config) error {
	queries.DeleteEnvs(ctx)
	for _, job := range config.Jobs {
		for _, env := range job.Envs {
			maxRetries := 5
			for attempt := 1; attempt <= maxRetries; attempt++ {
				_, err := queries.CreateEnv(ctx, jobs.CreateEnvParams{
					JobID: generateID(job.Name), // Generate a new ID for each attempt
					Key:   env.Key,
					Value: env.Value,
				})

				if err == nil {
					break
				}

				if attempt == maxRetries {
					return fmt.Errorf("failed to insert environment variable after %d attempts: %v", maxRetries, err)
				}
				log.Printf("environment variable conflict detected, retrying (%d/%d)...", attempt, maxRetries)
			}
		}
	}
	return nil
}

func createUpdateOrDeleteCommands(ctx context.Context, queries *jobs.Queries, config *config.Config) error {
	queries.DeleteCommands(ctx)
	for _, job := range config.Jobs {
		for _, command := range job.Commands {
			maxRetries := 5
			for attempt := 1; attempt <= maxRetries; attempt++ {
				_, err := queries.CreateCommand(ctx, jobs.CreateCommandParams{
					JobID:   generateID(job.Name),
					Command: command.Command,
				})

				if err == nil {
					break
				}

				if attempt == maxRetries {
					return fmt.Errorf("failed to insert command after %d attempts: %v", maxRetries, err)
				}
				log.Printf("command conflict detected, retrying (%d/%d)...", attempt, maxRetries)
			}
		}
	}
	return nil
}

func (js *JobService) GetQueries() *jobs.Queries {
	return js.Queries
}

func (js *JobService) ExecuteJobs() {
	jobs, _ := js.Queries.ListJobs(context.Background())
	for i := 0; i < len(jobs); i++ {
		js.ExecuteJob(&jobs[i])
	}
}

func (js *JobService) ExecuteJob(job *jobs.Job) {
	ctx := context.Background()

	run, _ := js.Queries.CreateRun(ctx, jobs.CreateRunParams{
		JobID:    job.ID,
		StatusID: int64(Running),
	})
	status := Finished

	envs, _ := js.Queries.ListEnvsByJobID(ctx, job.ID)
	for _, command := range envs {
		js.Queries.CreateLog(ctx, jobs.CreateLogParams{
			RunID:      run.ID,
			SeverityID: int64(Debug),
			Message:    fmt.Sprintf("Setting environment variable: \"%s\"", command.Key),
		})
		os.Setenv(command.Key, command.Value)
	}

	cmds, _ := js.Queries.ListCommandsByJobID(ctx, job.ID)
	for _, command := range cmds {
		program, args, err := commands.PrepareCommand(command.Command)
		if err != nil {
		}
		js.Queries.CreateLog(ctx, jobs.CreateLogParams{
			RunID:      run.ID,
			SeverityID: int64(Debug),
			Message:    fmt.Sprintf("Executing command: \"%s\" - \"%s\"", program, args),
		})
		out, err := commands.ExecuteCommand(program, args)
		severity := Info
		if err != nil {
			severity = Error
		}
		js.Queries.CreateLog(ctx, jobs.CreateLogParams{
			RunID:      run.ID,
			SeverityID: int64(severity),
			Message:    out,
		})
		if err != nil {
			status = Stopped
			break
		}
		// Wait for 2 seconds to make sure the command is finished
		time.Sleep(2 * time.Second)
	}

	js.Queries.UpdateRun(ctx, jobs.UpdateRunParams{
		StatusID: int64(status),
		EndTime:  sql.NullTime{Time: time.Now().UTC(), Valid: true},
		ID:       run.ID,
	})
}
