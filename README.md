# Kubernetes Custom Resource Definition

This repository is an example of how to create/list/update/delete Kubernetes Custom Resource Definition.


## Environment
1. Go: >= v1.7.0
2. Kubernetes: v1.8.0 or v1.8.1
3. Assume you have already a Kubernetes cluster and its kubeconfig file can be reached via system variable `KUBECONFIG`.


## Dependency Package
The CRD is mainly developed in repository [apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver) which depends on [client-go](https://github.com/kubernetes/client-go), [apimachinery](https://github.com/kubernetes/apimachinery) and [api](https://github.com/kubernetes/api). Please note that it is very time consuming and headache to make all of them work well for Kubernetes v1.7.x (see [kubernetes/apiextensions-apiserver#3](https://github.com/kubernetes/apiextensions-apiserver/issues/3) and [kubernetes/client-go#247](https://github.com/kubernetes/client-go/issues/247)). So, my code is based on Kubernetes v1.8.1:

* **k8s.io/client-go** with version `v5.0.1`.
* **k8s.io/apimachinery** with version `kubernetes-1.8.1`.
* **k8s.io/apiextensions-apiserver** with version `kubernetes-1.8.1`.
* **k8s.io/api** with commit `fe29995db37613b9c5b2a647544cf627bfa8d299`.


## Step-by-step Instruction
### Define CRD Object
We need firstly create the struct of CRD. The CRD object structure has these components:
* Metadata
    Standard Kubernetes properties like name, namespace, labels, etc.

* Spec
    CR configuration.

    Each instance of our CR has an attached Spec, which should be defined via a `struct{}` to provide data format validation. In practice, this Spec is arbitrary key-value data that specifies the configuration/behavior of the CR.

* Status
    Used by the CR controller in response to Spec updates.

```go
// apis/test/v1/types.go
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Test is the CRD.
type Test struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              TestSpec   `json:"spec"`
	Status            TestStatus `json:"status,omitempty"`
}

// Spec
type TestSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}

// Status
type TestStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}
```

From Kubernetes v1.8.0, we need to use its automatically code generated script to create the deep copy methods for CRD object (`pkg\apis\test\v1\zz_generated.deepcopy.go`)

You can get **code-generator** from [GitHub](https://github.com/kubernetes/code-generator).

In my example, I run following command to generate that file:
```bash
$ ./generate-groups.sh deepcopy github.com/jinghzhu/k8scrd/client github.com/jinghzhu/k8scrd/apis "test:v1"
```

And it may also need to create file [boilerplate.go.txt](https://github.com/kubernetes/kubernetes/blob/release-1.8/hack/boilerplate/boilerplate.go.txt).

Please note that as I use dep to manage dependency packages, we can't directly run the script, generate-groups. It will throw the error that it can't find necessary vendor packages.

The solution is we also need to deploy the dependency packages in `$GOPATH`. For more details, please view the [kubernetes/code-generator#21](https://github.com/kubernetes/code-generator/issues/21).


### Register CRD
The CRD name (`TestCRDName`) is the combination of CR plural (`TestResourcePlural`) and CR group(`GroupName`) which can be used for the reference of Kubernetes CLI or API. CR group and version also define API endpoints.

```go
// apis/test/v1/types.go
const (
	TestResourcePlural string = "tests"
	GroupName        string = "test.io"
	TestCRDName      string = TestResourcePlural + "." + GroupName
	version          string = "v1"
)
```

The following method covers the logic to register CRD into Kubernetes:

```go
// apis/test/v1/crd.go
func CreateCustomResourceDefinition(clientSet apiextensionsclient.Interface) (*apiextensionsv1beta1.CustomResourceDefinition, error) {
	crd := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: TestCRDName,
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   GroupName,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: TestResourcePlural,
				Kind:   reflect.TypeOf(Test{}).Name(),
			},
		},
	}
	_, err := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Create(crd)
	if err != nil {
		fmt.Println("Fail to create CRD: " + err.Error())
		return nil, err
	}

	// Wait for CRD creation.
	err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crd, err = clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Get(TestCRDName, metav1.GetOptions{})
		if err != nil {
			fmt.Println("Fail to wait for CRD creation: " + err.Error())
			return false, err
		}
		for _, cond := range crd.Status.Conditions {
			switch cond.Type {
			case apiextensionsv1beta1.Established:
				if cond.Status == apiextensionsv1beta1.ConditionTrue {
					return true, err
				}
			case apiextensionsv1beta1.NamesAccepted:
				if cond.Status == apiextensionsv1beta1.ConditionFalse {
					fmt.Println(fmt.Sprintf("Name conflict while wait for CRD creation: %v, %v", cond.Reason, err))
				}
			}
		}
		return false, err
	})
	if err != nil {
		deleteErr := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(TestCRDName, nil)
		if deleteErr != nil {
			fmt.Println("Fail to delete CRD: " + deleteErr.Error())
			return nil, errors.NewAggregate([]error{err, deleteErr})
		}
		return nil, err
	}
	return crd, nil
}
```

Please note that Kubernetes doesn't immediately register CRD. So, it is better to add logic to wait for the creation like shown in the code.



### CRD Client
After creating CRD, we can access via CLI. For easily usage, we hope it can also be accessed via API. So I develop some methods to wrapper some codes for CRD **Create**, **Update**, **Delete**, **Get**, and **List**. You can view them at `client/client.go`.

Also, to initalize this client, we need to let it be aware of our CRD schema:

```go
// client/client.go
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
```


### CRD Controller
Now, I have registered CRD and access it via both CLI and API. With them, users can write some yaml files to add CR instances. But if I want to leverage CRD to help us monitor and perform actions on the CR events, I introduce a controller (`controller/controller.go`) to do it.

It can only watch **Add/Update/Delete** event:
```go
// controller/controller.go
func (c *TestController) watch(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.TestClient,
		testv1.TestResourcePlural,
		corev1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		source,
		&testv1.Test{},
		0,
		// CRD event handlers.
		cache.ResourceEventHandlerFuncs{  // <--- watch events
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		},
	)

	go controller.Run(ctx.Done())
	return controller, nil
}
```


## Main Logic to Use CRD
I'll go through the main code (`main.go`) to show the main logic of everything set before for us to use CRD.

1. connect to Kubernetes.

    ```go
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
    ```

2. register the CRD. Please note that it can only be accessed by CLI now as mentioned before.

    ```go
    // Init a CRD.
	crd, err := testv1.CreateCustomResourceDefinition(apiextensionsClientSet)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}
    ```

3. create the API client to help access the CRD.

    ```go
    // Make a new config for extension's API group and use the first one as the baseline.
	testClient, testScheme, err := client.NewClient(clientConfig)
	if err != nil {
		panic(err)
    }

    // Create a CRD client interface.
	crdClient := client.NewCrdClient(testClient, testScheme, testv1.DefaultNamespace)
    ```

4. asynchronous create the CR events controller. You can also do it later.

    ```go
    // Start CRD controller.
	controller := k8scrdcontroller.TestController{
		TestClient: testClient,
		TestScheme: testScheme,
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
    go controller.Run(ctx)
    ```

5. declare a CR object for test.

    ```go
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
    ```

6. use the API client created in step 3 to create new CR.

    ```go
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
	testList, err := crdClient.List(metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
    fmt.Printf("LIST: %#v\n", testList)
    ```


## Result
Now, let's check CRD.

```bash
$ kubectl get crd
NAME            KIND
tests.test.io   CustomResourceDefinition.v1beta1.apiextensions.k8s.io

$ kubectl describe crd tests.test.io
Name:		tests.test.io
Namespace:	
Labels:		<none>
Annotations:	<none>
API Version:	apiextensions.k8s.io/v1beta1
Kind:		CustomResourceDefinition
Metadata:
  Creation Timestamp:	2017-12-13T13:24:25Z
  Resource Version:	1818148
  Self Link:		/apis/apiextensions.k8s.io/v1beta1/customresourcedefinitions/tests.test.io
  UID:			f092e852-e008-11e7-a465-02000455d788
Spec:
  Group:	test.io
  Names:
    Kind:	Test
    List Kind:	TestList
    Plural:	tests
    Singular:	test
  Scope:	Namespaced
  Version:	v1
Status:
  Accepted Names:
    Kind:	Test
    List Kind:	TestList
    Plural:	tests
    Singular:	test
  Conditions:
    Last Transition Time:	<nil>
    Message:			no conflicts found
    Reason:			NoConflicts
    Status:			True
    Type:			NamesAccepted
    Last Transition Time:	2017-12-13T13:24:25Z
    Message:			the initial names have been accepted
    Reason:			InitialNamesAccepted
    Status:			True
    Type:			Established
Events:				<none>
```

```
$ kubectl proxy
Starting to serve on 127.0.0.1:8001

$ curl -i 127.0.0.1:8001/apis/test.io/v1
HTTP/1.1 200 OK
Content-Length: 391
Content-Type: application/json
Date: Wed, 13 Dec 2017 13:27:08 GMT

{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "test.io/v1",
  "resources": [
    {
      "name": "tests",
      "singularName": "test",
      "namespaced": true,
      "kind": "Test",
      "verbs": [
        "delete",
        "deletecollection",
        "get",
        "list",
        "patch",
        "create",
        "update",
        "watch"
      ]
    }
  ]
}
```

