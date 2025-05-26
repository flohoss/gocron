package commands

import (
	"os"
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
			Name:    "apprise",
			Version: extract(run("apprise --version"), `(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/caronc/apprise/releases",
		},
		Information{
			Name:    "borg",
			Version: extract(run("borg --version"), `(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/borgbackup/borg/releases",
		},
		Information{
			Name:    "curl",
			Version: extract(run("curl -V"), `(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/curl/curl/releases",
		},
		Information{
			Name:    "docker",
			Version: extract(run("docker version --format {{.Server.Version}}"), `(\d+\.\d+\.\d)`),
			Repo:    "https://docs.docker.com/engine/release-notes/",
		},
		Information{
			Name:    "docker-compose",
			Version: extract(run("docker compose version --short"), `(\d+\.\d+\.\d)`),
			Repo:    "https://docs.docker.com/compose/releases/release-notes/",
		},
		Information{
			Name:    "gocron",
			Version: extract(os.Getenv("APP_VERSION"), `(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/flohoss/gocron/releases",
		},
		Information{
			Name:    "podman",
			Version: extract(run("podman -v"), `(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/containers/podman/releases",
		},
		Information{
			Name:    "podman-compose",
			Version: extract(run("podman-compose -v"), `podman-compose.*(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/containers/podman/releases",
		},
		Information{
			Name:    "rclone",
			Version: extract(run("rclone version"), `(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/rclone/rclone/releases",
		},
		Information{
			Name:    "rdiff-backup",
			Version: extract(run("rdiff-backup -V"), `(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/rdiff-backup/rdiff-backup/releases",
		},
		Information{
			Name:    "restic",
			Version: extract(run("restic version"), `(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/restic/restic/releases",
		},
		Information{
			Name:    "rsync",
			Version: extract(run("rsync -V"), `(\d+\.\d+\.\d)`),
			Repo:    "https://github.com/RsyncProject/rsync/releases",
		},
		Information{
			Name:    "wget",
			Version: extract(run("wget -V"), `(\d+\.\d+\.\d)`),
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
	found := re.FindStringSubmatch(content)
	if len(found) >= 1 {
		return found[1]
	}
	return ""
}
