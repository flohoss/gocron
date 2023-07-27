package system

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"gitlab.unjx.de/flohoss/gobackup/database"
)

type Data struct {
	Versions      Versions          `json:"versions" validate:"required"`
	Configuration Configuration     `json:"configuration" validate:"required"`
	Stats         database.JobStats `json:"job_stats" validate:"required"`
}

type Versions struct {
	Go       string `json:"go" validate:"required"`
	GoBackup string `json:"gobackup" validate:"required"`
	Rclone   string `json:"rclone" validate:"required"`
	Restic   string `json:"restic" validate:"required"`
	Docker   string `json:"docker" validate:"required"`
	Compose  string `json:"compose" validate:"required"`
}

type Configuration struct {
	Hostname         string `json:"hostname" validate:"required"`
	RcloneConfigFile string `json:"rclone_config_file" validate:"required"`
}

var SystemData Data

func init() {
	hostname, _ := os.Hostname()
	SystemData = Data{
		Versions: Versions{
			Go:       strings.Replace(runtime.Version(), "go", "v", 1),
			GoBackup: os.Getenv("APP_VERSION"),
			Rclone:   FindVersion("rclone", `((?:\d{1,2}.){2}\d{1,2})`, "version"),
			Restic:   FindVersion("restic", `((?:\d{1,2}.){2}\d{1,2})`, "version"),
			Docker:   FindVersion("docker", `Engine:\s*Version:\s*((?:\d{1,2}.){2}\d{1,2})`, "version"),
			Compose:  FindVersion("docker", `((?:\d{1,2}.){2}\d{1,2})`, "compose", "version"),
		},
		Configuration: Configuration{
			Hostname:         hostname,
			RcloneConfigFile: RcloneConfigFile(),
		},
	}
}

func FindVersion(name string, vRegexStr string, commands ...string) string {
	out, err := exec.Command(name, commands...).CombinedOutput()
	if err != nil {
		return ""
	}
	vRegex := regexp.MustCompile(vRegexStr)
	groups := vRegex.FindSubmatch(out)
	if len(groups) < 1 {
		return ""
	}
	return "v" + string(groups[1])
}

func RcloneConfigFile() string {
	filePath := "~/.config/rclone/rclone.conf"
	if _, err := os.Stat(filePath); err != nil {
		return ""
	}
	return filePath
}
