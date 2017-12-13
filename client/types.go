package client

import (
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
)

type CRDClient struct {
	restClient     *rest.RESTClient
	namespace      string
	plural         string
	parameterCodec runtime.ParameterCodec
}
