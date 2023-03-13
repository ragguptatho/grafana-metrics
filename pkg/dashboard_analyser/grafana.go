package dashboard_analyser

import (
	"fmt"
	"regexp"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/grafana-tools/sdk"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/promql/parser"
)

// parse the metrics into Board given the ConsumerMetrics
func ParseMetricsInBoard(mig *ConsumerMetrics, board sdk.Board) {
	var parseErrors []error
	metrics := make(map[string]Metric)

	// Iterate through all the panels and collect metrics
	for _, panel := range board.Panels {
		parseErrors = append(parseErrors, metricsFromPanel(*panel, metrics)...)
		if panel.RowPanel != nil {
			for _, subPanel := range panel.RowPanel.Panels {
				parseErrors = append(parseErrors, metricsFromPanel(subPanel, metrics)...)
			}
		}
	}

	// Iterate through all the rows and collect metrics
	for _, row := range board.Rows {
		for _, panel := range row.Panels {
			parseErrors = append(parseErrors, metricsFromPanel(panel, metrics)...)
		}
	}

	// Process metrics in templating
	parseErrors = append(parseErrors, metricsFromTemplating(board.Templating, metrics)...)

	var parseErrs []string
	for _, err := range parseErrors {
		parseErrs = append(parseErrs, err.Error())
	}
	
	for metric := range metrics {
		if metric == "" {
			continue
		}
		mig.Metrics[metric] = metrics[metric]
	}

	log.Print(parseErrs)
}

// parse the metrics from the templated panels
func metricsFromTemplating(templating sdk.Templating, metrics map[string]Metric) []error {
	parseErrors := []error{}
	for _, templateVar := range templating.List {
		if templateVar.Type != "query" {
			continue
		}
		if query, ok := templateVar.Query.(string); ok {
			// label_values
			if strings.Contains(query, "label_values") {
				re := regexp.MustCompile(`label_values\(([a-zA-Z0-9_]+)`)
				sm := re.FindStringSubmatch(query)
				// In case of really gross queries, like - https://github.com/grafana/jsonnet-libs/blob/e97ab17f67ab40d5fe3af7e59151dd43be03f631/hass-mixin/dashboard.libsonnet#L93
				if len(sm) > 0 {
					query = sm[1]
				}
			}
			// query_result
			if strings.Contains(query, "query_result") {
				re := regexp.MustCompile(`query_result\((.+)\)`)
				query = re.FindStringSubmatch(query)[1]
			}
			err := parseQuery(query, metrics)
			if err != nil {
				parseErrors = append(parseErrors, errors.Wrapf(err, "query=%v", query))
				log.Debugln("msg", "promql parse error", "err", err, "query", query)
				continue
			}
		} else {
			err := fmt.Errorf("templating type error: name=%v", templateVar.Name)
			parseErrors = append(parseErrors, err)
			log.Debugln("msg", "templating parse error", "err", err)
			continue
		}
	}
	return parseErrors
}

// parse the metrics from the normal panels
func metricsFromPanel(panel sdk.Panel, metrics map[string]Metric) []error {
	var parseErrors []error

	targets := panel.GetTargets()
	if targets == nil {
		parseErrors = append(parseErrors,
			fmt.Errorf("unsupported panel type: %q", panel.CommonPanel.Type))
		return parseErrors
	}

	for _, target := range *targets {
		// Prometheus has this set.
		if target.Expr == "" {
			continue
		}
		query := target.Expr
		err := parseQuery(query, metrics)
		if err != nil {
			parseErrors = append(parseErrors, errors.Wrapf(err, "query=%v", query))
			log.Debugln("msg", "promql parse error", "err", err, "query", query)
			continue
		}
	}

	return parseErrors
}

// parses the promQl
func parseQuery(query string, metrics map[string]Metric) error {
	query = strings.ReplaceAll(query, `$__interval`, "5m")
	query = strings.ReplaceAll(query, `$interval`, "5m")
	query = strings.ReplaceAll(query, `$resolution`, "5s")
	query = strings.ReplaceAll(query, "$__rate_interval", "15s")
	query = strings.ReplaceAll(query, "$__range", "1d")
	query = strings.ReplaceAll(query, "${__range_s:glob}", "30")
	query = strings.ReplaceAll(query, "${__range_s}", "30")
	expr, err := parser.ParseExpr(query)
	if err != nil {
		return err
	}

	parser.Inspect(expr, func(node parser.Node, path []parser.Node) error {
		if n, ok := node.(*parser.VectorSelector); ok {
			labelKeys := make(map[string]void)
			for _, matcher := range n.LabelMatchers {
				if matcher.Name != "__name__" {
					labelKeys[matcher.Name] = void{}
				}
			}
			// only keep the older labels which are not present in new labelKeys
			for label := range metrics[n.Name].LabelKeys {
				if _, ok := labelKeys[label]; !ok {
					labelKeys[label] = void{}
				}
			}
			metrics[n.Name] = Metric{LabelKeys: labelKeys}
		}

		return nil
	})
	return nil
}
