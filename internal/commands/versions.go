package commands

import (
	"database/sql"
	"regexp"
)

type Versions []Information

type Information struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Repo    string `json:"repo"`
}

var v *Versions

func init() {
	v = &Versions{
		Information{
			Name:    "restic",
			Version: extract(run("restic", []string{"version"}), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/restic/restic/releases",
		},
		Information{
			Name:    "borg",
			Version: extract(run("borg", []string{"--version"}), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/borgbackup/borg/releases",
		},
		Information{
			Name:    "rclone",
			Version: extract(run("rclone", []string{"version"}), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/rclone/rclone/releases",
		},
		Information{
			Name:    "curl",
			Version: extract(run("curl", []string{"-V"}), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/curl/curl/releases",
		},
		Information{
			Name:    "rsync",
			Version: extract(run("rsync", []string{"-V"}), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/RsyncProject/rsync/releases",
		},
		Information{
			Name:    "rdiff-backup",
			Version: extract(run("rdiff-backup", []string{"-V"}), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/rdiff-backup/rdiff-backup/releases",
		},
		Information{
			Name:    "docker",
			Version: run("docker", []string{"version", "--format", "{{.Server.Version}}"}),
			Repo:    "https://docs.docker.com/engine/release-notes/",
		},
		Information{
			Name:    "compose",
			Version: run("docker", []string{"compose", "version", "--short"}),
			Repo:    "https://docs.docker.com/compose/releases/release-notes/",
		},
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
