package dashboard_analyser

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
)

func createPact() dsl.Pact {
	return dsl.Pact{
		Consumer: "grafana",
		Provider: "prometheus",
		LogLevel: "INFO",
	}
}

var pact dsl.Pact = createPact()

func TestGrafana(t *testing.T) {

	analysedMetrics := pact.AddMessage()

	dashFile := "../../sample/outputs/node-expoter-output.json"

	buffer, err := loadFile(dashFile)

	var metricsInGrafana ConsumerMetrics

	json.Unmarshal(buffer, &metricsInGrafana)

	if err != nil {
		t.Errorf("got some error %s", err)
	}

	analysedMetrics.
		Given("analyse metrics").
		ExpectsToReceive("analyse metrics").
		WithContent(metricsInGrafana).AsType(&ConsumerMetrics{})

	pact.VerifyMessageConsumer(t, analysedMetrics, func(m dsl.Message) error {

		ok := m.Content.(*ConsumerMetrics)
		if ok == nil {
			return errors.New("Content is nil")
		}

		return nil

	})

}
