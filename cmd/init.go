package cmd

import (
	"django-cli/internal"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init [название_проекта]",
	Short: "Генерирует каркас Django-проекта",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		projectName := args[0]
		apps, _ := cmd.Flags().GetStringSlice("apps")
		useDocker, _ := cmd.Flags().GetBool("docker")
		useDuplicate, _ := cmd.Flags().GetBool("duplicate")
		internal.CreateProject(projectName, apps, useDocker, useDuplicate)
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().StringSliceP("apps", "a", []string{}, "Список аппок через запятую (напр: users,api)")
	initCmd.Flags().BoolP("docker", "d", false, "Сгенерировать Dockerfile и compose")
	initCmd.Flags().BoolP("duplicate", "u", false, "Создать классическую структуру с дублем (project/project)")
}
