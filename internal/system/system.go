package system

import (
	"os"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

type Data struct {
	Versions      Versions      `json:"versions"`
	Configuration Configuration `json:"configuration"`
	Disk          Disk          `json:"disk"`
}

type Versions struct {
	Go       string `json:"go"`
	GoBackup string `json:"gobackup"`
	Rclone   string `json:"rclone"`
	Restic   string `json:"restic"`
	Docker   string `json:"docker"`
	Compose  string `json:"compose"`
}

type Configuration struct {
	Hostname         string `json:"hostname"`
	RcloneConfigFile string `json:"rclone_config_file"`
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
	filePath := "/root/.config/rclone/rclone.conf"
	if _, err := os.Stat(filePath); err != nil {
		return ""
	}
	return filePath
}
