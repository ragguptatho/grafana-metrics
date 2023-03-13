package dashboard_analyser

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/grafana-tools/sdk"
)

type DashboardAnalyser struct {
	DashFilesList []string
	DashOutFile   string
}

// NewDashboardAnalyser returns an instance of a dashboardAnalyser
func NewDashboardAnalyser(dashFilesList []string, dashOutFile string) *DashboardAnalyser {
	return &DashboardAnalyser{
		DashFilesList: dashFilesList,
		DashOutFile:   dashOutFile,
	}
}

// GetDashboardFiles gets all the dashboards in a given path
func (dashboardAnalyser *DashboardAnalyser) Analyse() error {
	output := &ConsumerMetrics{}
	output.Metrics = make(map[string]Metric)

	for _, file := range dashboardAnalyser.DashFilesList {
		var board sdk.Board
		buf, err := loadFile(file)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(buf, &board); err != nil {
			fmt.Fprintf(os.Stderr, "%s for %s\n", err, file)
			continue
		}
		ParseMetricsInBoard(output, board)
	}

	err := writeOut(output, dashboardAnalyser.DashOutFile)
	if err != nil {
		return err
	}

	return nil
}

// marshal the metrics into json format and write to the output file
func writeOut(mig *ConsumerMetrics, outputFile string) error {

	out, err := json.MarshalIndent(mig, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(outputFile, out, os.FileMode(int(0666))); err != nil {
		return err
	}

	return nil
}
