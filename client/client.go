package client

import (
	"time"

	"github.com/jinghzhu/GoUtils/logger"
	crd "github.com/jinghzhu/k8scrd/apis/test0.io/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
)

func NewClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := crd.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	config := *cfg
	config.GroupVersion = &crd.SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		logger.Error("Fail to generate REST client: " + err.Error())
		return nil, nil, err
	}

	return client, scheme, nil
}

func WaitForInstanceProcessed(testClient *rest.RESTClient, name string) error {
	return wait.Poll(100*time.Millisecond, 10*time.Second, func() (bool, error) {
		var instance crd.Test
		err := testClient.Get().
			Resource(crd.TestResourcePlural).
			Namespace(corev1.NamespaceDefault).
			Name(name).
			Do().Into(&instance)

		if err == nil && instance.Status.State == crd.StateProcessed {
			return true, nil
		}

		return false, err
	})
}
