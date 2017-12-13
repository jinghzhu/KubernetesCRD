package client

import (
	"fmt"
	"time"

	testv1 "github.com/jinghzhu/k8scrd/apis/test/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
)

func NewClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := testv1.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	config := *cfg
	config.GroupVersion = &testv1.SchemeGroupVersion
	config.APIPath = "/apis"
	config.ContentType = runtime.ContentTypeJSON
	config.NegotiatedSerializer = serializer.DirectCodecFactory{CodecFactory: serializer.NewCodecFactory(scheme)}

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		fmt.Println("Fail to generate REST client: " + err.Error())
		return nil, nil, err
	}

	return client, scheme, nil
}

func WaitForInstanceProcessed(testClient *rest.RESTClient, name string) error {
	fmt.Println("Wait for CRD instance processed...")
	return wait.Poll(100*time.Millisecond, 20*time.Second, func() (bool, error) {
		var instance testv1.Test
		err := testClient.Get().
			Resource(testv1.TestResourcePlural).
			Namespace(corev1.NamespaceDefault).
			Name(name).
			Do().Into(&instance)

		if err == nil && instance.Status.State == testv1.StateProcessed {
			return true, nil
		}

		return false, err
	})
}
