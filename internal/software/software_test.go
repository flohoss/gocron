package software

import (
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"sync"
	"testing"
)

type recordingRunner struct {
	mu       sync.Mutex
	commands []string
	failOn   func(cmd string) error
}

func (r *recordingRunner) run(cmd string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.commands = append(r.commands, cmd)
	if r.failOn != nil {
		return r.failOn(cmd)
	}
	return nil
}

func (r *recordingRunner) recorded() []string {
	r.mu.Lock()
	defer r.mu.Unlock()
	return slices.Clone(r.commands)
}

func newInstaller() (installer, *recordingRunner) {
	runner := &recordingRunner{}
	return installer{run: runner.run, installed: func(string) bool { return false }}, runner
}

func TestInstaller_RecordsEveryCommand(t *testing.T) {
	i, runner := newInstaller()

	if err := i.apprise(""); err != nil {
		t.Fatal(err)
	}

	recorded := runner.recorded()
	if len(recorded) != 1 {
		t.Fatalf("expected 1 command, got %d", len(recorded))
	}
	if !strings.Contains(recorded[0], "pipx install apprise") {
		t.Fatalf("unexpected command: %s", recorded[0])
	}
}

func TestInstaller_PinsVersionWhenProvided(t *testing.T) {
	tests := []struct {
		name    string
		install func(i installer) error
		want    []string
	}{
		{
			name:    "apprise",
			install: func(i installer) error { return i.apprise("1.9.0") },
			want:    []string{"pipx install apprise==1.9.0"},
		},
		{
			name:    "borgbackup",
			install: func(i installer) error { return i.borgBackup("1.4.0") },
			want:    []string{"apt-get install -y borgbackup=1.4.0"},
		},
		{
			name:    "git",
			install: func(i installer) error { return i.git("2.47.0") },
			want:    []string{"apt-get -y install git=2.47.0"},
		},
		{
			name:    "podman",
			install: func(i installer) error { return i.podman("5.0") },
			want:    []string{"apt-get -y install podman-compose podman=5.0"},
		},
		{
			name:    "rdiff-backup",
			install: func(i installer) error { return i.rdiffBackup("2.2.6") },
			want:    []string{"apt-get -y install rdiff-backup=2.2.6"},
		},
		{
			name:    "rsync",
			install: func(i installer) error { return i.rsync("3.4.0") },
			want:    []string{"apt-get -y install rsync=3.4.0"},
		},
		{
			name:    "logrotate",
			install: func(i installer) error { return i.logrotate("3.21.0") },
			want:    []string{"apt-get -y install logrotate=3.21.0"},
		},
		{
			name:    "sqlite3",
			install: func(i installer) error { return i.sqlite3("3.46.0") },
			want:    []string{"apt-get -y install sqlite3=3.46.0"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i, runner := newInstaller()
			if err := tt.install(i); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			recorded := runner.recorded()
			if len(recorded) != len(tt.want) || !slices.Equal(recorded, tt.want) {
				t.Fatalf("expected %v, got %v", tt.want, recorded)
			}
		})
	}
}

func TestInstaller_DockerSteps(t *testing.T) {
	i, runner := newInstaller()

	if err := i.docker("27.0.1"); err != nil {
		t.Fatal(err)
	}

	recorded := runner.recorded()
	if len(recorded) != 6 {
		t.Fatalf("expected 6 commands, got %d: %v", len(recorded), recorded)
	}
	if !strings.Contains(recorded[0], "install -m 0755 -d /etc/apt/keyrings") {
		t.Fatalf("expected keyring directory setup, got %s", recorded[0])
	}
	if !strings.Contains(recorded[4], "apt-get update") {
		t.Fatalf("expected package update before install, got %s", recorded[4])
	}
	if !strings.Contains(recorded[5], "docker-ce-cli=27.0.1") {
		t.Fatalf("expected pinned docker-ce-cli, got %s", recorded[5])
	}
}

func TestInstaller_DockerWithoutVersion(t *testing.T) {
	i, runner := newInstaller()

	if err := i.docker(""); err != nil {
		t.Fatal(err)
	}

	recorded := runner.recorded()
	if !strings.Contains(recorded[5], "apt-get install -y docker-ce-cli docker-compose-plugin") {
		t.Fatalf("expected unpinned docker install, got %s", recorded[5])
	}
}

func TestInstaller_ResticSelfUpdates(t *testing.T) {
	i, runner := newInstaller()

	if err := i.restic("0.17.0"); err != nil {
		t.Fatal(err)
	}

	recorded := runner.recorded()
	if len(recorded) != 2 {
		t.Fatalf("expected 2 commands, got %d: %v", len(recorded), recorded)
	}
	if !strings.Contains(recorded[0], "restic=0.17.0") {
		t.Fatalf("expected pinned restic, got %s", recorded[0])
	}
	if !strings.Contains(recorded[1], "restic self-update") {
		t.Fatalf("expected self-update, got %s", recorded[1])
	}
}

func TestInstaller_RcloneUsesScriptWithoutVersion(t *testing.T) {
	i, runner := newInstaller()

	if err := i.rclone(""); err != nil {
		t.Fatal(err)
	}

	recorded := runner.recorded()
	if len(recorded) != 1 || !strings.Contains(recorded[0], "rclone.org/install.sh") {
		t.Fatalf("expected install script, got %v", recorded)
	}
}

func TestInstaller_RcloneUsesAptWithVersion(t *testing.T) {
	i, runner := newInstaller()

	if err := i.rclone("1.67.0"); err != nil {
		t.Fatal(err)
	}

	recorded := runner.recorded()
	if len(recorded) != 1 || !strings.Contains(recorded[0], "apt-get install -y rclone=1.67.0") {
		t.Fatalf("expected pinned rclone apt install, got %v", recorded)
	}
}

func TestInstaller_KopiaSteps(t *testing.T) {
	i, runner := newInstaller()

	if err := i.kopia("0.17.0"); err != nil {
		t.Fatal(err)
	}

	recorded := runner.recorded()
	if len(recorded) != 6 {
		t.Fatalf("expected 6 commands, got %d: %v", len(recorded), recorded)
	}
	if !strings.Contains(recorded[5], "kopia=0.17.0") {
		t.Fatalf("expected pinned kopia, got %s", recorded[5])
	}
}

func TestSupported_CoversAllDocumentedSoftware(t *testing.T) {
	supported := installer{}.supported()

	want := []string{
		"apprise", "borgbackup", "docker", "git", "podman", "rclone",
		"rdiff-backup", "restic", "rsync", "logrotate", "sqlite3", "kopia",
	}

	for _, name := range want {
		get, ok := supported[name]
		if !ok {
			t.Fatalf("expected %q to be supported", name)
		}
		if get == nil {
			t.Fatalf("expected %q to have an installer", name)
		}
	}
}

func TestInstall_FlowUpdatesInstallsAndCleansUp(t *testing.T) {
	i, runner := newInstaller()

	list := []Software{{Name: "git", Version: "2.47.0"}}
	i.install(list)

	recorded := runner.recorded()
	if len(recorded) != 3 {
		t.Fatalf("expected 3 commands (update, install, cleanup), got %d: %v", len(recorded), recorded)
	}
	if recorded[0] != "apt-get update" {
		t.Fatalf("expected update first, got %s", recorded[0])
	}
	if !strings.Contains(recorded[1], "git=2.47.0") {
		t.Fatalf("expected pinned git install, got %s", recorded[1])
	}
	if !strings.Contains(recorded[2], "rm -rf") {
		t.Fatalf("expected cleanup last, got %s", recorded[2])
	}
}

func TestInstall_FlowSkipsAlreadyInstalled(t *testing.T) {
	runner := &recordingRunner{}
	i := installer{run: runner.run, installed: func(name string) bool { return name == "git" }}

	list := []Software{{Name: "git", Version: ""}}
	i.install(list)

	recorded := runner.recorded()
	for _, cmd := range recorded {
		if strings.Contains(cmd, "apt-get -y install git") {
			t.Fatalf("git should be skipped, but ran %s", cmd)
		}
	}
}

func TestInstall_FlowContinuesAfterFailure(t *testing.T) {
	runner := &recordingRunner{
		failOn: func(cmd string) error {
			if strings.Contains(cmd, "git=") {
				return os.ErrDeadlineExceeded
			}
			return nil
		},
	}
	i := installer{run: runner.run, installed: func(string) bool { return false }}

	list := []Software{
		{Name: "git", Version: "2.47.0"},
		{Name: "rsync", Version: ""},
	}
	i.install(list)

	recorded := runner.recorded()
	foundGit := false
	foundRsync := false
	for _, cmd := range recorded {
		if strings.Contains(cmd, "git=") {
			foundGit = true
		}
		if strings.Contains(cmd, "rsync") {
			foundRsync = true
		}
	}
	if !foundGit {
		t.Fatal("expected git install to be attempted")
	}
	if !foundRsync {
		t.Fatal("expected rsync install to continue after git failed")
	}
}

func TestInstall_FlowSkipsUnsupportedSoftware(t *testing.T) {
	i, runner := newInstaller()

	list := []Software{{Name: "not-a-package", Version: ""}}
	i.install(list)

	recorded := runner.recorded()
	if len(recorded) != 2 {
		t.Fatalf("expected only update and cleanup, got %v", recorded)
	}
}

func TestOSReleaseIsDebian(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    bool
	}{
		{"debian id", "ID=debian\n", true},
		{"quoted debian", `ID="debian"` + "\n", true},
		{"indented debian", "\t ID=debian \n", true},
		{"ubuntu", "ID=ubuntu\n", false},
		{"prefix only", "ID_LIKE=debian\n", false},
		{"empty", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			file := filepath.Join(t.TempDir(), "os-release")
			if err := os.WriteFile(file, []byte(tt.content), 0o644); err != nil {
				t.Fatal(err)
			}

			if got := osReleaseIsDebian(file); got != tt.want {
				t.Fatalf("expected %v, got %v", tt.want, got)
			}
		})
	}
}

func TestIsInstalled(t *testing.T) {
	if !isInstalled("sh") {
		t.Fatal("expected sh to be installed")
	}
	if isInstalled("definitely-not-a-real-binary-xyz") {
		t.Fatal("expected unknown binary to not be installed")
	}
}

func TestCanInstallOnOS_RejectsEachConditionIndependently(t *testing.T) {
	if runtime.GOOS != "linux" {
		if canInstallOnOS("/etc/os-release") {
			t.Fatal("expected non-Linux OS to be rejected")
		}
		return
	}

	ubuntu := filepath.Join(t.TempDir(), "os-release")
	if err := os.WriteFile(ubuntu, []byte("ID=ubuntu\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if canInstallOnOS(ubuntu) {
		t.Fatal("expected non-Debian Linux to be rejected")
	}

	debian := filepath.Join(t.TempDir(), "os-release")
	if err := os.WriteFile(debian, []byte("ID=debian\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if !canInstallOnOS(debian) {
		t.Fatal("expected Debian to be allowed")
	}
}
