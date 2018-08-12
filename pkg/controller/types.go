package controller

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// ExampleController is a watch on resource create/update/delete events.
type ExampleController struct {
	ExampleClient *rest.RESTClient
	ExampleScheme *runtime.Scheme
}
