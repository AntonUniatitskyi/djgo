package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "djgo",
	Short: "djgo - умный скаффолдер для Django проектов",
	Long:  `Быстрый CLI-инструмент для генерации структуры Django с аппками и Docker-конфигом.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
