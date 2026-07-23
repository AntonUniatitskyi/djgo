package internal

import (
	"fmt"
	"os"
	"strings"
)

func addAppsToSettings(settingsPath string, apps []string) error {
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать settings.py: %w", err)
	}
	lines := strings.Split(string(content), "\n")
	var newLines []string
	inInstalledApps := false
	inserted := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.Contains(line, "INSTALLED_APPS = [") {
			inInstalledApps = true
		}

		if inInstalledApps && !inserted && trimmed == "]" {
			for _, app := range apps {
				newAppLine := fmt.Sprintf("    '%s',", app)
				newLines = append(newLines, newAppLine)
			}
			inserted = true
			inInstalledApps = false
		}

		newLines = append(newLines, line)
	}

	if !inserted {
		return fmt.Errorf("не удалось найти закрывающую скобку блока INSTALLED_APPS")
	}

	output := strings.Join(newLines, "\n")
	if err := os.WriteFile(settingsPath, []byte(output), 0644); err != nil {
		return fmt.Errorf("ошибка записи в settings.py: %w", err)
	}

	return nil
}
