package dashboard_analyser

// to implement set as map[string]void
type void struct{}

type ConsumerMetrics struct {
	Metrics    map[string]Metric  `json:"metrics"`
}

type Metric struct {
	LabelKeys map[string]void `json:"label_keys"`
}
