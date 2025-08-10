package software

import (
	"bufio"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/viper"
)

type Software struct {
	Name    string
	Version string
}

func isDebian() bool {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "ID=") {
			id := strings.TrimPrefix(line, "ID=")
			id = strings.Trim(id, "\"") // remove quotes if present
			return id == "debian"
		}
	}

	return false
}

func isInstalled(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func execute(cmd string) error {
	command := exec.Command("sh", "-c", cmd)
	out, err := command.CombinedOutput()
	if err != nil {
		return errors.New(err.Error() + " - " + string(out))
	}
	return nil
}

func supportedSoftware() map[string]func(version string) error {
	return map[string]func(version string) error{
		"apprise":      apprise,
		"borgbackup":   borgBackup,
		"docker":       docker,
		"git":          git,
		"podman":       podman,
		"rclone":       rclone,
		"rdiff-backup": rdiffBackup,
		"restic":       restic,
		"rsync":        rsync,
	}
}

func Install() {
	if runtime.GOOS != "linux" && !isDebian() {
		slog.Warn("OS not supported for software installation, skipping", "os", runtime.GOOS)
		return
	}

	var softwareList []Software
	err := viper.UnmarshalKey("software", &softwareList)
	if err != nil {
		slog.Error("Unable to decode software into struct", "err", err.Error())
		return
	}

	if len(softwareList) == 0 {
		slog.Debug("No software to install, skipping")
		return
	}

	updatePackages()
	for _, software := range softwareList {
		if get, ok := supportedSoftware()[software.Name]; ok {
			if isInstalled(software.Name) {
				continue
			}
			slog.Info("Installing software", "name", software.Name)
			err := get(software.Version)
			if err != nil {
				slog.Error("Failed", "err", err.Error())
				continue
			}
			slog.Info("Done")
		} else {
			slog.Error("Not supported, skipping", "name", software.Name)
		}
	}
	cleanup()
}

func updatePackages() {
	slog.Debug("Updating system packages")
	execute("apt-get update")
}

func cleanup() {
	// Clean up common documentation and cache directories to reduce image size
	slog.Debug("Cleaning up documentation and cache directories")
	execute("rm -rf /usr/share/doc /usr/share/man /usr/share/locale /var/cache/*")
}

func apprise(version string) error {
	install := "apprise"
	if version != "" {
		install = "apprise==" + version
	}
	return execute("pipx install " + install)
}

func borgBackup(version string) error {
	install := "borgbackup"
	if version != "" {
		install = "borgbackup=" + version
	}
	return execute("apt-get install -y " + install)
}

func docker(version string) error {
	if err := execute("install -m 0755 -d /etc/apt/keyrings"); err != nil {
		return err
	}
	if err := execute("curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc"); err != nil {
		return err
	}
	if err := execute("chmod a+r /etc/apt/keyrings/docker.asc"); err != nil {
		return err
	}
	if err := execute("echo deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable | tee /etc/apt/sources.list.d/docker.list > /dev/null"); err != nil {
		return err
	}
	updatePackages()
	if version != "" {
		return execute(fmt.Sprintf("apt-get install -y docker-ce-cli=%s docker-compose-plugin", version))
	}
	return execute("apt-get install -y docker-ce-cli docker-compose-plugin")
}

func git(version string) error {
	install := "git"
	if version != "" {
		install = "git=" + version
	}
	return execute("apt-get -y install " + install)
}

func podman(version string) error {
	install := "podman"
	if version != "" {
		install = "podman=" + version
	}
	return execute("apt-get -y install podman-compose " + install)
}

func rclone(version string) error {
	if version != "" {
		return execute("apt-get install -y " + version)
	}
	return execute("curl https://rclone.org/install.sh | bash")
}

func rdiffBackup(version string) error {
	install := "rdiff-backup"
	if version != "" {
		install = "rdiff-backup=" + version
	}
	return execute("apt-get -y install " + install)
}

func restic(version string) error {
	install := "restic"
	if version != "" {
		install = "restic=" + version
	}
	if err := execute("apt-get install -y " + install); err != nil {
		return err
	}
	return execute("restic self-update")
}

func rsync(version string) error {
	install := "rsync"
	if version != "" {
		install = "rsync=" + version
	}
	return execute("apt-get -y install " + install)
}
