# Kubernetes Custom Resource Definition

This repository is an example of how to create/list/update/delete Kubernetes [Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).


## Environment
1. Go: >= v1.9.0
2. Kubernetes: > v1.9.0
3. Assume you have already a Kubernetes cluster and its kubeconfig file can be reached via system variable `KUBECONFIG`.


## Dependency Package
The CRD is mainly developed in repository [apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver) which depends on [client-go](https://github.com/kubernetes/client-go), [apimachinery](https://github.com/kubernetes/apimachinery) and [api](https://github.com/kubernetes/api).

This code is based on Kubernetes v1.9.6:

* **k8s.io/client-go** with version `kubernetes-1.9.6`.
* **k8s.io/apimachinery** with version `kubernetes-1.9.6`.
* **k8s.io/apiextensions-apiserver** with version `kubernetes-1.9.6`.
* **k8s.io/api** with version `kubernetes-1.9.6`.


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
// pkg/apis/example/v1/types.go
import metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

// Example is the CRD.
type Example struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
	Spec              ExampleSpec   `json:"spec"`
	Status            ExampleStatus `json:"status,omitempty"`
}

// Spec
type ExampleSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}

// Status
type ExampleStatus struct {
	State   string `json:"state,omitempty"`
	Message string `json:"message,omitempty"`
}
```

We need to use its automatically code generated script to create the deep copy methods for CRD object (see `pkg\apis\example\v1\zz_generated.deepcopy.go`)

You can get **code-generator** from [GitHub](https://github.com/kubernetes/code-generator).

In my example, I run following command to generate that file:
```bash
$ ./generate-groups.sh deepcopy github.com/jinghzhu/k8scrd/pkg/client github.com/jinghzhu/k8scrd/pkg/apis "example:v1"
```

And it may also need to create file [boilerplate.go.txt](https://github.com/kubernetes/kubernetes/blob/release-1.8/hack/boilerplate/boilerplate.go.txt).

Please note that as I use dep to manage dependency packages, we can't directly run the script, generate-groups. It will throw the error that it can't find necessary vendor packages.

The solution is we also need to deploy the dependency packages in `$GOPATH`. For more details, please view the [kubernetes/code-generator#21](https://github.com/kubernetes/code-generator/issues/21).


### Register CRD
The CRD name (`ExampleCRDName`) is the combination of CR plural (`ExampleResourcePlural`) and CR group(`GroupName`) which can be used for the reference of Kubernetes CLI or API. CR group and version also define API endpoints.

```go
// pkg/apis/example/v1/types.go
const (
	ExampleResourcePlural string = "examples"
	GroupName        string = "jinghzhu.io"
	ExampleCRDName      string = ExampleResourcePlural + "." + GroupName
	version          string = "v1"
)
```

The following method covers the logic to register CRD into Kubernetes:

```go
// pkg/apis/example/v1/crd.go
func CreateCustomResourceDefinition(clientSet apiextensionsclient.Interface) (*apiextensionsv1beta1.CustomResourceDefinition, error) {
	crd := &apiextensionsv1beta1.CustomResourceDefinition{
		ObjectMeta: metav1.ObjectMeta{
			Name: ExampleCRDName,
		},
		Spec: apiextensionsv1beta1.CustomResourceDefinitionSpec{
			Group:   GroupName,
			Version: SchemeGroupVersion.Version,
			Scope:   apiextensionsv1beta1.NamespaceScoped,
			Names: apiextensionsv1beta1.CustomResourceDefinitionNames{
				Plural: ExampleResourcePlural,
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
		crd, err = clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Get(ExampleCRDName, metav1.GetOptions{})
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
		deleteErr := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(ExampleCRDName, nil)
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
// pkg/client/client.go
func NewClient(cfg *rest.Config) (*rest.RESTClient, *runtime.Scheme, error) {
	scheme := runtime.NewScheme()
	if err := examplev1.AddToScheme(scheme); err != nil {
		return nil, nil, err
	}

	config := *cfg
	config.GroupVersion = &examplev1.SchemeGroupVersion
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
// pkg/controller/controller.go
func (c *ExampleController) watch(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.TestClient,
		examplev1.ExampleResourcePlural,
		corev1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		source,
		&examplev1.Test{},
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

1. Connect to Kubernetes.

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

2. Register the CRD. Please note that it can only be accessed by CLI now as mentioned before.

    ```go
    // Init a CRD.
	crd, err := examplev1.CreateCustomResourceDefinition(apiextensionsClientSet)
	if err != nil && !apierrors.IsAlreadyExists(err) {
		panic(err)
	}
    ```

3. Create the API client to help access the CRD.

    ```go
    // Make a new config for extension's API group and use the first one as the baseline.
	testClient, testScheme, err := client.NewClient(clientConfig)
	if err != nil {
		panic(err)
    }

    // Create a CRD client interface.
	crdClient := client.NewCrdClient(testClient, testScheme, testv1.DefaultNamespace)
    ```

4. Asynchronous create the CR events controller. You can also do it later.

    ```go
    // Start CRD controller.
	controller := k8scrdcontroller.ExampleController{
		ExampleClient: exampleClient,
		ExampleScheme: exampleScheme,
	}
	ctx, cancelFunc := context.WithCancel(context.Background())
	defer cancelFunc()
    go controller.Run(ctx)
    ```

5. Declare a CR object for test.

    ```go
    // Create an instance of CRD.
	instanceName := "example1"
	exampleInstance := &examplev1.Example{
		ObjectMeta: metav1.ObjectMeta{
			Name: instanceName,
		},
		Spec: examplev1.ExampleSpec{
			Foo: "hello",
			Bar: true,
		},
		Status: examplev1.ExampleStatus{
			State:   examplev1.StateCreated,
			Message: "Created but not processed yet",
		},
    }
    ```

6. Use the API client created in step 3 to create new CR.

    ```go
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
    ```


## Result
Now, let's check CRD.

```bash
$ kubectl get crd
NAME                            AGE
examples.jinghzhu.io            20m

$ kubectl describe crd examples.jinghzhu.io
Name:         examples.jinghzhu.io
Namespace:    
Labels:       <none>
Annotations:  <none>
API Version:  apiextensions.k8s.io/v1beta1
Kind:         CustomResourceDefinition
Metadata:
  Creation Timestamp:  2018-08-12T07:37:04Z
  Generation:          1
  Resource Version:    7499713
  Self Link:           /apis/apiextensions.k8s.io/v1beta1/customresourcedefinitions/examples.jinghzhu.io
  UID:                 825b6c34-9e02-11e8-9c65-020053ae2682
Spec:
  Group:  jinghzhu.io
  Names:
    Kind:       Example
    List Kind:  ExampleList
    Plural:     examples
    Singular:   example
  Scope:        Namespaced
  Version:      v1
Status:
  Accepted Names:
    Kind:       Example
    List Kind:  ExampleList
    Plural:     examples
    Singular:   example
  Conditions:
    Last Transition Time:  2018-08-12T07:37:04Z
    Message:               no conflicts found
    Reason:                NoConflicts
    Status:                True
    Type:                  NamesAccepted
    Last Transition Time:  2018-08-12T07:37:04Z
    Message:               the initial names have been accepted
    Reason:                InitialNamesAccepted
    Status:                True
    Type:                  Established
Events:                    <none>
```

```
$ kubectl proxy
Starting to serve on 127.0.0.1:8001

$ curl -i 127.0.0.1:8001/apis/jinghzhu.io/v1
HTTP/1.1 200 OK
Content-Length: 404
Content-Type: application/json
Date: Sun, 12 Aug 2018 07:59:24 GMT

{
  "kind": "APIResourceList",
  "apiVersion": "v1",
  "groupVersion": "jinghzhu.io/v1",
  "resources": [
    {
      "name": "examples",
      "singularName": "example",
      "namespaced": true,
      "kind": "Example",
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

