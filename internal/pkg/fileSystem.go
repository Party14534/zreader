package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetAppDataDir(appName string) (string, error) {
	var baseDir string
	var err error

    // Get the default path for where to store data
	switch runtime.GOOS {
	case "windows":
		baseDir = os.Getenv("AppData")
		if baseDir == "" {
			baseDir = os.Getenv("LocalAppData") // fallback
		}
	case "darwin":
		baseDir = filepath.Join(os.Getenv("HOME"), "Library", "Application Support")
	case "linux":
		baseDir = os.Getenv("XDG_DATA_HOME")
		if baseDir == "" {
			baseDir = filepath.Join(os.Getenv("HOME"), ".local", "share")
		}
	default:
		return "", fmt.Errorf("Unsupported platform")
	}

	if baseDir == "" {
		return "", fmt.Errorf("Unable to determine application data directory")
	}

	appDir := filepath.Join(baseDir, appName)
	if err = os.MkdirAll(appDir, os.ModePerm); err != nil {
		return "", fmt.Errorf("Failed to create application directory: %w", err)
	}

	return appDir, nil
}
