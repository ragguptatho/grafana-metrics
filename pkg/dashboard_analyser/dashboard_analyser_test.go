package dashboard_analyser

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"os"
	"testing"

	"github.com/grafana-tools/sdk"
	"github.com/stretchr/testify/assert"
)

func TestNewDashboardAnalyser_Analyse(t *testing.T) {

	// preparing dashboard
	board := sdk.NewBoard("Sample dashboard")
	board.ID = 1
	board.Time.From = "now-30m"
	board.Time.To = "now"
	row := board.AddRow("Sample row title")
	graphOne := sdk.NewGraph("Sample graph")
	metrics := []string{"node_boot_time_seconds", "node_context_switches_total", "process_resident_memory_max_bytes"}
	labels := []string{"instance", "job"}

	targets := []sdk.Target{
		{
			RefID:      "A",
			Expr:       fmt.Sprintf("%s{%s=\"\",%s=\"\"}", metrics[0], labels[0], labels[1]),
			Datasource: "prometheus",
			Instant:    true,
		},
		{
			RefID:      "B",
			Expr:       fmt.Sprintf("%s{%s=\"\",%s=\"\"}", metrics[1], labels[0], labels[1]),
			Datasource: "prometheus",
			Instant:    true,
		},
		{
			RefID:      "C",
			Expr:       fmt.Sprintf("%s{%s=\"\",%s=\"\"}", metrics[2], labels[0], labels[1]),
			Datasource: "prometheus",
			Instant:    true,
		},
	}
	graphOne.AddTarget(&targets[0])
	graphOne.AddTarget(&targets[1])
	graphOne.AddTarget(&targets[2])

	row.Add(graphOne)

	// writing the dashboard to a file
	boardContent, err := json.Marshal(board)
	if err != nil {
		t.Errorf("got some error while marhsalling the dashboard %s",err)
	}

	tempDir := os.TempDir()
	inputFile, err := os.CreateTemp(tempDir,"sample_dashboard.json")
	defer os.Remove(inputFile.Name())
	if err != nil {
		t.Errorf("got some error while creating input temp file %s",err)
	}
	outputFile, err := os.CreateTemp(tempDir,"sample_dashboard_output.json")
	if err != nil {
		t.Errorf("got some error while creating output temp file %s",err)
	}
	defer os.Remove(outputFile.Name())

	dashInputFile := inputFile.Name()
	dashOutputFile := outputFile.Name()

	if err != nil {
		t.Errorf("got some error %s", err)
	}

	err = os.WriteFile(dashInputFile, boardContent, fs.FileMode(int(0666)))

	if err != nil {
		t.Errorf("got some error %s", err)
	}

	dashInputFiles := []string{dashInputFile}

	dashboardAnalyser := DashboardAnalyser{
		DashFilesList: dashInputFiles,
		DashOutFile:   dashOutputFile,
	}
	// calling the analyse fn which will read the input file and output the content in an outputfile
	if err := dashboardAnalyser.Analyse(); err!=nil{
		t.Errorf("got some error %s", err)
	}

	actualContent, err := os.ReadFile(dashOutputFile)

	if err != nil {
		t.Errorf("got some error %s", err)
	}

	var metricsInGrafana ConsumerMetrics

	if err = json.Unmarshal(actualContent, &metricsInGrafana); err != nil {
		t.Errorf("got some error %s", err)
	}

	if metricsInGrafana.Metrics == nil {
		t.Errorf("expected metrics to be not nil")
	}

	assert.Equal(t, len(metricsInGrafana.Metrics), len(targets))

	metricNames := make([]string, 0)
	labelKeys := make(map[string]void)

	for metric, labels := range metricsInGrafana.Metrics {
		metricNames = append(metricNames, metric)
		for each_label := range labels.LabelKeys {
			labelKeys[each_label] = void{}
		}
	}
	assert.Subset(t, metrics, metricNames)
	assert.Equal(t, len(labelKeys), len(labels))
}