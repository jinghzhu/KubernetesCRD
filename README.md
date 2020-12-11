# Kubernetes Custom Resource Definition

This repository is an example of how to create/list/update/delete Kubernetes [Custom Resource Definition](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/). If you also want to learn how to develop your own [Operator](https://coreos.com/operators/), please visit:
* [CRD Operator](https://github.com/jinghzhu/KubernetesCRDOperator)
* [Pod Operator](https://github.com/jinghzhu/KubernetesPodOperator)

The CRD example I write here is like the [ReplicaSet](https://kubernetes.io/docs/concepts/workloads/controllers/replicaset/). You can use my CRD with my Operators mentioned before to run following scenarios:
1. Register CRD.
2. The Operator can automatically reconcile current state to desired state. By saying state, it means the number of running Pods.
3. High availability guarantee, which means if you manually delete a Pod via kubectl, the Operator will automatically reconcile it.



# Environment
1. Go: >= v1.9.0
2. Kubernetes: >= v1.18.0
3. Assume you have already a Kubernetes cluster and its kubeconfig file can be reached via system variable `CRD_KUBECONFIG` and it will create a CRD instance in namespace `crd` by default, which you can also modify via system variable `CRD_NAMESPACE`. For more details, please check package `pkg/config`/



# Dependency Package
The CRD is mainly developed in repository [apiextensions-apiserver](https://github.com/kubernetes/apiextensions-apiserver) which depends on [client-go](https://github.com/kubernetes/client-go) and [apimachinery](https://github.com/kubernetes/apimachinery).

This code is based on Kubernetes v1.18.12:

* **k8s.io/client-go** with version `v0.18.12`.
* **k8s.io/apimachinery** with version `v0.18.12`.
* **k8s.io/apiextensions-apiserver** with version `v0.18.12`.



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
	// Desired is the desired Pod number.
	Desired int `json:"desired"`
	// Current is the number of Pod currently running.
	Current int `json:"current"`
	// PodList is the name list of current Pods.
	PodList []string `json:"podList"`
}

// JinghzhuStatus describes the lifecycle status of Jinghzhu.
type JinghzhuStatus struct {
	State   string `json:"state"`
	Message string `json:"message"`
}
```

We need to leverage the automatically code generated script to create the deep copy methods for CRD object. The result can be found at:
* `pkg/crd/jinghzhu/v1/zz_generated.deepcopy.go`
* `pkg/crd/jinghzhu/v1/apis`

You can get **code-generator** from [GitHub](https://github.com/kubernetes/code-generator).

In my example, I run following command to generate that file:
```bash
$ ./k8s.io/code-generator/generate-groups.sh all github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu/v1/apis github.com/jinghzhu/KubernetesCRD/pkg/crd "jinghzhu:v1"
```

Please note that the code-generator requires annotation to work as expected. If you carefully read my code, you can find the annotations at:

* `pkg/crd/jinghzhu/v1/types.go`:

```go
// +genclient
// +genclient:noStatus
// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=jinghzhu

...

// +k8s:deepcopy-gen=true

...

// +k8s:deepcopy-gen=true
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +resource:path=jinghzhu
```

* `pkg/crd/jinghzhu/v1/doc.go`:

```go
// +k8s:deepcopy-gen=package,register
// +k8s:defaulter-gen=TypeMeta
// +k8s:openapi-gen=true

// Package v1 is the v1 version of the API.
// +groupName=jinghzhu.io
```

For more details about annotations for code-generator, please read Kubernetes documents or Google it.

> If you use dep to manage dependency packages, you may not directly run the script, generate-groups. It will throw the error that it can't find necessary vendor packages.
> 
> The solution is to deploy the dependency packages in `$GOPATH`. For more details, please view the [kubernetes/code-generator#21](https://github.com/kubernetes/code-generator/issues/21).


## Register CRD
The CRD name (`CRDName string = Plural + "." + GroupName`) is the combination of CR plural (`Plural string = "jinghzhus"`) and CR group(`GroupName string = "jinghzhu.io"`) which can be used for the reference of Kubernetes CLI or API. CR group and version also define API endpoints.

`pkg/crd/jinghzhu/v1/register.go`:

```go
const (
	// Kind is normally the CamelCased singular type. The resource manifest uses this.
	Kind string = "Jinghzhu"
	// GroupVersion is the version.
	GroupVersion string = "v1"
	// Plural is the plural name used in /apis/<group>/<version>/<plural>
	Plural string = "jinghzhus"
	// Singular is used as an alias on kubectl for display.
	Singular string = "jinghzhu"
	// CRDName is the CRD name for Jinghzhu.
	CRDName string = Plural + "." + crdjinghzhu.GroupName
	// ShortName is the short alias for the CRD.
	ShortName string = "jh"
)
```

The method, `CreateCustomResourceDefinition` at `pkg/crd/jinghzhu/v1/crd.go`, covers the logic to register CRD into Kubernetes. Please note that Kubernetes doesn't immediately register CRD. So, it is better to add logic to wait for the creation like shown in the code.

Meanwhile, I also leverage OpenAPI v3 to perform validation check. You can find related logic at `CreateCustomResourceDefinition@pkg/crd/jinghzhu/v1`:

```go
      Validation: &apiextensionsv1beta1.CustomResourceValidation{
				OpenAPIV3Schema: &apiextensionsv1beta1.JSONSchemaProps{
					Type: "object",
					Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
						"spec": {
							Type: "object",
							Properties: map[string]apiextensionsv1beta1.JSONSchemaProps{
								"desired": {Type: "integer", Format: "int"},
								"current": {Type: "integer", Format: "int"},
								"podList": {
									Type: "array",
									Items: &apiextensionsv1beta1.JSONSchemaPropsOrArray{
										Schema: &apiextensionsv1beta1.JSONSchemaProps{Type: "string"},
									},
								},
							},
							Required: []string{"desired"},
						},
					},
				},
			},
```


## CRD Client
After creating CRD, we can access via CLI. For easily usage, we hope it can also be accessed via API. So I develop some methods to wrapper some codes for CRD **Create**, **Update**, **Delete**, **Get**, and **List**. You can view them at `pkg/crd/jinghzhu/v1/client/client.go`.



# Main Logic to Use CRD
Here, I'll go through the main code (`main.go`) to show the main logic of everything set before for us to use CRD.

1. Connect to Kubernetes.

    ```go
    kubeConfigPath := os.Getenv("KUBECONFIG")

	  // Use kubeconfig to create client config.
	  clientConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
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
	  if _, err = crdjinghzhuv1.CreateCustomResourceDefinition(apiextensionsClientSet); err != nil {
		  panic(err)
	  }
    ```

3. Create the API client to help access the CRD.

    ```go
    // Create a CRD client interface for Jinghzhu v1.
	  crdClient, err := jinghzhuv1client.NewClient(ctx, kubeconfigPath, cfg.GetCRDNamespace())
	  if err != nil {
		  panic(err)
	  }
    ```

4. Declare a CR object for test.

    ```go
    instanceName := "jinghzhu-example-"
	  exampleInstance := &crdjinghzhuv1.Jinghzhu{
		  ObjectMeta: metav1.ObjectMeta{
			  GenerateName: instanceName,
		  },
		  Spec: crdjinghzhuv1.JinghzhuSpec{
			  Desired: 1,
			  Current: 0,
			  PodList: make([]string, 0),
		  },
		  Status: crdjinghzhuv1.JinghzhuStatus{
			  State:   types.StatePending,
			  Message: "Created but not processed yet",
		  },
	  }
    ```

5. Use the API client created in step 3 to create new CR.

    ```go
    result, err := crdClient.CreateDefault(exampleInstance)
	  if err != nil && apierrors.IsAlreadyExists(err) {
		  fmt.Printf("ALREADY EXISTS: %#v\n", result)
	  } else if err != nil {
		  panic(err)
	  }
	  crdInstanceName := result.GetName()
	  fmt.Println("CREATED: " + result.String())

	  // Wait until the CRD object is handled by controller and its status is changed to Processed.
	  err = crdClient.WaitForInstanceProcessed(crdInstanceName)
	  if err != nil {
		  panic(err)
	  }
	  fmt.Println("Processed " + crdInstanceName)

	  // Get the list of CRs.
	  exampleList, err := crdClient.List(metav1.ListOptions{})
	  if err != nil {
		  panic(err)
	  }
	  fmt.Printf("LIST: %#v\n", exampleList)
    ```



## Result
If everything goes well, you should see logs like:

```bash
$ go run cmd/crd/main.go
CRD Jinghzhu is created
CREATED:        Name = jinghzhu-example-f7wgv
        Resource Version = 2343534
        Desired = 1
        Current = 0
        PodList =
        State = Pending
        Message = Created but not processed yet

Processed jinghzhu-example-f7wgv
LIST: &v1.JinghzhuList{TypeMeta:v1.TypeMeta{Kind:"", APIVersion:""}, ListMeta:v1.ListMeta{SelfLink:"/apis/jinghzhu.io/v1/namespaces/crd/jinghzhus", ResourceVersion:"2343539", Continue:"", RemainingItemCount:(*int64)(nil)}, Items:[]v1.Jinghzhu{v1.Jinghzhu{TypeMeta:v1.TypeMeta{Kind:"Jinghzhu", APIVersion:"jinghzhu.io/v1"}, ObjectMeta:v1.ObjectMeta{Name:"jinghzhu-example-f7wgv", GenerateName:"jinghzhu-example-", Namespace:"crd", SelfLink:"/apis/jinghzhu.io/v1/namespaces/crd/jinghzhus/jinghzhu-example-f7wgv", UID:"c9a14049-e255-45be-99f7-033a5ea87787", ResourceVersion:"2343534", Generation:1, CreationTimestamp:v1.Time{Time:time.Time{wall:0x0, ext:63743288791, loc:(*time.Location)(0x1c802e0)}}, DeletionTimestamp:(*v1.Time)(nil), DeletionGracePeriodSeconds:(*int64)(nil), Labels:map[string]string(nil), Annotations:map[string]string(nil), OwnerReferences:[]v1.OwnerReference(nil), Finalizers:[]string(nil), ClusterName:"", ManagedFields:[]v1.ManagedFieldsEntry{v1.ManagedFieldsEntry{Manager:"main", Operation:"Update", APIVersion:"jinghzhu.io/v1", Time:(*v1.Time)(0xc00039d1c0), FieldsType:"FieldsV1", FieldsV1:(*v1.FieldsV1)(0xc00039d1a0)}}}, Spec:v1.JinghzhuSpec{Desired:1, Current:0, PodList:[]string{}}, Status:v1.JinghzhuStatus{State:"Pending", Message:"Created but not processed yet"}}}}
```

Now, let's check CRD.

```bash
$ kubectl get crd
NAME                    CREATED AT
jinghzhus.jinghzhu.io   2020-12-11T13:06:25Z
```

```bash
$ kubectl describe crd jinghzhus.jinghzhu.com
Name:         jinghzhus.jinghzhu.io
Namespace:
Labels:       <none>
Annotations:  <none>
API Version:  apiextensions.k8s.io/v1
Kind:         CustomResourceDefinition
Metadata:
  Creation Timestamp:  2020-12-11T13:06:25Z
  Generation:          1
  Managed Fields:
    API Version:  apiextensions.k8s.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:status:
        f:acceptedNames:
          f:kind:
          f:listKind:
          f:plural:
          f:shortNames:
          f:singular:
        f:conditions:
    Manager:      kube-apiserver
    Operation:    Update
    Time:         2020-12-11T13:06:25Z
    API Version:  apiextensions.k8s.io/v1beta1
    Fields Type:  FieldsV1
    fieldsV1:
      f:spec:
        f:conversion:
          .:
          f:strategy:
        f:group:
        f:names:
          f:kind:
          f:listKind:
          f:plural:
          f:shortNames:
          f:singular:
        f:preserveUnknownFields:
        f:scope:
        f:validation:
          .:
          f:openAPIV3Schema:
            .:
            f:properties:
              .:
              f:spec:
                .:
                f:properties:
                  .:
                  f:current:
                    .:
                    f:format:
                    f:type:
                  f:desired:
                    .:
                    f:format:
                    f:type:
                  f:podList:
                    .:
                    f:items:
                    f:type:
                f:required:
                f:type:
            f:type:
        f:version:
        f:versions:
      f:status:
        f:storedVersions:
    Manager:         main
    Operation:       Update
    Time:            2020-12-11T13:06:25Z
  Resource Version:  2343519
  Self Link:         /apis/apiextensions.k8s.io/v1/customresourcedefinitions/jinghzhus.jinghzhu.io
  UID:               1afe3057-4398-45f7-9200-90b879d52f6a
Spec:
  Conversion:
    Strategy:  None
  Group:       jinghzhu.io
  Names:
    Kind:       Jinghzhu
    List Kind:  JinghzhuList
    Plural:     jinghzhus
    Short Names:
      jh
    Singular:               jinghzhu
  Preserve Unknown Fields:  true
  Scope:                    Namespaced
  Versions:
    Name:  v1
    Schema:
      openAPIV3Schema:
        Properties:
          Spec:
            Properties:
              Current:
                Format:  int
                Type:    integer
              Desired:
                Format:  int
                Type:    integer
              Pod List:
                Items:
                  Type:  string
                Type:    array
            Required:
              desired
            Type:  object
        Type:      object
    Served:        true
    Storage:       true
Status:
  Accepted Names:
    Kind:       Jinghzhu
    List Kind:  JinghzhuList
    Plural:     jinghzhus
    Short Names:
      jh
    Singular:  jinghzhu
  Conditions:
    Last Transition Time:  2020-12-11T13:06:25Z
    Message:               no conflicts found
    Reason:                NoConflicts
    Status:                True
    Type:                  NamesAccepted
    Last Transition Time:  2020-12-11T13:06:25Z
    Message:               the initial names have been accepted
    Reason:                InitialNamesAccepted
    Status:                True
    Type:                  Established
  Stored Versions:
    v1
Events:  <none>
```

```bash
$ kubectl proxy
Starting to serve on 127.0.0.1:8001

$ curl -i 127.0.0.1:8001/apis/jinghzhu.io
HTTP/1.1 200 OK
Cache-Control: no-cache, private
Content-Length: 253
Content-Type: application/json
Date: Fri, 11 Dec 2020 13:11:34 GMT

{
  "kind": "APIGroup",
  "apiVersion": "v1",
  "name": "jinghzhu.io",
  "versions": [
    {
      "groupVersion": "jinghzhu.io/v1",
      "version": "v1"
    }
  ],
  "preferredVersion": {
    "groupVersion": "jinghzhu.io/v1",
    "version": "v1"
  }
}
```

```bash
$ kubectl -n crd get jh
NAME                     AGE
jinghzhu-example-f7wgv   6m20s

$  kubectl -n crd describe jh jinghzhu-example-f7wgv
Name:         jinghzhu-example-f7wgv
Namespace:    crd
Labels:       <none>
Annotations:  <none>
API Version:  jinghzhu.io/v1
Kind:         Jinghzhu
Metadata:
  Creation Timestamp:  2020-12-11T13:06:31Z
  Generate Name:       jinghzhu-example-
  Generation:          1
  Managed Fields:
    API Version:  jinghzhu.io/v1
    Fields Type:  FieldsV1
    fieldsV1:
      f:metadata:
        f:generateName:
      f:spec:
        .:
        f:current:
        f:desired:
        f:podList:
      f:status:
        .:
        f:message:
        f:state:
    Manager:         main
    Operation:       Update
    Time:            2020-12-11T13:06:31Z
  Resource Version:  2343534
  Self Link:         /apis/jinghzhu.io/v1/namespaces/crd/jinghzhus/jinghzhu-example-f7wgv
  UID:               c9a14049-e255-45be-99f7-033a5ea87787
Spec:
  Current:  0
  Desired:  1
  Pod List:
Status:
  Message:  Created but not processed yet
  State:    Pending
Events:     <none>
```
