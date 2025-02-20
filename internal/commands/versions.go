package commands

import (
	"database/sql"
	"regexp"
)

type Versions struct {
	Body struct {
		Restic  Information `json:"restic"`
		Borg    Information `json:"borg"`
		Rclone  Information `json:"rclone"`
		Docker  Information `json:"docker"`
		Compose Information `json:"compose"`
	}
}

type Information struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Repo    string `json:"repo"`
}

var v *Versions

func init() {
	v = &Versions{}
	v.Body.Restic = Information{
		Name:    "restic",
		Version: extract(run("restic", []string{"version"}), `\d+\.\d+\.\d`),
		Repo:    "https://github.com/restic/restic/releases/tag/v",
	}
	v.Body.Borg = Information{
		Name:    "borg",
		Version: extract(run("borg", []string{"--version"}), `\d+\.\d+\.\d`),
		Repo:    "https://github.com/borgbackup/borg/releases/tag/",
	}
	v.Body.Rclone = Information{
		Name:    "rclone",
		Version: extract(run("rclone", []string{"version"}), `v\d+\.\d+\.\d`),
		Repo:    "https://github.com/rclone/rclone/releases/tag/",
	}
	v.Body.Docker = Information{
		Name:    "docker",
		Version: run("docker", []string{"version", "--format", "{{.Server.Version}}"}),
		Repo:    "https://docs.docker.com/engine/release-notes/",
	}
	v.Body.Compose = Information{
		Name:    "compose",
		Version: run("docker", []string{"compose", "version", "--short"}),
		Repo:    "https://docs.docker.com/compose/releases/release-notes/",
	}
}

func GetVersions() *Versions {
	return v
}

func run(program string, args []string) string {
	res, _ := ExecuteCommand(program, args, sql.NullString{Valid: false})
	return res
}

func extract(content string, regex string) string {
	re := regexp.MustCompile(regex)
	return re.FindString(content)
}
