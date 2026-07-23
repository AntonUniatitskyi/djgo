package internal

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

var templateFiles embed.FS

type ProjectData struct {
	ProjectName string
}

func GenerateFile(tmplName, outPath string, data ProjectData) error {
	tmpl, err := template.ParseFS(templateFiles, "templates/"+tmplName)
	if err != nil {
		return fmt.Errorf("ошибка чтения шаблона %s: %w", tmplName, err)
	}
	file, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("ошибка создания файла %s: %w", outPath, err)
	}
	defer file.Close()
	return tmpl.Execute(file, data)
}

func GenerateInfrastructure(projectPath, projectName string, useDocker bool) error {
	data := ProjectData{ProjectName: projectName}
	GenerateFile("gitignore.tmpl", filepath.Join(projectPath, ".gitignore"), data)
	if useDocker {
		GenerateFile("docker.tmpl", filepath.Join(projectPath, "Dockerfile"), data)
		GenerateFile("compose.tmpl", filepath.Join(projectPath, "docker-compose.yml"), data)
		GenerateFile("dockerignore.tmpl", filepath.Join(projectPath, ".dockerignore"), data)
	}

	return nil
}
