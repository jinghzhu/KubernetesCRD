package controller

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

// TestController is a watch on resource create/update/delete events.
type TestController struct {
	TestClient *rest.RESTClient
	TestScheme *runtime.Scheme
}
