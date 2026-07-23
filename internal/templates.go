package internal

import (
	"embed"
	"fmt"
	"os"
	"path/filepath"
	"text/template"
)

//go:embed templates/*
var templateFiles embed.FS

type ProjectData struct {
	ProjectName   string
	PythonVersion string
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

func GenerateInfrastructure(projectPath, projectName string, useDocker bool, pythonVersion string) error {
	data := ProjectData{
		ProjectName:   projectName,
		PythonVersion: pythonVersion + "-slim",
	}
	if err := GenerateFile("gitignore.tmpl", filepath.Join(projectPath, ".gitignore"), data); err != nil {
		return err
	}

	if useDocker {
		if err := GenerateFile("docker.tmpl", filepath.Join(projectPath, "Dockerfile"), data); err != nil {
			return err
		}
		if err := GenerateFile("compose.tmpl", filepath.Join(projectPath, "docker-compose.yml"), data); err != nil {
			return err
		}
		if err := GenerateFile("dockerignore.tmpl", filepath.Join(projectPath, ".dockerignore"), data); err != nil {
			return err
		}
	}

	return nil
}
