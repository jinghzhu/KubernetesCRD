package main

import (
	"context"
	"fmt"
	"os"

	logger "github.com/jinghzhu/GoUtils/logger"
	crd "github.com/jinghzhu/k8scrd/apis/test0/v1"
	"github.com/jinghzhu/k8scrd/client"
	"github.com/jinghzhu/k8scrd/controller"
	corev1 "k8s.io/api/core/v1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := os.Getenv("KUBECONFIG")

	// Use kubeconfig to create client config.
	clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	apiextensionsClientSet, err := apiextensionsclient.NewForConfig(clientConfig)
	if err != nil {
		panic(err)
	}

	// Init a CRD.
	crd, err := crd.CreateCustomResourceDefinition(apiextensionsClientSet)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}

	// Make a new config for extension's API group and use the first one as the baseline.
	testClient, testScheme, err := client.NewClient(clientConfig)
	if err != nil {
		panic(err)
	}

	// Start CRD controller.
	controller := controller.TestController{
		TestClient: testClient,
		TestScheme: testScheme,
	}
	ctx := context.Background()
	go controller.Run(ctx)

	// Create an instance of CRD.
	instanceName := "test1"
	testInstance := crd.Test{
		ObjectMeta: metav1.ObjectMeta{
			Name: instanceName,
		},
		Spec: crd.TestSpec{
			Foo: "hello",
			Bar: true,
		},
		Status: crd.TestStatus{
			State:   crd.StateCreated,
			Message: "Created but not processed yet",
		},
	}
	var result crd.Test
	err = testClient.Post().
		Resource(crd.TestResourcePlural).
		Namespace(corev1.NamespaceDefault).
		Body(testInstance).
		Do().Into(&result)
	if err == nil {
		logger.Info(fmt.Sprintf("CREATED: %#v", result))
	} else if apierrors.IsAlreadyExists(err) {
		logger.Info(fmt.Sprintf("ALREADY EXISTS: %#v", result))
	} else {
		panic(err)
	}

	// Wait until the CRD object is handled by controller and its status is changed to Processed.
	err = client.WaitForInstanceProcessed(testClient, instanceName)
	if err != nil {
		panic(err)
	}
	logger.Info("Porcessed")

	// Get the list of CRs.
	testList := crd.TestList{}
	err = testClient.Get().Resource(crd.TestResourcePlural).Do().Into(&testList)
	if err != nil {
		panic(err)
	}
	logger.Info(fmt.Sprintf("LIST: %#v", testList))
}
