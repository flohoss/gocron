package buildinfo

import "fmt"

var (
	Version   = "dev"
	BuildTime = "unknown"
	RepoURL   = ""
)

func Summary() string {
	if RepoURL == "" {
		return fmt.Sprintf("gocron %s (built %s)", Version, BuildTime)
	}

	return fmt.Sprintf("gocron %s (built %s, source %s)", Version, BuildTime, RepoURL)
}
