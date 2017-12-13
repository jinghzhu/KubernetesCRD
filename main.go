package main

import (
	"context"
	"fmt"
	"os"
	"time"

	testv1 "github.com/jinghzhu/k8scrd/apis/test/v1"
	"github.com/jinghzhu/k8scrd/client"
	k8scrdcontroller "github.com/jinghzhu/k8scrd/controller"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeConfigPath := os.Getenv("KUBECONFIG")

	// Use kubeconfig to create client config.
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeConfigPath)
	if err != nil {
		panic(err)
	}

	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	// Init a CRD.
	crd, err := testv1.CreateCustomResourceDefinition(apiextensionsClientSet)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}
	// Just for cleanup.
	defer func() {
		fmt.Println("Exit and clean " + crd.Name)
		apiextensionsClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, nil)
	}()

	// Make a new config for extension's API group and use the first one as the baseline.
	testClient, testScheme, err := client.NewClient(clientConfig)
	if err != nil {
		panic(err)
	}

	// Start CRD controller.
	controller := k8scrdcontroller.TestController{
		TestClient: testClient,
		TestScheme: testScheme,
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go controller.Run(ctx)

	// Create a CRD client interface.
	crdClient := client.NewCrdClient(testClient, testScheme, testv1.DefaultNamespace)
	// Create an instance of CRD.
	instanceName := "test1"
	testInstance := &testv1.Test{
		ObjectMeta: metav1.ObjectMeta{
			Name: instanceName,
		},
		Spec: testv1.TestSpec{
			Foo: "hello",
			Bar: true,
		},
		Status: testv1.TestStatus{
			State:   testv1.StateCreated,
			Message: "Created but not processed yet",
		},
	}
	// var result testv1.Test
	// err = testClient.Post().
	// 	Resource(testv1.TestResourcePlural).
	// 	Namespace(corev1.NamespaceDefault).
	// 	Body(testInstance).
	// 	Do().Into(&result)
	// if err == nil {
	// 	fmt.Printf("CREATED: %#v", result)
	// } else if apierrors.IsAlreadyExists(err) {
	// 	fmt.Printf("ALREADY EXISTS: %#v", result)
	// } else {
	// 	panic(err)
	// }
	result, err := crdClient.Create(testInstance)
	if err == nil {
		fmt.Printf("CREATED: %#v", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v", result)
	} else {
		panic(err)
	}

	// Wait until the CRD object is handled by controller and its status is changed to Processed.
	err = client.WaitForInstanceProcessed(testClient, instanceName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Porcessed")

	// Get the list of CRs.
	// testList := testv1.TestList{}
	// err = testClient.Get().Resource(testv1.TestResourcePlural).Do().Into(&testList)
	// if err != nil {
	// 	panic(err)
	// }
	testList, err := crdClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", testList)

	// As there is a cleanup logic before, here it sleeps for a while for example view.
	sleepDuration := 5 * time.Second
	fmt.Printf("Sleep for %s...\n", sleepDuration.String())
	time.Sleep(sleepDuration)
}
