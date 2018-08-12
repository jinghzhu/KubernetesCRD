package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jinghzhu/k8scrd/client"
	"k8s.io/client-go/tools/clientcmd"

	crdexamplev1 "github.com/jinghzhu/k8scrd/apis/example/v1"
	k8scrdcontroller "github.com/jinghzhu/k8scrd/controller"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
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
	crd, err := crdexamplev1.CreateCustomResourceDefinition(apiextensionsClientSet)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}
	// Just for cleanup.
	defer func() {
		fmt.Println("Exit and clean " + crd.Name)
		apiextensionsClientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(crd.Name, nil)
	}()

	// Make a new config for extension's API group and use the first one as the baseline.
	exampleClient, exampleScheme, err := client.NewClient(clientConfig)
	if err != nil {
		panic(err)
	}

	// Start CRD controller.
	controller := k8scrdcontroller.ExampleController{
		ExampleClient: exampleClient,
		ExampleScheme: exampleScheme,
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
	go controller.Run(ctx)

	// Create a CRD client interface.
	crdClient := client.NewCrdClient(exampleClient, exampleScheme, crdexamplev1.DefaultNamespace)
	// Create an instance of CRD.
	instanceName := "example1"
	exampleInstance := &crdexamplev1.Example{
		ObjectMeta: metav1.ObjectMeta{
			Name: instanceName,
		},
		Spec: crdexamplev1.ExampleSpec{
			Foo: "hello",
			Bar: true,
		},
		Status: crdexamplev1.ExampleStatus{
			State:   crdexamplev1.StateCreated,
			Message: "Created but not processed yet",
		},
	}
	result, err := crdClient.Create(exampleInstance)
	if err == nil {
		fmt.Printf("CREATED: %#v", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#v", result)
	} else {
		panic(err)
	}

	// Wait until the CRD object is handled by controller and its status is changed to Processed.
	err = client.WaitForInstanceProcessed(exampleClient, instanceName)
	if err != nil {
		panic(err)
	}
	fmt.Println("Porcessed")

	// Get the list of CRs.
	exampleList, err := crdClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("LIST: %#v\n", exampleList)

	// As there is a cleanup logic before, here it sleeps for a while for example view.
	sleepDuration := 5 * time.Second
	fmt.Printf("Sleep for %s...\n", sleepDuration.String())
	time.Sleep(sleepDuration)
}
