package controller

import (
	"k8s.io/client-go/rest"
	"k8s.io/apimachinery/pkg/runtime"
)

// TestController is a watch on resource create/update/delete events.
type TestController struct {
	TestClient *rest.RESTClient
	TestScheme *runtime.Scheme
}