package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const (
	flagDashFilesList = "dash-files-list"
	flagDashOutFile   = "dash-out-file"
)

type flags struct {
	dashFilesList []string
	dashOutFile   string
}

var (
	globalFlags flags
	rootCmd     = &cobra.Command{
		Use:   "grafana-metrics-analyzer",
		Short: "utility that helps to analyze prometheus metrics in grafana dashboards",
		Long:  "utility that helps to analyze prometheus metrics in grafana dashboards",
	}
)

func init() {
	globalFlags = flags{}
}

func Execute() {
	if _, err := rootCmd.ExecuteC(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
