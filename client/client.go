package client

import (
	"fmt"
	"time"

	testv1 "github.com/jinghzhu/k8scrd/apis/test/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
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

// NewCrdClient returns a CRD client interface.
func NewCrdClient(client *rest.RESTClient, scheme *runtime.Scheme, namespace string) *CRDClient {
	return &CRDClient{
		restClient:     client,
		namespace:      namespace,
		plural:         testv1.TestResourcePlural,
		parameterCodec: runtime.NewParameterCodec(scheme),
	}
}

func (f *CRDClient) Create(obj *testv1.Test) (*testv1.Test, error) {
	var result testv1.Test
	err := f.restClient.Post().
		Namespace(f.namespace).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *CRDClient) Update(obj *testv1.Test) (*testv1.Test, error) {
	var result testv1.Test
	err := f.restClient.Put().
		Namespace(f.namespace).Resource(f.plural).
		Body(obj).Do().Into(&result)
	return &result, err
}

func (f *CRDClient) Delete(name string, options *metav1.DeleteOptions) error {
	return f.restClient.Delete().
		Namespace(f.namespace).Resource(f.plural).
		Name(name).Body(options).Do().
		Error()
}

func (f *CRDClient) Get(name string) (*testv1.Test, error) {
	var result testv1.Test
	err := f.restClient.Get().
		Namespace(f.namespace).Resource(f.plural).
		Name(name).Do().Into(&result)
	return &result, err
}

func (f *CRDClient) List(opts metav1.ListOptions) (*testv1.TestList, error) {
	var result testv1.TestList
	err := f.restClient.Get().
		Namespace(f.namespace).Resource(f.plural).
		VersionedParams(&opts, f.parameterCodec).
		Do().Into(&result)
	return &result, err
}

// Create a new List watch for our TPR
func (f *CRDClient) NewListWatch() *cache.ListWatch {
	return cache.NewListWatchFromClient(f.restClient, f.plural, f.namespace, fields.Everything())
}
