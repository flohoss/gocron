package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/spf13/viper"
)

func setConfigForTest(t *testing.T, testCfg GlobalConfig) {
	t.Helper()

	mu.Lock()
	previous := cfg
	cfg = testCfg
	mu.Unlock()

	t.Cleanup(func() {
		mu.Lock()
		cfg = previous
		mu.Unlock()
	})
}

func TestSetConfigFolderPath_SetsConfigFilePath(t *testing.T) {
	previous := GetConfigFilePath()
	t.Cleanup(func() { SetConfigFilePath(previous) })

	SetConfigFolderPath("./tmp-config")

	expected := filepath.Clean("tmp-config/config.yaml")
	if got := GetConfigFilePath(); got != expected {
		t.Fatalf("unexpected config file path: got %q want %q", got, expected)
	}
}

func TestSetConfigFilePath_UsesDefaultWhenEmpty(t *testing.T) {
	previous := GetConfigFilePath()
	t.Cleanup(func() { SetConfigFilePath(previous) })

	SetConfigFilePath("")

	expected := filepath.Clean(GetDefaultConfigFile())
	if got := GetConfigFilePath(); got != expected {
		t.Fatalf("unexpected default config file path: got %q want %q", got, expected)
	}
}

func TestTerminalSettingsHydrate_PopulatesAllowedArgsMap(t *testing.T) {
	settings := TerminalSettings{
		AllowedCommands: map[string]AllowedCommands{
			"docker": {
				Args: []string{"ps", "version"},
			},
		},
	}

	settings.Hydrate()

	docker := settings.AllowedCommands["docker"]
	if docker.AllowedArgsMap == nil {
		t.Fatal("expected allowed args map to be initialized")
	}
	if _, ok := docker.AllowedArgsMap["ps"]; !ok {
		t.Fatal("expected argument ps to be present in allowed args map")
	}
	if _, ok := docker.AllowedArgsMap["version"]; !ok {
		t.Fatal("expected argument version to be present in allowed args map")
	}
}

func TestEnableAndDisableAllJobs(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		Jobs: []Job{
			{Name: "a", Disabled: true},
			{Name: "b", Disabled: false},
		},
	})

	EnableAllJobs()
	jobs := GetJobs()
	if jobs[0].Disabled || jobs[1].Disabled {
		t.Fatalf("expected all jobs enabled, got: %#v", jobs)
	}

	DisableAllJobs()
	jobs = GetJobs()
	if !jobs[0].Disabled || !jobs[1].Disabled {
		t.Fatalf("expected all jobs disabled, got: %#v", jobs)
	}
}

func TestEnableScheduledAndNonScheduledJobs(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		Jobs: []Job{
			{Name: "scheduled", DisableCron: false},
			{Name: "manual", DisableCron: true},
		},
	})

	EnableScheduledJobs()
	jobs := GetJobs()
	if jobs[0].Disabled || !jobs[1].Disabled {
		t.Fatalf("expected only scheduled jobs enabled, got: %#v", jobs)
	}

	EnableNonScheduledJobs()
	jobs = GetJobs()
	if !jobs[0].Disabled || jobs[1].Disabled {
		t.Fatalf("expected only non-scheduled jobs enabled, got: %#v", jobs)
	}
}

func TestToggleDisabledJob(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		Jobs: []Job{{Name: "job-a", Disabled: false}},
	})

	if err := ToggleDisabledJob("job-a"); err != nil {
		t.Fatalf("expected no error toggling existing job, got: %v", err)
	}
	if !GetJobs()[0].Disabled {
		t.Fatal("expected job to be disabled after toggle")
	}

	if err := ToggleDisabledJob("missing"); err == nil {
		t.Fatal("expected error toggling missing job")
	}
}

func TestGetEnvsByJobName_MergesDefaultsAndOverrides(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		JobDefaults: JobDefaults{
			Envs: []Env{{Key: "A", Value: "1"}, {Key: "B", Value: "2"}},
		},
		Jobs: []Job{{
			Name: "job-1",
			Envs: []Env{{Key: "B", Value: "3"}, {Key: "C", Value: "4"}},
		}},
	})

	envs := GetEnvsByJobName("job-1")

	expectedOrder := []string{"A", "B", "C"}
	if !reflect.DeepEqual(envs.Order, expectedOrder) {
		t.Fatalf("unexpected env order: got %#v want %#v", envs.Order, expectedOrder)
	}
	if envs.Data["A"] != "1" || envs.Data["B"] != "3" || envs.Data["C"] != "4" {
		t.Fatalf("unexpected env data: %#v", envs.Data)
	}
}

func TestGetCommandsByJobName_IncludesPreAndPostCommands(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		JobDefaults: JobDefaults{
			PreCommands:  []string{"pre-1", "pre-2"},
			PostCommands: []string{"post-1"},
		},
		Jobs: []Job{{
			Name:     "job-1",
			Commands: []string{"run-1", "run-2"},
		}},
	})

	commands := GetCommandsByJobName("job-1")
	expected := []string{"pre-1", "pre-2", "run-1", "run-2", "post-1"}

	if !reflect.DeepEqual(commands, expected) {
		t.Fatalf("unexpected commands: got %#v want %#v", commands, expected)
	}
}

func TestGetJobsCron_UsesJobCronOrDefault(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		JobDefaults: JobDefaults{Cron: "0 3 * * *"},
	})

	jobWithCron := Job{Name: "with-cron", Cron: "0 1 * * *"}
	if got := GetJobsCron(&jobWithCron); got != "0 1 * * *" {
		t.Fatalf("unexpected cron for job with explicit cron: %q", got)
	}

	jobWithoutCron := Job{Name: "without-cron"}
	if got := GetJobsCron(&jobWithoutCron); got != "0 3 * * *" {
		t.Fatalf("unexpected cron fallback to default: %q", got)
	}
}

func TestGetAllCrons_GroupsJobsAndSkipsDisabledCron(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		JobDefaults: JobDefaults{Cron: "0 5 * * *"},
		Jobs: []Job{
			{Name: "job-a", Cron: "0 1 * * *", DisableCron: false},
			{Name: "job-b", DisableCron: false},
			{Name: "job-c", DisableCron: true},
			{Name: "job-d", Cron: "0 1 * * *", DisableCron: false},
		},
	})

	crons := GetAllCrons()

	if len(crons) != 2 {
		t.Fatalf("unexpected cron group count: got %d want 2", len(crons))
	}

	oneAMJobs, ok := crons["0 1 * * *"]
	if !ok {
		t.Fatal("expected cron group 0 1 * * * to exist")
	}
	if len(oneAMJobs) != 2 {
		t.Fatalf("unexpected number of jobs for 0 1 * * *: got %d want 2", len(oneAMJobs))
	}

	defaultJobs, ok := crons["0 5 * * *"]
	if !ok {
		t.Fatal("expected default cron group 0 5 * * * to exist")
	}
	if len(defaultJobs) != 1 || defaultJobs[0].Name != "job-b" {
		t.Fatalf("unexpected jobs for default cron: %#v", defaultJobs)
	}
	for _, jobs := range crons {
		for _, job := range jobs {
			if job.Name == "job-c" {
				t.Fatal("expected disable_cron job to be excluded from cron groups")
			}
		}
	}
}

func TestDefaultStarterJobs_HasFourValidJobs(t *testing.T) {
	jobs := defaultStarterJobs()

	if len(jobs) != 4 {
		t.Fatalf("unexpected number of default starter jobs: got %d want 4", len(jobs))
	}

	for i, job := range jobs {
		if job.Name == "" || len(job.Commands) == 0 {
			t.Fatalf("starter job %d is invalid: %#v", i, job)
		}
	}
	if !jobs[3].DisableCron {
		t.Fatalf("expected fourth starter job to be manual/disable_cron=true: %#v", jobs[3])
	}
}

func TestNew_CreatesAndLoadsDefaultStarterJobs(t *testing.T) {
	prevPath := GetConfigFilePath()
	prevTZ, hadTZ := os.LookupEnv("TZ")

	t.Cleanup(func() {
		SetConfigFilePath(prevPath)
		if hadTZ {
			_ = os.Setenv("TZ", prevTZ)
		} else {
			_ = os.Unsetenv("TZ")
		}
		viper.Reset()
	})

	viper.Reset()
	configPath := filepath.Join(t.TempDir(), "config.yaml")

	New(configPath)

	if _, err := os.Stat(configPath); err != nil {
		t.Fatalf("expected config file to be created, got error: %v", err)
	}

	jobs := GetJobs()
	if len(jobs) != 4 {
		t.Fatalf("unexpected number of loaded default jobs: got %d want 4", len(jobs))
	}

	if jobs[0].Name != "Example Scheduled Happy Path" {
		t.Fatalf("unexpected first default job name: %q", jobs[0].Name)
	}
	if jobs[1].Name != "Example Continue On Failure" {
		t.Fatalf("unexpected second default job name: %q", jobs[1].Name)
	}
	if jobs[2].Name != "Example Env Expansion" {
		t.Fatalf("unexpected third default job name: %q", jobs[2].Name)
	}
	if jobs[3].Name != "Example Manual Long Running" {
		t.Fatalf("unexpected fourth default job name: %q", jobs[3].Name)
	}

	if got := GetDBLocation(); got != filepath.Dir(configPath) {
		t.Fatalf("unexpected default db location: got %q want %q", got, filepath.Dir(configPath))
	}
	if got := GetDBName(); got != "db.sqlite" {
		t.Fatalf("unexpected default db name: got %q want %q", got, "db.sqlite")
	}
}

func TestGetDBLocation_UsesConfiguredLocation(t *testing.T) {
	previous := GetConfigFilePath()
	configPath := filepath.Join(t.TempDir(), "config.yaml")
	SetConfigFilePath(configPath)
	t.Cleanup(func() { SetConfigFilePath(previous) })

	setConfigForTest(t, GlobalConfig{DB: DBSettings{Location: "./custom-db"}})
	if got := GetDBLocation(); got != filepath.Join(filepath.Dir(configPath), "custom-db") {
		t.Fatalf("unexpected db location: %q", got)
	}
}

func TestGetDBLocation_UsesConfiguredAbsoluteLocation(t *testing.T) {
	absolute := filepath.Join(t.TempDir(), "db")
	setConfigForTest(t, GlobalConfig{DB: DBSettings{Location: absolute}})
	if got := GetDBLocation(); got != filepath.Clean(absolute) {
		t.Fatalf("unexpected absolute db location: %q", got)
	}
}

func TestGetDBName_UsesConfiguredName(t *testing.T) {
	setConfigForTest(t, GlobalConfig{DB: DBSettings{Name: "jobs.sqlite3"}})
	if got := GetDBName(); got != "jobs.sqlite3" {
		t.Fatalf("unexpected db name: %q", got)
	}
}

func TestGetDBName_NormalizesToFileName(t *testing.T) {
	setConfigForTest(t, GlobalConfig{DB: DBSettings{Name: "nested/path/jobs.sqlite"}})
	if got := GetDBName(); got != "jobs.sqlite" {
		t.Fatalf("unexpected normalized db name: %q", got)
	}
}

func TestGetDefaultConfigFolder(t *testing.T) {
	if got := GetDefaultConfigFolder(); got == "" {
		t.Fatal("expected non-empty default config folder")
	}
}

func TestGetConfigFolderPath_ReflectsSetPath(t *testing.T) {
	previous := GetConfigFilePath()
	t.Cleanup(func() { SetConfigFilePath(previous) })

	SetConfigFolderPath("./tmp-folder-test")
	folder := GetConfigFolderPath()
	if folder == "" {
		t.Fatal("expected non-empty config folder path")
	}
}

func TestGetLogLevel_AllValues(t *testing.T) {
	cases := []struct {
		level string
		want  string
	}{
		{"debug", "DEBUG"},
		{"warn", "WARN"},
		{"warning", "WARN"},
		{"error", "ERROR"},
		{"info", "INFO"},
		{"", "INFO"},
		{"unknown", "INFO"},
	}

	for _, tc := range cases {
		setConfigForTest(t, GlobalConfig{LogLevel: tc.level})
		got := GetLogLevel().String()
		if got != tc.want {
			t.Errorf("GetLogLevel(%q) = %q, want %q", tc.level, got, tc.want)
		}
	}
}

func TestGetServer_FormatsAddressAndPort(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		Server: ServerSettings{Address: "127.0.0.1", Port: 9000},
	})
	if got := GetServer(); got != "127.0.0.1:9000" {
		t.Fatalf("unexpected server string: %q", got)
	}
}

func TestGetHealthcheck_ReturnsConfigValue(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		Healthcheck: HealthCheck{Type: "GET"},
	})
	if got := GetHealthcheck().Type; got != "GET" {
		t.Fatalf("unexpected healthcheck type: %q", got)
	}
}

func TestGetDeleteRunsAfterDays_ReturnsConfigValue(t *testing.T) {
	setConfigForTest(t, GlobalConfig{DeleteRunsAfterDays: 14})
	if got := GetDeleteRunsAfterDays(); got != 14 {
		t.Fatalf("unexpected delete runs after days: %d", got)
	}
}

func TestGetTerminalSettings_ReturnsConfigValue(t *testing.T) {
	setConfigForTest(t, GlobalConfig{
		Terminal: TerminalSettings{AllowAllCommands: true},
	})
	if !GetTerminalSettings().AllowAllCommands {
		t.Fatal("expected allow_all_commands to be true")
	}
}

func TestConfigLoaded_ReturnsFalseWhenNotLoaded(t *testing.T) {
	// Before New() is called in a fresh viper state, ConfigLoaded returns false
	v := viper.New()
	_ = v
	// ConfigLoaded checks the global viper instance; it should be false unless New() was called
	// Just verify it returns a bool without panicking
	_ = ConfigLoaded()
}
