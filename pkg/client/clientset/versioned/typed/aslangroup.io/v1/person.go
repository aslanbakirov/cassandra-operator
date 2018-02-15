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

package v1

import (
	v1 "github.com/aslanbekirov/personcrd/pkg/apis/aslangroup.io/v1"
	scheme "github.com/aslanbekirov/personcrd/pkg/client/clientset/versioned/scheme"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	rest "k8s.io/client-go/rest"
)

// PersonsGetter has a method to return a PersonInterface.
// A group's client should implement this interface.
type PersonsGetter interface {
	Persons(namespace string) PersonInterface
}

// PersonInterface has methods to work with Person resources.
type PersonInterface interface {
	Create(*v1.Person) (*v1.Person, error)
	Update(*v1.Person) (*v1.Person, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error
	Get(name string, options meta_v1.GetOptions) (*v1.Person, error)
	List(opts meta_v1.ListOptions) (*v1.PersonList, error)
	Watch(opts meta_v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Person, err error)
	PersonExpansion
}

// persons implements PersonInterface
type persons struct {
	client rest.Interface
	ns     string
}

// newPersons returns a Persons
func newPersons(c *AslangroupV1Client, namespace string) *persons {
	return &persons{
		client: c.RESTClient(),
		ns:     namespace,
	}
}

// Get takes name of the person, and returns the corresponding person object, and an error if there is any.
func (c *persons) Get(name string, options meta_v1.GetOptions) (result *v1.Person, err error) {
	result = &v1.Person{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("persons").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Persons that match those selectors.
func (c *persons) List(opts meta_v1.ListOptions) (result *v1.PersonList, err error) {
	result = &v1.PersonList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("persons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested persons.
func (c *persons) Watch(opts meta_v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("persons").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a person and creates it.  Returns the server's representation of the person, and an error, if there is any.
func (c *persons) Create(person *v1.Person) (result *v1.Person, err error) {
	result = &v1.Person{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("persons").
		Body(person).
		Do().
		Into(result)
	return
}

// Update takes the representation of a person and updates it. Returns the server's representation of the person, and an error, if there is any.
func (c *persons) Update(person *v1.Person) (result *v1.Person, err error) {
	result = &v1.Person{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("persons").
		Name(person.Name).
		Body(person).
		Do().
		Into(result)
	return
}

// Delete takes name of the person and deletes it. Returns an error if one occurs.
func (c *persons) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("persons").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *persons) DeleteCollection(options *meta_v1.DeleteOptions, listOptions meta_v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("persons").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched person.
func (c *persons) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1.Person, err error) {
	result = &v1.Person{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("persons").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}
