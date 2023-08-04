package system

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"

	"gitlab.unjx.de/flohoss/gobackup/internal/env"
)

type Versions struct {
	Go       string `json:"go" validate:"required"`
	GoBackup string `json:"gobackup" validate:"required"`
	Rclone   string `json:"rclone" validate:"required"`
	Restic   string `json:"restic" validate:"required"`
	Docker   string `json:"docker" validate:"required"`
	Compose  string `json:"compose" validate:"required"`
}

type SystemConfig struct {
	Config           env.Config `json:"config" validate:"required"`
	RcloneConfigFile string     `json:"rclone_config_file" validate:"required"`
	Hostname         string     `json:"hostname" validate:"required"`
	Versions         Versions   `json:"versions" validate:"required"`
}

var SystemConf SystemConfig

func init() {
	hostname, _ := os.Hostname()
	SystemConf = SystemConfig{
		Hostname:         hostname,
		RcloneConfigFile: RcloneConfigFile(),
		Versions: Versions{
			Go:       strings.Replace(runtime.Version(), "go", "v", 1),
			GoBackup: os.Getenv("APP_VERSION"),
			Rclone:   FindVersion("rclone", `((?:\d{1,2}.){2}\d{1,2})`, "version"),
			Restic:   FindVersion("restic", `((?:\d{1,2}.){2}\d{1,2})`, "version"),
			Docker:   FindVersion("docker", `Engine:\s*Version:\s*((?:\d{1,2}.){2}\d{1,2})`, "version"),
			Compose:  FindVersion("docker", `((?:\d{1,2}.){2}\d{1,2})`, "compose", "version"),
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
	filePath := "/root/.config/rclone/rclone.conf"
	if _, err := os.Stat(filePath); err != nil {
		return ""
	}
	return filePath
}
