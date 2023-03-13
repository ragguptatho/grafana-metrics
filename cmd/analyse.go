package cmd

import (
	"log"

	"github.com/shakirshakiel/grafana-metrics-analyser/pkg/dashboard_analyser"
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(analyseCommand)
	analyseCommand.Flags().StringSliceVar(&globalFlags.dashFilesList, flagDashFilesList, []string{}, "dashboard files to analyze")
	analyseCommand.Flags().StringVar(&globalFlags.dashOutFile, flagDashOutFile, "", "output dashboard file path")
	if err := analyseCommand.MarkFlagRequired(flagDashFilesList); err != nil {
		log.Println(err)
	}
	if err := analyseCommand.MarkFlagRequired(flagDashOutFile); err != nil {
		log.Println(err)
	}
}

var analyseCommand = &cobra.Command{
	Use:   "analyse",
	Short: "to analyse prometheus metrics in grafana dashboards",
	Long:  "to analyse prometheus metrics in grafana dashboards",
	RunE: func(cmd *cobra.Command, args []string) error {
		dashboardAnalyser := dashboard_analyser.NewDashboardAnalyser(globalFlags.dashFilesList, globalFlags.dashOutFile)
		return dashboardAnalyser.Analyse()
	},
}
