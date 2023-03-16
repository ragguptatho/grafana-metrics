package dashboard_analyser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/pact-foundation/pact-go/dsl"
	types "github.com/pact-foundation/pact-go/types"
)

func createPact() dsl.Pact {
	return dsl.Pact{
		Consumer: "grafana",
		Provider: "prometheus",
		LogLevel: "INFO",
		PactDir: "./pacts",
	}
}

var pact dsl.Pact = createPact()

func TestGrafana(t *testing.T) {

	var (
		dashboardOutputFile string 
		pactBrokerUrl string 
		consumerVersion string 
		current_env string 
	)

	dashboardOutputFile = os.Getenv("DASH_OUTPUT_FILE")

	pactBrokerUrl = os.Getenv("PACT_BROKER_URL")
	consumerVersion = os.Getenv("CONSUMER_TAG")
	current_env = os.Getenv("ENV")

	if dashboardOutputFile == ""{
		t.Errorf("Please provide the dashboard output file name using DASH_OUTPUT_FILE variable.")
	}

	if pactBrokerUrl == "" || consumerVersion == "" || current_env == ""{
		t.Errorf("Please provide the value for PACT_BROKER_URL,CONSUMER_TAG and ENV inorder to publish the contract to the pact broker")
	}

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

	p := dsl.Publisher{}

	// pact local path
	pactUrl := fmt.Sprintf("%s/%s-%s.json",pact.PactDir,pact.Consumer,pact.Provider)
	
	p.Publish(types.PublishRequest{
		PactURLs:	[]string{pactUrl},
		PactBroker:	pactBrokerUrl,
		ConsumerVersion: consumerVersion,
		Tags:		[]string{current_env},
	})
}
