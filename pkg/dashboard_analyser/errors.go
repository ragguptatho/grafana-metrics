package dashboard_analyser

import "fmt"

type ParseError struct {
	dashboardId string
	panelId     string
	message     string
}

func (p ParseError) Error() string {
	return fmt.Sprintf("DashboardId: %s,PanelId: %s,ErrorMessage: %s", p.dashboardId, p.panelId, p.message)
}
