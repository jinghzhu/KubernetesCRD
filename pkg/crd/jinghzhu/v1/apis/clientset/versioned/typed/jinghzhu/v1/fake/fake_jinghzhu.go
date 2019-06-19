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

package fake

import (
	jinghzhu_v1 "github.com/jinghzhu/KubernetesCRD/pkg/crd/jinghzhu/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakeJinghzhus implements JinghzhuInterface
type FakeJinghzhus struct {
	Fake *FakeJinghzhuV1
	ns   string
}

var jinghzhusResource = schema.GroupVersionResource{Group: "jinghzhu.com", Version: "v1", Resource: "jinghzhus"}

var jinghzhusKind = schema.GroupVersionKind{Group: "jinghzhu.com", Version: "v1", Kind: "Jinghzhu"}

// Get takes name of the jinghzhu, and returns the corresponding jinghzhu object, and an error if there is any.
func (c *FakeJinghzhus) Get(name string, options v1.GetOptions) (result *jinghzhu_v1.Jinghzhu, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(jinghzhusResource, c.ns, name), &jinghzhu_v1.Jinghzhu{})

	if obj == nil {
		return nil, err
	}
	return obj.(*jinghzhu_v1.Jinghzhu), err
}

// List takes label and field selectors, and returns the list of Jinghzhus that match those selectors.
func (c *FakeJinghzhus) List(opts v1.ListOptions) (result *jinghzhu_v1.JinghzhuList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(jinghzhusResource, jinghzhusKind, c.ns, opts), &jinghzhu_v1.JinghzhuList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &jinghzhu_v1.JinghzhuList{}
	for _, item := range obj.(*jinghzhu_v1.JinghzhuList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested jinghzhus.
func (c *FakeJinghzhus) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(jinghzhusResource, c.ns, opts))

}

// Create takes the representation of a jinghzhu and creates it.  Returns the server's representation of the jinghzhu, and an error, if there is any.
func (c *FakeJinghzhus) Create(jinghzhu *jinghzhu_v1.Jinghzhu) (result *jinghzhu_v1.Jinghzhu, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(jinghzhusResource, c.ns, jinghzhu), &jinghzhu_v1.Jinghzhu{})

	if obj == nil {
		return nil, err
	}
	return obj.(*jinghzhu_v1.Jinghzhu), err
}

// Update takes the representation of a jinghzhu and updates it. Returns the server's representation of the jinghzhu, and an error, if there is any.
func (c *FakeJinghzhus) Update(jinghzhu *jinghzhu_v1.Jinghzhu) (result *jinghzhu_v1.Jinghzhu, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(jinghzhusResource, c.ns, jinghzhu), &jinghzhu_v1.Jinghzhu{})

	if obj == nil {
		return nil, err
	}
	return obj.(*jinghzhu_v1.Jinghzhu), err
}

// Delete takes name of the jinghzhu and deletes it. Returns an error if one occurs.
func (c *FakeJinghzhus) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(jinghzhusResource, c.ns, name), &jinghzhu_v1.Jinghzhu{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakeJinghzhus) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(jinghzhusResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &jinghzhu_v1.JinghzhuList{})
	return err
}

// Patch applies the patch and returns the patched jinghzhu.
func (c *FakeJinghzhus) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *jinghzhu_v1.Jinghzhu, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(jinghzhusResource, c.ns, name, data, subresources...), &jinghzhu_v1.Jinghzhu{})

	if obj == nil {
		return nil, err
	}
	return obj.(*jinghzhu_v1.Jinghzhu), err
}
