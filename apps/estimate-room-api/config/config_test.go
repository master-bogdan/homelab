package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfigReadsDotEnvFile(t *testing.T) {
	unsetEnv(t, "PORT")
	unsetEnv(t, "HOST")
	unsetEnv(t, "PASETO_SYMMETRIC_KEY")

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")

	content := "PORT=9000\nHOST=127.0.0.1\nPASETO_SYMMETRIC_KEY=\"0123456789abcdef0123456789abcdef\"\n"
	if err := os.WriteFile(envPath, []byte(content), 0644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(wd); chdirErr != nil {
			t.Fatalf("restore cwd: %v", chdirErr)
		}
	})

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	if cfg.Server.Port != "9000" {
		t.Fatalf("expected PORT from .env, got %q", cfg.Server.Port)
	}

	if cfg.Server.Host != "127.0.0.1" {
		t.Fatalf("expected HOST from .env, got %q", cfg.Server.Host)
	}

	if cfg.Server.PasetoSymmetricKey != "0123456789abcdef0123456789abcdef" {
		t.Fatalf("expected PASETO key from .env, got %q", cfg.Server.PasetoSymmetricKey)
	}
}

func TestLoadConfigDoesNotOverrideExistingEnvironment(t *testing.T) {
	t.Setenv("PORT", "7777")

	tmpDir := t.TempDir()
	envPath := filepath.Join(tmpDir, ".env")
	if err := os.WriteFile(envPath, []byte("PORT=9000\n"), 0644); err != nil {
		t.Fatalf("write .env: %v", err)
	}

	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() {
		if chdirErr := os.Chdir(wd); chdirErr != nil {
			t.Fatalf("restore cwd: %v", chdirErr)
		}
	})

	if err := os.Chdir(tmpDir); err != nil {
		t.Fatalf("chdir temp dir: %v", err)
	}

	cfg, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}

	if cfg.Server.Port != "7777" {
		t.Fatalf("expected existing environment to win, got %q", cfg.Server.Port)
	}
}

func unsetEnv(t *testing.T, key string) {
	t.Helper()

	value, exists := os.LookupEnv(key)
	if err := os.Unsetenv(key); err != nil {
		t.Fatalf("unset %s: %v", key, err)
	}

	t.Cleanup(func() {
		var err error
		if exists {
			err = os.Setenv(key, value)
		} else {
			err = os.Unsetenv(key)
		}
		if err != nil {
			t.Fatalf("restore %s: %v", key, err)
		}
	})
}
