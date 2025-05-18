package commands

import (
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
			Version: extract(run("restic version"), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/restic/restic/releases",
		},
		Information{
			Name:    "borg",
			Version: extract(run("borg --version"), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/borgbackup/borg/releases",
		},
		Information{
			Name:    "rclone",
			Version: extract(run("rclone version"), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/rclone/rclone/releases",
		},
		Information{
			Name:    "rsync",
			Version: extract(run("rsync -V"), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/RsyncProject/rsync/releases",
		},
		Information{
			Name:    "rdiff-backup",
			Version: extract(run("rdiff-backup -V"), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/rdiff-backup/rdiff-backup/releases",
		},
		Information{
			Name:    "docker",
			Version: run("docker version --format {{.Server.Version}}"),
			Repo:    "https://docs.docker.com/engine/release-notes/",
		},
		Information{
			Name:    "compose",
			Version: run("docker compose version --short"),
			Repo:    "https://docs.docker.com/compose/releases/release-notes/",
		},
		Information{
			Name:    "apprise",
			Version: extract(run("apprise --version"), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/caronc/apprise/releases",
		},
		Information{
			Name:    "curl",
			Version: extract(run("curl -V"), `\d+\.\d+\.\d`),
			Repo:    "https://github.com/curl/curl/releases",
		},
		Information{
			Name:    "wget",
			Version: extract(run("wget -V"), `\d+\.\d+\.\d`),
			Repo:    "https://ftp.gnu.org/gnu/wget/",
		},
	}
}

func GetVersions() *Versions {
	return v
}

func run(cmdString string) string {
	res, _ := ExecuteCommand(cmdString)
	return res
}

func extract(content string, regex string) string {
	re := regexp.MustCompile(regex)
	return re.FindString(content)
}
