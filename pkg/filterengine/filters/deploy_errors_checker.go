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
	Message           = "You can see your pod's errors in kibana:"
	KibanaUrlTemplate = "https://kibana.empayre.com/app/discover#/?_g=(filters:!(),refreshInterval:(pause:!t,value:0),time:(from:now-60m,to:now))&_a=(columns:!('@severity',message,kubernetes.pod.name),filters:!(),index:filebeat-rbkmoney-processing,interval:auto,query:(language:kuery,query:'kubernetes.pod.name:%20%22%s%22'),sort:!())"
)

type DeployErrorsChecker struct {
	Description string
}

func (d DeployErrorsChecker) Run(object interface{}, event *events.Event) {
	if event.Kind != "Pod" || event.Type != config.ErrorEvent {
		return
	}
	var podObj coreV1.Pod
	err := utils.TransformIntoTypedObject(object.(*unstructured.Unstructured), &podObj)
	if err != nil {
		log.Errorf("Unable to transform object type: %v, into type: %v", reflect.TypeOf(object), reflect.TypeOf(podObj))
	}
	event.Recommendations = append(event.Recommendations, fmt.Sprintf(Message+KibanaUrlTemplate, podObj.Name))
}

func (d DeployErrorsChecker) Describe() string {
	return d.Description
}

// Register filter
func init() {
	filterengine.DefaultFilterEngine.Register(DeployErrorsChecker{
		Description: "Checks if errors occurred while deployment and adds link to kibana for that pod.",
	})
}
