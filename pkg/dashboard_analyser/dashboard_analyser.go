package dashboard_analyser

import (
	"encoding/json"
	"fmt"
	"os"
	"log"
	"github.com/grafana-tools/sdk"
	"io/ioutil"
)

type DashboardAnalyser struct{
    DashFilesList []string
}

// NewGenerator returns an instance of a dashboard
func NewDashboardAnalyser(dashFilesList []string) *DashboardAnalyser {
	return &DashboardAnalyser{
	    DashFilesList: dashFilesList,
	}
}

// GetDashboardFiles gets all the dashboards in a given path
func (dashboardAnalyser *DashboardAnalyser) Analyse() error {
    output := &MetricsInGrafana{}
    output.OverallMetrics = make([]Metric, 0)

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
		log.Println(board.ID)
		ParseMetricsInBoard(output, board)
	}

	err := writeOut(output, "grafana-metrics-analyser.json")
    if err != nil {
    	return err
    }

    return nil
}

func loadFile(filename string) ([]byte, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	fileinfo, err := file.Stat()
	if err != nil {
		return nil, err
	}

	filesize := fileinfo.Size()
	buffer := make([]byte, filesize)

	_, err = file.Read(buffer)
	if err != nil {
		return nil, err
	}

	return buffer, nil
}

func writeOut(mig *MetricsInGrafana, outputFile string) error {
	metricsUsed := make([]Metric, 0)
	for _, metric := range mig.OverallMetrics {
		metricsUsed = append(metricsUsed, metric)
	}
	mig.MetricsUsed = metricsUsed

	out, err := json.MarshalIndent(mig, "", "  ")
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(outputFile, out, os.FileMode(int(0666))); err != nil {
		return err
	}

	return nil
}
