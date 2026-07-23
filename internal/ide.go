package internal

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

func findExecutable(names ...string) string {
	for _, name := range names {
		if path, err := exec.LookPath(name); err == nil {
			return path
		}
	}
	return ""
}

func startProcess(execPath, projectPath string) error {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "start", "", execPath, projectPath)
	} else {
		cmd = exec.Command(execPath, projectPath)
	}
	return cmd.Start()
}

func waitForPyCharm(projectPath string) {
	fmt.Print("⏳ Ждем, пока PyCharm развернет интерфейс проекта")
	workspaceFile := filepath.Join(projectPath, ".idea", "workspace.xml")
	timeout := time.After(20 * time.Second)
	ticker := time.NewTicker(250 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-timeout:
			fmt.Println("\n⚠️ Время ожидания истекло! PyCharm грузится слишком долго, продолжаю в фоне...")
			return
		case <-ticker.C:
			if _, err := os.Stat(workspaceFile); err == nil {
				fmt.Println("\n✅ PyCharm инициализировал рабочую область!")
				fmt.Println("☕️ Даем ему еще 2.5 секунды на отрисовку боковой панели...")
				time.Sleep(2500 * time.Millisecond)
				return
			}
			fmt.Print(".")
		}
	}
}

func chooseAndOpenIDE(projectPath string) {
	fmt.Println("\n🖥  В какой IDE открыть проект для визуализации?")
	fmt.Println("  1) 🟢 PyCharm (JetBrains)")
	fmt.Println("  2) 🔵 VS Code")
	fmt.Println("  3) 🟣 Cursor")
	fmt.Println("  0) ⚪️ Не открывать (продолжить в фоне)")
	fmt.Print("👉 Выбери цифру [по умолчанию 1 - PyCharm]: ")
	reader := bufio.NewReader(os.Stdin)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	switch input {
	case "2":
		fmt.Println("🚀 Запускаем VS Code...")
		if path := findExecutable("code", "code.cmd"); path != "" {
			startProcess(path, projectPath)
			time.Sleep(1000 * time.Millisecond)
		} else {
			fmt.Println("⚠️ Команда 'code' не найдена в PATH.")
		}
	case "3":
		fmt.Println("🚀 Запускаем Cursor...")
		if path := findExecutable("cursor", "cursor.cmd"); path != "" {
			startProcess(path, projectPath)
			time.Sleep(1000 * time.Millisecond)
		} else {
			fmt.Println("⚠️ Команда 'cursor' не найдена в PATH.")
		}
	case "0":
		fmt.Println("⏩ Пропускаем запуск IDE.")
		return
	default:
		fmt.Println("🚀 Запускаем PyCharm...")
		pycharmPath := findExecutable("pycharm64", "pycharm64.exe", "pycharm", "pycharm.cmd", "pycharm.bat")
		if pycharmPath == "" {
			fmt.Println("⚠️ Не удалось найти исполняемый файл PyCharm в PATH.")
			fmt.Println("ℹ️ Убедись, что путь к bin/ добавлен в переменные среды Windows, а пока продолжаем...")
			return
		}
		if err := startProcess(pycharmPath, projectPath); err != nil {
			fmt.Printf("⚠️ Ошибка запуска (%v). Продолжаем в консоли...\n", err)
			return
		}

		waitForPyCharm(projectPath)
	}
}
