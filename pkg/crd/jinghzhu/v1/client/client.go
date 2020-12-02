package client

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	jinghzhuv1 "github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachinerytypes "k8s.io/apimachinery/pkg/types"

	"k8s.io/apimachinery/pkg/util/wait"
)

// WaitForInstanceProcessed is used for monitor the creation of a CRD instance.
func (c *Client) WaitForInstanceProcessed(ctx context.Context, name string) error {
	return wait.Poll(time.Second, 3*time.Second, func() (bool, error) {
		instance, err := c.Get(ctx, name, metav1.GetOptions{})
		if err == nil && instance.Status.State == jinghzhuv1.StatePending {
			return true, nil
		}
		fmt.Printf("Fail to wait for CRD instance processed: %+v\n", err)

		return false, err
	})
}

// Create post an instance of CRD into Kubernetes.
func (c *Client) Create(ctx context.Context, obj *jinghzhuv1.Jinghzhu, opts metav1.CreateOptions) (*jinghzhuv1.Jinghzhu, error) {
	return c.clientset.JinghzhuV1().Jinghzhus(c.namespace).Create(ctx, obj, opts)
}

// Update puts new instance of CRD to replace the old one.
func (c *Client) Update(ctx context.Context, obj *jinghzhuv1.Jinghzhu, opts metav1.UpdateOptions) (*jinghzhuv1.Jinghzhu, error) {
	return c.clientset.JinghzhuV1().Jinghzhus(c.namespace).Update(ctx, obj, opts)
}

// UpdateSpecAndStatus updates the spec and status filed of Jinghzhu v1.
// If only want to update some sub-resource, please use Patch instead.
func (c *Client) UpdateSpecAndStatus(ctx context.Context, name string, jinghzhuSpec *jinghzhuv1.JinghzhuSpec, jinghzhuStatus *jinghzhuv1.JinghzhuStatus) (*jinghzhuv1.Jinghzhu, error) {
	instance, err := c.Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}

	instance.Spec = *jinghzhuSpec
	instance.Status = *jinghzhuStatus

	return c.Update(ctx, instance, metav1.UpdateOptions{})
}

// Patch applies the patch and returns the patched Jinghzhu v1 instance.
func (c *Client) Patch(ctx context.Context, name string, pt apimachinerytypes.PatchType, data []byte, subresources ...string) (*jinghzhuv1.Jinghzhu, error) {
	var result jinghzhuv1.Jinghzhu
	err := c.clientset.RESTClient().Patch(pt).
		Namespace(c.namespace).
		Resource(c.plural).
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do(ctx).
		Into(&result)

	return &result, err
}

// PatchJSONType uses JSON Type (RFC6902) in PATCH.
func (c *Client) PatchJSONType(ctx context.Context, name string, ops []PatchJSONTypeOps) (*jinghzhuv1.Jinghzhu, error) {
	patchBytes, err := json.Marshal(ops)
	if err != nil {
		return nil, err
	}

	return c.clientset.JinghzhuV1().Jinghzhus(c.namespace).Patch(ctx, name, apimachinerytypes.JSONPatchType, patchBytes, metav1.PatchOptions{})
}

// PatchSpec only updates the spec field of Jinghzhu v1, which is /spec.
func (c *Client) PatchSpec(ctx context.Context, name string, jinghzhuSpec *jinghzhuv1.JinghzhuSpec) (*jinghzhuv1.Jinghzhu, error) {
	ops := make([]PatchJSONTypeOps, 1, 1)
	ops[0].Op = PatchJSONTypeReplace
	ops[0].Path = "/spec"
	ops[0].Value = jinghzhuSpec

	return c.PatchJSONType(ctx, name, ops)
}

// PatchStatus only updates the status field of Jinghzhu v1, which is /status.
func (c *Client) PatchStatus(ctx context.Context, name string, jinghzhuStatus *jinghzhuv1.JinghzhuStatus) (*jinghzhuv1.Jinghzhu, error) {
	ops := make([]PatchJSONTypeOps, 1, 1)
	ops[0].Op = PatchJSONTypeReplace
	ops[0].Path = "/status"
	ops[0].Value = jinghzhuStatus

	return c.PatchJSONType(ctx, name, ops)
}

// PatchSpecAndStatus performs patch for both spec and status field of Jinghzhu.
func (c *Client) PatchSpecAndStatus(
	ctx context.Context,
	name string,
	jinghzhuSpec *jinghzhuv1.JinghzhuSpec,
	jinghzhuStatus *jinghzhuv1.JinghzhuStatus,
) (*jinghzhuv1.Jinghzhu, error) {
	ops := make([]PatchJSONTypeOps, 2, 2)
	ops[0].Op = PatchJSONTypeReplace
	ops[0].Path = "/spec"
	ops[0].Value = jinghzhuSpec
	ops[1].Op = PatchJSONTypeReplace
	ops[1].Path = "/status"
	ops[1].Value = jinghzhuStatus

	return c.PatchJSONType(ctx, name, ops)
}

// Delete removes the CRD instance by given name and delete options.
func (c *Client) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.clientset.JinghzhuV1().Jinghzhus(c.namespace).Delete(ctx, name, opts)
}

// Get returns a pointer to the CRD instance.
func (c *Client) Get(ctx context.Context, name string, opts metav1.GetOptions) (*jinghzhuv1.Jinghzhu, error) {
	return c.clientset.JinghzhuV1().Jinghzhus(c.namespace).Get(ctx, name, opts)
}

// GetWithoutOps retrieves the Jinghzhu instance without any GetOptions.
func (c *Client) GetWithoutOps(ctx context.Context, name string) (*jinghzhuv1.Jinghzhu, error) {
	return c.clientset.JinghzhuV1().Jinghzhus(c.namespace).Get(ctx, name, metav1.GetOptions{})
}

// List returns a list of CRD instances by given list options.
func (c *Client) List(ctx context.Context, opts metav1.ListOptions) (*jinghzhuv1.JinghzhuList, error) {
	return c.clientset.JinghzhuV1().Jinghzhus(c.namespace).List(ctx, opts)
}
