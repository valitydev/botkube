package utils

import (
	"sync"

	"github.com/infracloudio/botkube/pkg/log"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/types"
)

// TODO: Use redis or key-value DB for caching to work on scale

// TODO: Make async thread safe
type Resources struct {
	ResourceMap sync.Map
}

func (r *Resources) Load(usObj *unstructured.Unstructured) {
	if usObj == nil || usObj.GetUID() == "" {
		return
	}
	r.ResourceMap.Store(usObj.GetUID(), usObj)
}

func (r *Resources) Get(uID types.UID) (*unstructured.Unstructured, bool) {
	val, ok := r.ResourceMap.Load(uID)
	if !ok {
		return nil, false
	}
	robj, ok := val.(*unstructured.Unstructured)
	if !ok {
		log.Error("Failed to case cache value into runtime object")
		return nil, false
	}
	return robj, true
}

func (r *Resources) Delete(uID types.UID) {
	r.ResourceMap.Delete(uID)
}
