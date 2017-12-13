package controller

import (
	"context"
	"fmt"

	crd "github.com/jinghzhu/k8scrd/apis/test/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/tools/cache"
)

// Run starts a CRD resource controller.
func (c *TestController) Run(ctx context.Context) error {
	fmt.Println("Watch CRD objects...")

	// Watch CRD objects
	_, err := c.watch(ctx)
	if err != nil {
		fmt.Printf("Failed to register watch for Example resource: %v\n", err)
		return err
	}

	<-ctx.Done()
	return ctx.Err()
}

func (c *TestController) watch(ctx context.Context) (cache.Controller, error) {
	source := cache.NewListWatchFromClient(
		c.TestClient,
		crd.TestResourcePlural,
		corev1.NamespaceAll,
		fields.Everything(),
	)

	_, controller := cache.NewInformer(
		source,
		&crd.Test{},
		// Every resyncPeriod, all resources in the cache will retrigger events.
		// Set to 0 to disable the resync.
		0,
		// CRD event handlers.
		cache.ResourceEventHandlerFuncs{
			AddFunc:    c.onAdd,
			UpdateFunc: c.onUpdate,
			DeleteFunc: c.onDelete,
		},
	)

	go controller.Run(ctx.Done())
	return controller, nil
}

func (c *TestController) onAdd(obj interface{}) {
	test := obj.(*crd.Test)
	fmt.Println("[CONTROLLER] OnAdd " + test.ObjectMeta.SelfLink)

	// Use DeepCopy() to make a deep copy of original object and modify this copy
	// or create a copy manually for better performance.
	testCopy := test.DeepCopy()
	testCopy.Status = crd.TestStatus{
		State:   crd.StateProcessed,
		Message: "Successfully processed by controller",
	}

	err := c.TestClient.Put().
		Name(test.ObjectMeta.Name).
		Namespace(test.ObjectMeta.Namespace).
		Resource(crd.TestResourcePlural).
		Body(testCopy).
		Do().
		Error()

	if err != nil {
		fmt.Println("ERROR updating status: " + err.Error())
	} else {
		fmt.Println("UPDATED status: " + testCopy.SelfLink)
	}
}

func (c *TestController) onUpdate(oldObj, newObj interface{}) {
	old := oldObj.(*crd.Test)
	new := newObj.(*crd.Test)
	fmt.Println("[CONTROLLER] OnUpdate old: " + old.ObjectMeta.SelfLink)
	fmt.Println("[CONTROLLER] OnUpdate new: " + new.ObjectMeta.SelfLink)
}

func (c *TestController) onDelete(obj interface{}) {
	test := obj.(*crd.Test)
	fmt.Println("[CONTROLLER] OnDelete " + test.ObjectMeta.SelfLink)
}
