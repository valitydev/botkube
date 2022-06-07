package filters

import (
	"fmt"
	"github.com/infracloudio/botkube/pkg/config"
	"github.com/infracloudio/botkube/pkg/events"
	"github.com/infracloudio/botkube/pkg/filterengine"
	"github.com/infracloudio/botkube/pkg/log"
	"github.com/infracloudio/botkube/pkg/utils"
	coreV1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"reflect"
)

const (
	//Message string representing message template
	Message = "You can see your pod's errors in OpenSearch: "
)

// DeployErrorsChecker checks if some errors occurred during deployment
type DeployErrorsChecker struct {
	Description string
}

// Run filter and generate message
func (d DeployErrorsChecker) Run(object interface{}, event *events.Event) {
	if event.Kind != "Pod" || event.Type != config.ErrorEvent {
		return
	}
	commConfig, confErr := config.NewCommunicationsConfig()
	if confErr != nil {
		log.Errorf("Error in loading configuration. %s", confErr.Error())
		return
	}
	if commConfig == nil {
		log.Errorf("Error in loading configuration.")
		return
	}
	var podObj coreV1.Pod
	err := utils.TransformIntoTypedObject(object.(*unstructured.Unstructured), &podObj)
	if err != nil {
		log.Errorf("Unable to transform object type: %v, into type: %v", reflect.TypeOf(object), reflect.TypeOf(podObj))
	}
	searchURLTemplate := commConfig.Communications.PodLogsDashboard.URL
	event.LogsURLMsg = fmt.Sprintf(Message+"[LOGS URL]("+searchURLTemplate+")", podObj.Name)
}

// Describe filter
func (d DeployErrorsChecker) Describe() string {
	return d.Description
}

// Register filter
func init() {
	filterengine.DefaultFilterEngine.Register(DeployErrorsChecker{
		Description: "Checks if errors occurred while deployment and adds link to kibana for that pod.",
	})
}
