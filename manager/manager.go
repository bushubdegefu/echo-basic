package manager

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	goBlueCmd = &cobra.Command{
		Use:           "goblue",
		Short:         "goblue – command-line tool to aid structure you fiber backend projects with gorm",
		Long:          ``,
		Version:       "0.0.0",
		SilenceErrors: true,
		SilenceUsage:  true,
	}
)

func Execute() {
	if err := goBlueCmd.Execute(); err != nil {

		fmt.Println(err)
		os.Exit(1)
	}
}
