package dashboard

import (
	"log"
	"strings"
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
    log.Printf(strings.Join(dashboardAnalyser.DashFilesList, " "))
    return nil
}