package internal

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func getVenvExec(projectName, execName string) string {
	var relPath string
	if runtime.GOOS == "windows" {
		relPath = filepath.Join(projectName, "env", "Scripts", execName+".exe")
	} else {
		relPath = filepath.Join(projectName, "env", "bin", execName)
	}

	absPath, err := filepath.Abs(relPath)
	if err != nil {
		return relPath
	}
	return absPath
}

func getPythonVersion(pythonExec string) string {
	cmd := exec.Command(pythonExec, "-c", "import sys; print(f'{sys.version_info.major}.{sys.version_info.minor}')")
	out, err := cmd.Output()
	if err != nil {
		return "3.11"
	}
	return strings.TrimSpace(string(out))
}

func runCommandQuiet(workDir string, name string, args ...string) error {
	cmd := exec.Command(name, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("%w | Лог: %s", err, strings.TrimSpace(string(output)))
	}
	return nil
}

func runCommandOutput(workDir string, name string, args ...string) ([]byte, error) {
	cmd := exec.Command(name, args...)
	if workDir != "" {
		cmd.Dir = workDir
	}
	return cmd.Output()
}

func executeStep(icon, message string, action func() error) error {
	fmt.Printf("%s %s... ", icon, message)
	start := time.Now()
	err := action()
	duration := time.Since(start).Round(time.Millisecond)
	if err != nil {
		fmt.Printf("❌ Ошибка (%v)\n", duration)
		return err
	}
	fmt.Printf("✨ (%v)\n", duration)
	return nil
}

func launchPyCharmAtTheEnd(projectPath string) {
	fmt.Print("\n🚀 Запускаем PyCharm...")

	var pycharmCmd string
	names := []string{"pycharm64", "pycharm64.exe", "pycharm", "pycharm.cmd"}
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			pycharmCmd = path
			break
		}
	}

	if pycharmCmd == "" {
		fmt.Println("\n⚠️ PyCharm не найден в PATH. Открой проект вручную!")
		return
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", "", pycharmCmd, projectPath)
	} else {
		cmd = exec.Command(pycharmCmd, projectPath)
	}

	if err := cmd.Start(); err != nil {
		fmt.Printf("\n⚠️ Ошибка запуска PyCharm: %v\n", err)
	} else {
		fmt.Println(" Ок!")
	}
}

func CreateProject(projectName string, apps []string, useDocker bool, useDuplicate bool) {
	fmt.Printf("\n🔥 Конструируем Django-проект: [%s]\n\n", projectName)

	if err := os.MkdirAll(projectName, 0755); err != nil {
		log.Fatalf("❌ Ошибка создания папки: %v", err)
	}

	err := executeStep("🐍", "Создаем изолированное окружение", func() error {
		return runCommandQuiet(projectName, "python", "-m", "venv", "env")
	})
	if err != nil {
		log.Fatalf("\n%v", err)
	}

	venvPython := getVenvExec(projectName, "python")
	venvPip := getVenvExec(projectName, "pip")

	err = executeStep("📦", "Устанавливаем Django, psycopg2 и decouple", func() error {
		return runCommandQuiet("", venvPip, "install", "--no-cache-dir", "django", "psycopg2-binary", "python-decouple")
	})
	if err != nil {
		log.Fatalf("\n%v", err)
	}

	err = executeStep("📄", "Генерируем requirements.txt", func() error {
		freezeData, err := runCommandOutput("", venvPip, "freeze")
		if err == nil {
			reqPath := filepath.Join(projectName, "requirements.txt")
			os.WriteFile(reqPath, freezeData, 0644)
		}
		return err
	})
	if err != nil {
		log.Fatalf("\n%v", err)
	}

	var settingsPath, urlsPath string
	if useDuplicate {
		err = executeStep("🏗 ", "Генерируем классическую структуру", func() error {
			return runCommandQuiet("", venvPython, "-m", "django", "startproject", projectName)
		})
		settingsPath = filepath.Join(projectName, projectName, "settings.py")
		urlsPath = filepath.Join(projectName, projectName, "urls.py")
	} else {
		err = executeStep("🏗 ", "Создаем плоскую архитектуру (core)", func() error {
			return runCommandQuiet(projectName, venvPython, "-m", "django", "startproject", "core", ".")
		})
		settingsPath = filepath.Join(projectName, "core", "settings.py")
		urlsPath = filepath.Join(projectName, "core", "urls.py")
	}
	if err != nil {
		log.Fatalf("\n%v", err)
	}

	if err := executeStep("🔐", "Генерируем криптостойкий SECRET_KEY и настраиваем .env", func() error {
		if err := createEnvFiles(projectName); err != nil {
			return fmt.Errorf("ошибка генерации .env файлов: %v", err)
		}
		if err := patchSettingsForEnv(settingsPath); err != nil {
			return fmt.Errorf("ошибка настройки python-decouple: %v", err)
		}
		return nil
	}); err != nil {
		log.Fatalf("\n%v", err)
	}

	pyVersion := getPythonVersion(venvPython)
	if err := executeStep("🐳", "Разворачиваем .gitignore, .dockerignore и Dockerfile", func() error {
		return GenerateInfrastructure(projectName, projectName, useDocker, pyVersion)
	}); err != nil {
		log.Fatalf("\n%v", err)
	}

	if len(apps) > 0 {
		for _, app := range apps {
			if err := executeStep("🚀", fmt.Sprintf("Создаем аппку [%s] и генерируем urls.py", app), func() error {
				if err := runCommandQuiet(projectName, venvPython, "manage.py", "startapp", app); err != nil {
					return fmt.Errorf("ошибка startapp для %s: %v", app, err)
				}
				return createAppUrls(projectName, app)
			}); err != nil {
				log.Fatalf("\n%v", err)
			}
		}

		if err := executeStep("⚙️ ", "Интегрируем аппки в settings.py (INSTALLED_APPS)", func() error {
			return addAppsToSettings(settingsPath, apps)
		}); err != nil {
			log.Fatalf("\n%v", err)
		}

		if err := executeStep("🔗", "Связываем роутинги в главном urls.py", func() error {
			return addAppsToUrls(urlsPath, apps)
		}); err != nil {
			log.Fatalf("\n%v", err)
		}
	}

	if err := executeStep("📚", "Инициализируем Git-репозиторий", func() error {
		return runCommandQuiet(projectName, "git", "init")
	}); err != nil {
		log.Fatalf("\n%v", err)
	}
	fmt.Println("\n==================================================")
	fmt.Printf("🎉 Проект %s успешно собран и готов к разработке!\n", projectName)
	fmt.Println("==================================================")

	launchPyCharmAtTheEnd(projectName)
}
