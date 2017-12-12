package v1

import (
	"fmt"
	"reflect"
	"time"

	logger "github.com/jinghzhu/GoUtils/logger"
	apiextensionsv1beta1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1beta1"
	apiextensionsclient "k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/wait"
)

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
		logger.Error("Fail to create CRD: " + err.Error())
		return nil, err
	}

	// Wait for CRD creation.
	err = wait.Poll(500*time.Millisecond, 60*time.Second, func() (bool, error) {
		crd, err = clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Get(TestCRDName, metav1.GetOptions{})
		if err != nil {
			logger.Error("Fail to wait for CRD creation: " + err.Error())
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
					logger.Error(fmt.Sprintf("Name conflict while wait for CRD creation: %v, %v", cond.Reason, err))
				}
			}
		}
		return false, err
	})
	if err != nil {
		deleteErr := clientSet.ApiextensionsV1beta1().CustomResourceDefinitions().Delete(TestCRDName, nil)
		if deleteErr != nil {
			logger.Error("Fail to delete CRD: " + deleteErr.Error())
			return nil, errors.NewAggregate([]error{err, deleteErr})
		}
		return nil, err
	}
	return crd, nil
}

// addKnownTypes adds the set of types defined in this package to the supplied scheme.
func addKnownTypes(scheme *runtime.Scheme) error {
	scheme.AddKnownTypes(SchemeGroupVersion,
		&Test{},
		&TestList{},
	)
	metav1.AddToGroupVersion(scheme, SchemeGroupVersion)

	return nil
}
