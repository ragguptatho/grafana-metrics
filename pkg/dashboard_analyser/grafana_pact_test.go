package dashboard_analyser

import (
	"encoding/json"
	"errors"
	"os"
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


	var (
		dashboardOutputFile string 
	)

	if os.Getenv("DASH_OUTPUT_FILE") == ""{
		t.Errorf("Please provide the dashboard output file name using DASH_OUTPUT_FILE variable.")
	}

	dashboardOutputFile = os.Getenv("DASH_OUTPUT_FILE")


	analysedMetrics := pact.AddMessage()


	buffer, err := loadFile(dashboardOutputFile)

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
