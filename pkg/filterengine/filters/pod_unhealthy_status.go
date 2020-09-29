package filters

import (
	"strings"

	"github.com/infracloudio/botkube/pkg/config"
	"github.com/infracloudio/botkube/pkg/events"
	"github.com/infracloudio/botkube/pkg/filterengine"
)

// PodStatusChecker Checks and skip error event if pod status is Unhealthy
type PodStatusChecker struct {
	Description string
}

// Register filter
func init() {
	filterengine.DefaultFilterEngine.Register(PodStatusChecker{
		Description: "Checks and skip error event if pod status is Unhealthy",
	})
}

// Run filters and modifies event struct
func (f PodStatusChecker) Run(object interface{}, event *events.Event) {
	if event.Kind != "Pod" || event.Type != config.ErrorEvent || event.Reason != "Unhealthy" {
		return
	}

	if strings.HasPrefix(event.Messages[0], "Readiness probe failed:") || strings.HasPrefix(event.Messages[0], "Liveness probe failed:") {
		event.Skip = true
	}
}

// Describe filter
func (f PodStatusChecker) Describe() string {
	return f.Description
}
