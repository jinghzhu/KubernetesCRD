# Kubernetes Custom Resource Definition

This repository is an example of how to create/list/update/delete Kubernetes [Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/).



# Environment
1. Go: >= v1.9.0
2. Kubernetes: >= v1.9.0
3. Assume you have already a Kubernetes cluster and its kubeconfig file can be reached via system variable `KUBECONFIG`.



# Dependency Package
The CRD is mainly developed in repository [apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver) which depends on [client-go](https://github.com/kubernetes/client-go), [apimachinery](https://github.com/kubernetes/apimachinery) and [api](https://github.com/kubernetes/api).

This code is based on Kubernetes v1.9.6:

* **k8s.io/client-go** with version `kubernetes-1.9.6`.
* **k8s.io/apimachinery** with version `kubernetes-1.9.6`.
* **k8s.io/apiextensions-apiserver** with version `kubernetes-1.9.6`.
* **k8s.io/api** with version `kubernetes-1.9.6`.



# Step-by-step Instructions
## Define CRD Object
We need firstly create the struct of CRD. The CRD object structure has these components:

* `Metadata`

    Standard Kubernetes properties like name, namespace, labels, etc.

* `Spec`

    CR configuration.

    Each instance of our CR has an attached Spec, which should be defined via a `struct{}` to provide data format validation. In practice, this Spec is arbitrary key-value data that specifies the configuration/behavior of the CR.

* `Status`

    Used by the CR controller in response to Spec updates.

```go
// pkg/crd/jinghzhu/v1/types.go

// Jinghzhu is the CRD. Use this command to generate deepcopy for it:
// ./k8s.io/code-generator/generate-groups.sh all github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu/v1/apis github.com/jinghzhu/KubernetesCRD/pkg/crd "jinghzhu:v1"
type Jinghzhu struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	metav1.ObjectMeta `json:"metadata"`
	// Specification of the desired behavior of Jinghzhu.
	Spec JinghzhuSpec `json:"spec"`
	// Observed status of Jinghzhu.
	Status JinghzhuStatus `json:"status"`
}

// JinghzhuSpec is a desired state description of Jinghzhu.
type JinghzhuSpec struct {
	Foo string `json:"foo"`
	Bar bool   `json:"bar"`
}

// JinghzhuStatus describes the lifecycle status of Jinghzhu.
type JinghzhuStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}
```

We need to leverage the automatically code generated script to create the deep copy methods for CRD object. The resutl can be found at:
* `pkg/crd/jinghzhu/v1/zz_generated.deepcopy.go`
* `pkg/crd/jinghzhu/v1/apis`

You can get **code-generator** from [GitHub](https://github.com/kubernetes/code-generator).

In my example, I run following command to generate that file:
```bash
$ ./k8s.io/code-generator/generate-groups.sh all github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu/v1/apis github.com/jinghzhu/KubernetesCRD/pkg/crd "jinghzhu:v1"
```

> It requires you already have the code-generator at `$GOPATH/src/k8s.io`.

It may also need to create file [boilerplate.go.txt](https://github.com/kubernetes/kubernetes/blob/release-1.8/hack/boilerplate/boilerplate.go.txt). Please note that the code-generator requires annotation to work as expected. If you carefully read my code, you can find the annotations at:

* `pkg/crd/jinghzhu/v1/types.go`:

```go
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=jinghzhu

...

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

...

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=jinghzhu
```

* `pkg/crd/jinghzhu/v1/doc.go`:

```go
// +k8s:deepcopy-gen=package,register
// +k8s:defaulter-gen=TypeMeta
// +k8s:openapi-gen=true

// Package v1 is the v1 version of the API.
// +groupName=jinghzhu.com
```

For more details about annotations for code-generator, please read Kubernetes documents or Google it.

> If you use dep to manage dependency packages, you may not directly run the script, generate-groups. It will throw the error that it can't find necessary vendor packages.
> 
> The solution is to deploy the dependency packages in `$GOPATH`. For more details, please view the [kubernetes/code-generator#21](https://github.com/kubernetes/code-generator/issues/21).


## Register CRD
The CRD name (`CRDName string = Plural + "." + GroupName`) is the combination of CR plural (`Plural string = "jinghzhus"`) and CR group(`GroupName string = "jinghzhu.com"`) which can be used for the reference of Kubernetes CLI or API. CR group and version also define API endpoints.

`pkg/crd/jinghzhu/v1/register.go`:
```go
const (
	// GroupName is the group name used in this package.
	GroupName string = "jinghzhu.com"
	Kind      string = "Jinghzhu"
	// GroupVersion is the version.
	GroupVersion string = "v1"
	// Plural is the Plural for Jinghzhu.
	Plural string = "jinghzhus"
	// Singular is the singular for Jinghzhu.
	Singular string = "jinghzhu"
	// CRDName is the CRD name for Jinghzhu.
	CRDName string = Plural + "." + GroupName
)
```

The method, `CreateCustomResourceDefinition` at `pkg/crd/jinghzhu/v1/crd.go`, covers the logic to register CRD into Kubernetes. Please note that Kubernetes doesn't immediately register CRD. So, it is better to add logic to wait for the creation like shown in the code.


## CRD Client
After creating CRD, we can access via CLI. For easily usage, we hope it can also be accessed via API. So I develop some methods to wrapper some codes for CRD **Create**, **Update**, **Delete**, **Get**, and **List**. You can view them at `pkg/crd/jinghzhu/v1/client/client.go`.



# Main Logic to Use CRD
Here, I'll go through the main code (`main.go`) to show the main logic of everything set before for us to use CRD.

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
    // Init a CRD kind.
	if _, err = crdjinghzhuv1.CreateCustomResourceDefinition("crd-ns", apiextensionsClientSet); err != nil {
		panic(err)
	}
    ```

3. Create the API client to help access the CRD.

    ```go
    // Create a CRD client interface for Jinghzhu v1.
	crdClient, err := jinghzhuv1client.NewClient(kubeConfigPath, types.DefaultCRDNamespace)
	if err != nil {
		panic(err)
	}
    ```

4. Declare a CR object for test.

    ```go
    instanceName := "jinghzhu-example1"
	exampleInstance := &crdjinghzhuv1.Jinghzhu{
		ObjectMeta: metav1.ObjectMeta{
			Name: instanceName,
		},
		Spec: crdjinghzhuv1.JinghzhuSpec{
			Foo: "hello",
			Bar: true,
		},
		Status: crdjinghzhuv1.JinghzhuStatus{
			State:   crdjinghzhuv1.StatePending,
			Message: "Created but not processed yet",
		},
	}
    ```

5. Use the API client created in step 3 to create new CR.

    ```go
    result, err := crdClient.Create(exampleInstance)
	if err == nil {
		fmt.Printf("CREATED: %#v\n", result)
	} else if apierrors.IsAlreadyExists(err) {
		fmt.Printf("ALREADY EXISTS: %#\n", result)
	} else {
		panic(err)
	}

	// Wait until the CRD object is handled by controller and its status is changed to Processed.
	err = crdClient.WaitForInstanceProcessed(instanceName)
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

If everything goes well, you should see logs like:

```bash
$ KUBECONFIG=~/.kube/config go run cmd/crd/main.go
CRD Jinghzhu is created

CREATED: &v1.Jinghzhu{TypeMeta:v1.TypeMeta{Kind:"", APIVersion:""}, ObjectMeta:v1.ObjectMeta{Name:"jinghzhu-example1", GenerateName:"", Namespace:"crd-ns", SelfLink:"/apis/jinghzhu.com/v1/namespaces/crd-ns/jinghzhus/jinghzhu-example1", UID:"2095f6e4-8b76-11e9-92f8-02005ee2a828", ResourceVersion:"40895595", Generation:0, CreationTimestamp:v1.Time{Time:time.Time{wall:0x0, ext:63695764308, loc:(*time.Location)(0x23468a0)}}, DeletionTimestamp:(*v1.Time)(nil), DeletionGracePeriodSeconds:(*int64)(nil), Labels:map[string]string(nil), Annotations:map[string]string(nil), OwnerReferences:[]v1.OwnerReference(nil), Initializers:(*v1.Initializers)(nil), Finalizers:[]string(nil), ClusterName:""}, Spec:v1.JinghzhuSpec{Foo:"hello", Bar:true}, Status:v1.JinghzhuStatus{State:"Pending", Message:"Created but not processed yet"}}

Porcessed

LIST: &v1.JinghzhuList{TypeMeta:v1.TypeMeta{Kind:"", APIVersion:""}, ListMeta:v1.ListMeta{SelfLink:"/apis/jinghzhu.com/v1/namespaces/crd-ns/jinghzhus", ResourceVersion:"40895598", Continue:""}, Items:[]v1.Jinghzhu{v1.Jinghzhu{TypeMeta:v1.TypeMeta{Kind:"Jinghzhu", APIVersion:"jinghzhu.com/v1"}, ObjectMeta:v1.ObjectMeta{Name:"jinghzhu-example1", GenerateName:"", Namespace:"crd-ns", SelfLink:"/apis/jinghzhu.com/v1/namespaces/crd-ns/jinghzhus/jinghzhu-example1", UID:"2095f6e4-8b76-11e9-92f8-02005ee2a828", ResourceVersion:"40895595", Generation:0, CreationTimestamp:v1.Time{Time:time.Time{wall:0x0, ext:63695764308, loc:(*time.Location)(0x23468a0)}}, DeletionTimestamp:(*v1.Time)(nil), DeletionGracePeriodSeconds:(*int64)(nil), Labels:map[string]string(nil), Annotations:map[string]string(nil), OwnerReferences:[]v1.OwnerReference(nil), Initializers:(*v1.Initializers)(nil), Finalizers:[]string(nil), ClusterName:""}, Spec:v1.JinghzhuSpec{Foo:"hello", Bar:true}, Status:v1.JinghzhuStatus{State:"Pending", Message:"Created but not processed yet"}}}}
```



## Result
Now, let's check CRD.

```bash
$ kubectl get crd
NAME                            AGE
jinghzhus.jinghzhu.com          21m
```

```bash
$ kubectl describe crd jinghzhus.jinghzhu.com
Name:         jinghzhus.jinghzhu.com
Namespace:
Labels:       <none>
Annotations:  <none>
API Version:  apiextensions.k8s.io/v1beta1
Kind:         CustomResourceDefinition
Metadata:
  Creation Timestamp:  2019-06-10T11:51:13Z
  Generation:          1
  Resource Version:    40895519
  Self Link:           /apis/apiextensions.k8s.io/v1beta1/customresourcedefinitions/jinghzhus.jinghzhu.com
  UID:                 0bdc06fa-8b76-11e9-92f8-02005ee2a828
Spec:
  Group:  jinghzhu.com
  Names:
    Kind:       Jinghzhu
    List Kind:  JinghzhuList
    Plural:     jinghzhus
    Singular:   jinghzhu
  Scope:        Namespaced
  Version:      v1
Status:
  Accepted Names:
    Kind:       Jinghzhu
    List Kind:  JinghzhuList
    Plural:     jinghzhus
    Singular:   jinghzhu
  Conditions:
    Last Transition Time:  2019-06-10T11:51:13Z
    Message:               no conflicts found
    Reason:                NoConflicts
    Status:                True
    Type:                  NamesAccepted
    Last Transition Time:  2019-06-10T11:51:13Z
    Message:               the initial names have been accepted
    Reason:                InitialNamesAccepted
    Status:                True
    Type:                  Established
Events:                    <none>
```

```bash
$ kubectl proxy
Starting to serve on 127.0.0.1:8001

$ curl -i 127.0.0.1:8001/apis/jinghzhu.com
HTTP/1.1 200 OK
Content-Length: 294
Content-Type: application/json
Date: Mon, 10 Jun 2019 12:13:34 GMT

{
  "kind": "APIGroup",
  "apiVersion": "v1",
  "name": "jinghzhu.com",
  "versions": [
    {
      "groupVersion": "jinghzhu.com/v1",
      "version": "v1"
    }
  ],
  "preferredVersion": {
    "groupVersion": "jinghzhu.com/v1",
    "version": "v1"
  },
  "serverAddressByClientCIDRs": null
}
```

```bash
$ kubectl -n crd-ns get jinghzhus
NAME                AGE
jinghzhu-example1   22m

$ kubectl -n crd-ns describe jinghzhus jinghzhu-example1
Name:         jinghzhu-example1
Namespace:    crd-ns
Labels:       <none>
Annotations:  <none>
API Version:  jinghzhu.com/v1
Kind:         Jinghzhu
Metadata:
  Cluster Name:
  Creation Timestamp:  2019-06-10T11:51:48Z
  Resource Version:    40895595
  Self Link:           /apis/jinghzhu.com/v1/namespaces/crd-ns/jinghzhus/jinghzhu-example1
  UID:                 2095f6e4-8b76-11e9-92f8-02005ee2a828
Spec:
  Bar:  true
  Foo:  hello
Status:
  Message:  Created but not processed yet
  State:    Pending
Events:     <none>
```