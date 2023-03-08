package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/shakirshakiel/grafana-metrics-analyser/pkg/dashboard_analyser"
)

const (
    flagDashFilesList = "dash-files-list"
)

type flags struct {
    dashFilesList []string
}

var (
	cmd         *cobra.Command
	globalFlags flags
)

func init() {
    globalFlags = flags{}
	cmd = &cobra.Command{
		Use:   "grafana-metrics-analyzer",
		Short: "utility that helps to analyze prometheus metrics in grafana dashboards",
		Long:  "utility that helps to analyze prometheus metrics in grafana dashboards",
		Args:  cobra.MinimumNArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
		    log.Printf("Hello")
		    dashboardAnalyser := dashboard_analyser.NewDashboardAnalyser(globalFlags.dashFilesList)
		    return dashboardAnalyser.Analyse()
		},
	}
	cmd.Flags().StringSliceVarP(&globalFlags.dashFilesList, flagDashFilesList, "", []string{}, "dashboard file to analyze")

	if err := cmd.MarkFlagRequired(flagDashFilesList); err != nil {
    	log.Println(err)
    }
}

// Main will take the workload of executing/starting the cli, when the command is passed to it.
func Main() {
	if err := execute(os.Args[1:]); err != nil {
		log.Fatal(err)
	}
}

func execute(args []string) error {
	cmd.SetArgs(args)
	_, err := cmd.ExecuteC()
	if err != nil {
		return err
	}
	return nil
}
