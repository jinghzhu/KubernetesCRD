/*
Copyright 2019 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	v1 "github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu/v1"
	scheme "github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu/v1/apis/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// JinghzhusGetter has a method to return a JinghzhuInterface.
// A group's client should implement this interface.
type JinghzhusGetter interface {
	Jinghzhus(namespace string) JinghzhuInterface
}

// JinghzhuInterface has methods to work with Jinghzhu resources.
type JinghzhuInterface interface {
	Create(*v1.Jinghzhu) (*v1.Jinghzhu, error)
	Update(*v1.Jinghzhu) (*v1.Jinghzhu, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Jinghzhu, error)
	List(opts meta_v1.ListOptions) (*v1.JinghzhuList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Jinghzhu, err error)
	JinghzhuExpansion
}

// jinghzhus implements JinghzhuInterface
type jinghzhus struct {
	client rest.Interface
	ns     string
}

// newJinghzhus returns a Jinghzhus
func newJinghzhus(c *JinghzhuV1Client, namespace string) *jinghzhus {
	return &jinghzhus{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the jinghzhu, and returns the corresponding jinghzhu object, and an error if there is any.
func (c *jinghzhus) Get(name string, options meta_v1.GetOptions) (result *v1.Jinghzhu, err error) {
	result = &v1.Jinghzhu{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("jinghzhus").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Jinghzhus that match those selectors.
func (c *jinghzhus) List(opts meta_v1.ListOptions) (result *v1.JinghzhuList, err error) {
	result = &v1.JinghzhuList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("jinghzhus").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested jinghzhus.
func (c *jinghzhus) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("jinghzhus").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a jinghzhu and creates it.  Returns the server's representation of the jinghzhu, and an error, if there is any.
func (c *jinghzhus) Create(jinghzhu *v1.Jinghzhu) (result *v1.Jinghzhu, err error) {
	result = &v1.Jinghzhu{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("jinghzhus").
		Body(jinghzhu).
		Do().
		Into(result)
	return
}

// Update takes the representation of a jinghzhu and updates it. Returns the server's representation of the jinghzhu, and an error, if there is any.
func (c *jinghzhus) Update(jinghzhu *v1.Jinghzhu) (result *v1.Jinghzhu, err error) {
	result = &v1.Jinghzhu{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("jinghzhus").
		Name(jinghzhu.Name).
		Body(jinghzhu).
		Do().
		Into(result)
	return
}

// Delete takes name of the jinghzhu and deletes it. Returns an error if one occurs.
func (c *jinghzhus) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("jinghzhus").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *jinghzhus) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("jinghzhus").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched jinghzhu.
func (c *jinghzhus) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Jinghzhu, err error) {
	result = &v1.Jinghzhu{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("jinghzhus").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
