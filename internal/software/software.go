package software

import (
	"errors"
	"log/slog"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/flohoss/gocron/config"
)

type Software = config.Software

func canInstallOnOS(osReleasePath string) bool {
	if runtime.GOOS != "linux" {
		return false
	}
	return osReleaseIsDebian(osReleasePath)
}

func osReleaseIsDebian(path string) bool {
	data, err := os.ReadFile(path)
	if err != nil {
		return false
	}

	for line := range strings.SplitSeq(string(data), "\n") {
		line = strings.TrimSpace(line)

		if after, ok := strings.CutPrefix(line, "ID="); ok {
			id := strings.Trim(after, `"`)
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

type runner func(cmd string) error

type installer struct {
	run       runner
	installed func(name string) bool
}

func (i installer) shell(cmd string) error {
	return i.run(cmd)
}

func Install() {
	if !canInstallOnOS("/etc/os-release") {
		slog.Warn("Software installation is only supported on Debian Linux, skipping", "os", runtime.GOOS)
		return
	}

	softwareList := config.GetSoftware()
	if len(softwareList) == 0 {
		slog.Debug("No software to install, skipping")
		return
	}

	installer{run: execute, installed: isInstalled}.install(softwareList)
}

func (i installer) updatePackages() {
	slog.Debug("Updating system packages")
	i.shell("apt-get update")
}

func (i installer) cleanup() {
	// Clean up common documentation and cache directories to reduce image size
	slog.Debug("Cleaning up documentation and cache directories")
	i.shell("rm -rf /usr/share/doc /usr/share/man /usr/share/locale /var/cache/*")
}

func (i installer) supported() map[string]func(version string) error {
	return map[string]func(version string) error{
		"apprise":      i.apprise,
		"borgbackup":   i.borgBackup,
		"docker":       i.docker,
		"git":          i.git,
		"podman":       i.podman,
		"rclone":       i.rclone,
		"rdiff-backup": i.rdiffBackup,
		"restic":       i.restic,
		"rsync":        i.rsync,
		"logrotate":    i.logrotate,
		"sqlite3":      i.sqlite3,
		"kopia":        i.kopia,
	}
}

func (i installer) install(list []Software) {
	i.updatePackages()
	supported := i.supported()
	for _, software := range list {
		get, ok := supported[software.Name]
		if !ok {
			slog.Error("Not supported, skipping", "name", software.Name)
			continue
		}
		if i.installed(software.Name) {
			continue
		}
		slog.Info("Installing software", "name", software.Name)
		err := get(software.Version)
		if err != nil {
			slog.Error("Failed", "err", err.Error())
			continue
		}
		slog.Info("Done")
	}
	i.cleanup()
}

func (i installer) apprise(version string) error {
	pkg := "apprise"
	if version != "" {
		pkg = "apprise==" + version
	}
	return i.shell("pipx install " + pkg)
}

func (i installer) borgBackup(version string) error {
	pkg := "borgbackup"
	if version != "" {
		pkg = "borgbackup=" + version
	}
	return i.shell("apt-get install -y " + pkg)
}

func (i installer) docker(version string) error {
	if err := i.shell("install -m 0755 -d /etc/apt/keyrings"); err != nil {
		return err
	}
	if err := i.shell("curl -fsSL https://download.docker.com/linux/debian/gpg -o /etc/apt/keyrings/docker.asc"); err != nil {
		return err
	}
	if err := i.shell("chmod a+r /etc/apt/keyrings/docker.asc"); err != nil {
		return err
	}
	if err := i.shell("echo deb [arch=$(dpkg --print-architecture) signed-by=/etc/apt/keyrings/docker.asc] https://download.docker.com/linux/debian $(. /etc/os-release && echo \"$VERSION_CODENAME\") stable | tee /etc/apt/sources.list.d/docker.list > /dev/null"); err != nil {
		return err
	}
	i.updatePackages()
	if version != "" {
		return i.shell("apt-get install -y docker-ce-cli=" + version + " docker-compose-plugin")
	}
	return i.shell("apt-get install -y docker-ce-cli docker-compose-plugin")
}

func (i installer) git(version string) error {
	pkg := "git"
	if version != "" {
		pkg = "git=" + version
	}
	return i.shell("apt-get -y install " + pkg)
}

func (i installer) podman(version string) error {
	pkg := "podman"
	if version != "" {
		pkg = "podman=" + version
	}
	return i.shell("apt-get -y install podman-compose " + pkg)
}

func (i installer) rclone(version string) error {
	if version == "" {
		return i.shell("curl https://rclone.org/install.sh | bash")
	}
	return i.shell("apt-get install -y rclone=" + version)
}

func (i installer) rdiffBackup(version string) error {
	pkg := "rdiff-backup"
	if version != "" {
		pkg = "rdiff-backup=" + version
	}
	return i.shell("apt-get -y install " + pkg)
}

func (i installer) restic(version string) error {
	pkg := "restic"
	if version != "" {
		pkg = "restic=" + version
	}
	if err := i.shell("apt-get install -y " + pkg); err != nil {
		return err
	}
	return i.shell("restic self-update")
}

func (i installer) rsync(version string) error {
	pkg := "rsync"
	if version != "" {
		pkg = "rsync=" + version
	}
	return i.shell("apt-get -y install " + pkg)
}

func (i installer) logrotate(version string) error {
	pkg := "logrotate"
	if version != "" {
		pkg = "logrotate=" + version
	}
	return i.shell("apt-get -y install " + pkg)
}

func (i installer) sqlite3(version string) error {
	pkg := "sqlite3"
	if version != "" {
		pkg = "sqlite3=" + version
	}
	return i.shell("apt-get -y install " + pkg)
}

func (i installer) kopia(version string) error {
	if err := i.shell("install -m 0755 -d /etc/apt/keyrings"); err != nil {
		return err
	}
	if err := i.shell("curl -s https://kopia.io/signing-key | gpg --dearmor -o /etc/apt/keyrings/kopia-keyring.gpg"); err != nil {
		return err
	}
	if err := i.shell("chmod a+r /etc/apt/keyrings/kopia-keyring.gpg"); err != nil {
		return err
	}

	if err := i.shell("echo deb [signed-by=/etc/apt/keyrings/kopia-keyring.gpg] http://packages.kopia.io/apt/ stable main | tee /etc/apt/sources.list.d/kopia.list > /dev/null"); err != nil {
		return err
	}
	i.updatePackages()
	if version != "" {
		return i.shell("apt-get install -y kopia=" + version)
	}
	return i.shell("apt-get install -y kopia")
}
