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

func showStep(icon, message string, delayMs int) {
	fmt.Printf("%s %s...", icon, message)
	time.Sleep(time.Duration(delayMs) * time.Millisecond)
	fmt.Println(" ✨")
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

	showStep("🐍", "Создаем изолированное окружение (python -m venv env)", 400)
	if err := runCommandQuiet(projectName, "python", "-m", "venv", "env"); err != nil {
		log.Fatalf("\n❌ Ошибка создания venv: %v", err)
	}

	venvPython := getVenvExec(projectName, "python")
	venvPip := getVenvExec(projectName, "pip")

	showStep("📦", "Устанавливаем Django (и psycopg2-binary для БД) внутрь env", 800)
	if err := runCommandQuiet("", venvPip, "install", "--no-cache-dir", "django", "psycopg2-binary", "python-decouple"); err != nil {
		log.Fatalf("\n❌ Ошибка установки пакетов через pip: %v", err)
	}

	showStep("📄", "Генерируем честный requirements.txt через pip freeze", 200)
	freezeData, err := runCommandOutput("", venvPip, "freeze")
	if err != nil {
		log.Fatalf("\n❌ Ошибка pip freeze: %v", err)
	}
	reqPath := filepath.Join(projectName, "requirements.txt")
	os.WriteFile(reqPath, freezeData, 0644)
	var settingsPath, urlsPath string
	if useDuplicate {
		showStep("🏗 ", "Генерируем классическую структуру (duplication)", 300)
		if err := runCommandQuiet("", venvPython, "-m", "django", "startproject", projectName); err != nil {
			log.Fatalf("\n❌ Ошибка startproject: %v", err)
		}
		settingsPath = filepath.Join(projectName, projectName, "settings.py")
		urlsPath = filepath.Join(projectName, projectName, "urls.py")
	} else {
		showStep("🏗 ", "Создаем плоскую архитектуру (core workspace)", 300)
		if err := runCommandQuiet(projectName, venvPython, "-m", "django", "startproject", "core", "."); err != nil {
			log.Fatalf("\n❌ Ошибка startproject: %v", err)
		}
		settingsPath = filepath.Join(projectName, "core", "settings.py")
		urlsPath = filepath.Join(projectName, "core", "urls.py")
	}
	showStep("🔐", "Генерируем криптостойкий SECRET_KEY и настраиваем .env", 300)
	if err := createEnvFiles(projectName); err != nil {
		log.Fatalf("\n❌ Ошибка генерации .env файлов: %v", err)
	}
	if err := patchSettingsForEnv(settingsPath); err != nil {
		log.Fatalf("\n❌ Ошибка настройки python-decouple в settings.py: %v", err)
	}
	showStep("🐳", "Разворачиваем .gitignore, .dockerignore и Dockerfile", 300)
	GenerateInfrastructure(projectName, projectName, useDocker)

	if len(apps) > 0 {
		for _, app := range apps {
			showStep("🚀", fmt.Sprintf("Создаем аппку [%s] и генерируем urls.py", app), 250)
			if err := runCommandQuiet(projectName, venvPython, "manage.py", "startapp", app); err != nil {
				log.Fatalf("\n❌ Ошибка startapp для %s: %v", app, err)
			}
			createAppUrls(projectName, app)
		}

		showStep("⚙️ ", "Интегрируем аппки в settings.py (INSTALLED_APPS)", 200)
		if err := addAppsToSettings(settingsPath, apps); err != nil {
			log.Fatalf("\n❌ Ошибка settings.py: %v", err)
		}

		showStep("🔗", "Связываем роутинги в главном urls.py", 200)
		if err := addAppsToUrls(urlsPath, apps); err != nil {
			log.Fatalf("\n❌ Ошибка urls.py: %v", err)
		}
	}

	showStep("📚", "Инициализируем Git-репозиторий", 150)
	runCommandQuiet(projectName, "git", "init")

	fmt.Println("\n==================================================")
	fmt.Printf("🎉 Проект %s успешно собран и готов к разработке!\n", projectName)
	fmt.Println("==================================================")

	launchPyCharmAtTheEnd(projectName)
}
