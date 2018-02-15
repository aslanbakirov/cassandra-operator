/*
Copyright 2018 The Kubernetes Authors.

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
	aslangroup_io_v1 "github.com/aslanbekirov/personcrd/pkg/apis/aslangroup.io/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	labels "k8s.io/apimachinery/pkg/labels"
	schema "k8s.io/apimachinery/pkg/runtime/schema"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	testing "k8s.io/client-go/testing"
)

// FakePersons implements PersonInterface
type FakePersons struct {
	Fake *FakeAslangroupV1
	ns   string
}

var personsResource = schema.GroupVersionResource{Group: "aslangroup.io", Version: "v1", Resource: "persons"}

var personsKind = schema.GroupVersionKind{Group: "aslangroup.io", Version: "v1", Kind: "Person"}

// Get takes name of the person, and returns the corresponding person object, and an error if there is any.
func (c *FakePersons) Get(name string, options v1.GetOptions) (result *aslangroup_io_v1.Person, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewGetAction(personsResource, c.ns, name), &aslangroup_io_v1.Person{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aslangroup_io_v1.Person), err
}

// List takes label and field selectors, and returns the list of Persons that match those selectors.
func (c *FakePersons) List(opts v1.ListOptions) (result *aslangroup_io_v1.PersonList, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewListAction(personsResource, personsKind, c.ns, opts), &aslangroup_io_v1.PersonList{})

	if obj == nil {
		return nil, err
	}

	label, _, _ := testing.ExtractFromListOptions(opts)
	if label == nil {
		label = labels.Everything()
	}
	list := &aslangroup_io_v1.PersonList{}
	for _, item := range obj.(*aslangroup_io_v1.PersonList).Items {
		if label.Matches(labels.Set(item.Labels)) {
			list.Items = append(list.Items, item)
		}
	}
	return list, err
}

// Watch returns a watch.Interface that watches the requested persons.
func (c *FakePersons) Watch(opts v1.ListOptions) (watch.Interface, error) {
	return c.Fake.
		InvokesWatch(testing.NewWatchAction(personsResource, c.ns, opts))

}

// Create takes the representation of a person and creates it.  Returns the server's representation of the person, and an error, if there is any.
func (c *FakePersons) Create(person *aslangroup_io_v1.Person) (result *aslangroup_io_v1.Person, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewCreateAction(personsResource, c.ns, person), &aslangroup_io_v1.Person{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aslangroup_io_v1.Person), err
}

// Update takes the representation of a person and updates it. Returns the server's representation of the person, and an error, if there is any.
func (c *FakePersons) Update(person *aslangroup_io_v1.Person) (result *aslangroup_io_v1.Person, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewUpdateAction(personsResource, c.ns, person), &aslangroup_io_v1.Person{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aslangroup_io_v1.Person), err
}

// Delete takes name of the person and deletes it. Returns an error if one occurs.
func (c *FakePersons) Delete(name string, options *v1.DeleteOptions) error {
	_, err := c.Fake.
		Invokes(testing.NewDeleteAction(personsResource, c.ns, name), &aslangroup_io_v1.Person{})

	return err
}

// DeleteCollection deletes a collection of objects.
func (c *FakePersons) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	action := testing.NewDeleteCollectionAction(personsResource, c.ns, listOptions)

	_, err := c.Fake.Invokes(action, &aslangroup_io_v1.PersonList{})
	return err
}

// Patch applies the patch and returns the patched person.
func (c *FakePersons) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *aslangroup_io_v1.Person, err error) {
	obj, err := c.Fake.
		Invokes(testing.NewPatchSubresourceAction(personsResource, c.ns, name, data, subresources...), &aslangroup_io_v1.Person{})

	if obj == nil {
		return nil, err
	}
	return obj.(*aslangroup_io_v1.Person), err
}
