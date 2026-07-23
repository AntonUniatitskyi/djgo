package internal

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"os"
	"path/filepath"
	"strings"
)

func generateSecretKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*(-_=+)"
	key := make([]byte, 50)
	for i := range key {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		key[i] = charset[num.Int64()]
	}
	return string(key)
}

func createEnvFiles(projectPath string) error {
	secretKey := generateSecretKey()
	envContent := fmt.Sprintf("SECRET_KEY=%s\nDEBUG=True\n", secretKey)
	if err := os.WriteFile(filepath.Join(projectPath, ".env"), []byte(envContent), 0644); err != nil {
		return fmt.Errorf("ошибка создания .env: %w", err)
	}
	envExampleContent := "SECRET_KEY=your-secret-key-here\nDEBUG=True\n"
	if err := os.WriteFile(filepath.Join(projectPath, ".env.example"), []byte(envExampleContent), 0644); err != nil {
		return fmt.Errorf("ошибка создания .env.example: %w", err)
	}
	return nil
}

func patchSettingsForEnv(settingsPath string) error {
	content, err := os.ReadFile(settingsPath)
	if err != nil {
		return fmt.Errorf("не удалось прочитать settings.py: %w", err)
	}
	lines := strings.Split(string(content), "\n")
	var newLines []string
	importAdded := false
	for _, line := range lines {
		if !importAdded && strings.Contains(line, "from pathlib import Path") {
			newLines = append(newLines, line)
			newLines = append(newLines, "from decouple import config")
			importAdded = true
			continue
		}
		if strings.HasPrefix(line, "SECRET_KEY =") {
			newLines = append(newLines, "SECRET_KEY = config('SECRET_KEY')")
			continue
		}
		if strings.HasPrefix(line, "DEBUG =") {
			newLines = append(newLines, "DEBUG = config('DEBUG', default=False, cast=bool)")
			continue
		}
		newLines = append(newLines, line)
	}
	output := strings.Join(newLines, "\n")
	return os.WriteFile(settingsPath, []byte(output), 0644)
}
