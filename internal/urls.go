package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const appUrlsTemplate = `from django.urls import path
from . import views

urlpatterns = [
    # path('', views.index, name='index'),
]
`

func createAppUrls(projectName, appName string) error {
	appUrlsPath := filepath.Join(projectName, appName, "urls.py")
	return os.WriteFile(appUrlsPath, []byte(appUrlsTemplate), 0644)
}

func addAppsToUrls(urlsPath string, apps []string) error {
	content, err := os.ReadFile(urlsPath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать urls.py: %w", err)
	}

	lines := strings.Split(string(content), "\n")
	var newLines []string
	inUrlPatterns := false
	inserted := false
	importUpdated := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if !importUpdated && strings.HasPrefix(trimmed, "from django.urls import") {
			if !strings.Contains(line, "include") {
				line = strings.Replace(line, "import path", "import path, include", 1)
			}
			importUpdated = true
		}
		if strings.Contains(line, "urlpatterns = [") {
			inUrlPatterns = true
		}
		if inUrlPatterns && !inserted && trimmed == "]" {
			for _, app := range apps {
				newRoute := fmt.Sprintf("    path('%s/', include('%s.urls')),", app, app)
				newLines = append(newLines, newRoute)
			}
			inserted = true
			inUrlPatterns = false
		}

		newLines = append(newLines, line)
	}

	if !inserted {
		return fmt.Errorf("не удалось найти закрывающую скобку блока urlpatterns")
	}

	output := strings.Join(newLines, "\n")
	return os.WriteFile(urlsPath, []byte(output), 0644)
}
