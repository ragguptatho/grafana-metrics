package dashboard_analyser

import (
	"fmt"
	"testing"

	"github.com/grafana-tools/sdk"
	"github.com/stretchr/testify/assert"
)

func TestMetricsFromPanel(t *testing.T) {

	graph := sdk.NewGraph("")

	metrics := []string{"node_boot_time_seconds"}
	labels := []string{"instance", "job"}

	targets := []sdk.Target{
		{
			RefID:      "A",
			Expr:       fmt.Sprintf("%s{%s=\"\",%s=\"\"}", metrics[0], labels[0], labels[1]),
			Datasource: "prometheus",
			Instant:    true,
		},
	}

	graph.AddTarget(&targets[0])

	panel := &sdk.Panel{
		CommonPanel: sdk.CommonPanel{
			OfType: sdk.GraphType,
		},
		GraphPanel: graph.GraphPanel,
	}

	var metricsMap map[string]Metric = make(map[string]Metric)

	err := metricsFromPanel(*panel, metricsMap)

	if err != nil {
		t.Errorf("got some error %s", err)
	}

	metricNames := make([]string, 0)
	labelKeys := make(map[string]void)

	for metric, labels := range metricsMap {
		metricNames = append(metricNames, metric)
		for each_label := range labels.LabelKeys {
			labelKeys[each_label] = void{}
		}
	}
	assert.Subset(t, metrics, metricNames)
	assert.Equal(t, len(labelKeys), len(labels))

}

func TestMetricsFromTemplating(t *testing.T) {

	metrics := []string{"nodes"}
	labels := []string{"hostname"}

	query := fmt.Sprintf("label_values(%s,%s)", metrics[0], labels[0])

	template := sdk.Templating{
		List: []sdk.TemplateVar{{
			Type:       "query",
			Name:       "nodeName",
			Datasource: "prometheus",
			Query:      query,
		},
		},
	}

	var metricsMap map[string]Metric = make(map[string]Metric)

	err := metricsFromTemplating(template, metricsMap)

	if len(err) != 0 {
		t.Errorf("got some error %s", err)
	}

	metricNames := make([]string, 0)

	for metric := range metricsMap {
		metricNames = append(metricNames, metric)

	}

	assert.Subset(t, metrics, metricNames)

}