package config

import (
	"path/filepath"
	"testing"
)

func TestSetConfigFolderPath_SetsConfigFilePath(t *testing.T) {
	previous := GetConfigFilePath()
	t.Cleanup(func() { SetConfigFilePath(previous) })

	SetConfigFolderPath("./tmp-config")

	expected := filepath.Clean("tmp-config/config.yaml")
	if got := GetConfigFilePath(); got != expected {
		t.Fatalf("unexpected config file path: got %q want %q", got, expected)
	}
}

func TestSetConfigFilePath_UsesDefaultWhenEmpty(t *testing.T) {
	previous := GetConfigFilePath()
	t.Cleanup(func() { SetConfigFilePath(previous) })

	SetConfigFilePath("")

	expected := filepath.Clean(GetDefaultConfigFile())
	if got := GetConfigFilePath(); got != expected {
		t.Fatalf("unexpected default config file path: got %q want %q", got, expected)
	}
}

func TestTerminalSettingsHydrate_PopulatesAllowedArgsMap(t *testing.T) {
	settings := TerminalSettings{
		AllowedCommands: map[string]AllowedCommands{
			"docker": {
				Args: []string{"ps", "version"},
			},
		},
	}

	settings.Hydrate()

	docker := settings.AllowedCommands["docker"]
	if docker.AllowedArgsMap == nil {
		t.Fatal("expected allowed args map to be initialized")
	}
	if _, ok := docker.AllowedArgsMap["ps"]; !ok {
		t.Fatal("expected argument ps to be present in allowed args map")
	}
	if _, ok := docker.AllowedArgsMap["version"]; !ok {
		t.Fatal("expected argument version to be present in allowed args map")
	}
}
