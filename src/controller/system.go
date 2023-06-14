package controller

import (
	"net/http"
	"os"
	"regexp"
	"runtime"

	"github.com/labstack/echo/v4"
)

type SystemData struct {
	Title         string
	Versions      Versions
	Configuration Configuration
}

type Versions struct {
	Go      string
	Rclone  string
	Restic  string
	Docker  string
	Compose string
}

type Configuration struct {
	Hostname         string
	RcloneConfigFile string
}

func (c *Controller) findVersion(name string, vRegexStr string, command ...string) string {
	out, _ := c.executeCmd(name, command...)
	vRegex := regexp.MustCompile(vRegexStr)
	groups := vRegex.FindSubmatch(out)
	if len(groups) < 1 {
		return ""
	}
	return "v" + string(groups[1])
}

func (c *Controller) RenderSystem(ctx echo.Context) error {
	return ctx.Render(http.StatusOK, "system", SystemData{Title: c.env.Identifier + " - System", Versions: c.Versions, Configuration: c.Configuration})
}

func (c *Controller) setVersions() {
	c.Versions = Versions{
		Go:      runtime.Version(),
		Rclone:  c.findVersion("rclone", `((?:\d{1,2}.){2}\d{1,2})`, "version"),
		Restic:  c.findVersion("restic", `((?:\d{1,2}.){2}\d{1,2})`, "version"),
		Docker:  c.findVersion("docker", `Engine:\s*Version:\s*((?:\d{1,2}.){2}\d{1,2})`, "version"),
		Compose: c.findVersion("docker", `((?:\d{1,2}.){2}\d{1,2})`, "compose", "version"),
	}
	hostname, _ := os.Hostname()
	c.Configuration = Configuration{
		Hostname:         hostname,
		RcloneConfigFile: rcloneConfigFile(),
	}
}
