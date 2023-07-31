/*
Copyright 2021 The Kubernetes Authors.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// GenericIssuer is an interface for both issuer types to simplify controller code
// +k8s:deepcopy-gen=false
type GenericIssuer interface {
	runtime.Object
	metav1.Object

	GetObjectMeta() *metav1.ObjectMeta
	GetSpec() *CAIssuerSpec
	GetStatus() *CAIssuerStatus
}

var _ GenericIssuer = &CAIssuer{}
var _ GenericIssuer = &ClusterCAIssuer{}

// GetObjectMeta returns the k8s object metadata
func (c *ClusterCAIssuer) GetObjectMeta() *metav1.ObjectMeta {
	return &c.ObjectMeta
}

// GetSpec returns the issuer spec
func (c *ClusterCAIssuer) GetSpec() *CAIssuerSpec {
	return &c.Spec
}

// GetStatus returns the issuer status
func (c *ClusterCAIssuer) GetStatus() *CAIssuerStatus {
	return &c.Status
}

// SetSpec sets the issuer spec
func (c *ClusterCAIssuer) SetSpec(spec CAIssuerSpec) {
	c.Spec = spec
}

// SetStatus sets the issuer status
func (c *ClusterCAIssuer) SetStatus(status CAIssuerStatus) {
	c.Status = status
}

// Copy deep copies the issuer
func (c *ClusterCAIssuer) Copy() GenericIssuer {
	return c.DeepCopy()
}

// GetObjectMeta returns the k8s object metadata
func (c *CAIssuer) GetObjectMeta() *metav1.ObjectMeta {
	return &c.ObjectMeta
}

// GetSpec returns the issuer spec
func (c *CAIssuer) GetSpec() *CAIssuerSpec {
	return &c.Spec
}

// GetStatus returns the issuer status
func (c *CAIssuer) GetStatus() *CAIssuerStatus {
	return &c.Status
}

// SetSpec sets the issuer spec
func (c *CAIssuer) SetSpec(spec CAIssuerSpec) {
	c.Spec = spec
}

// SetStatus sets the issuer status
func (c *CAIssuer) SetStatus(status CAIssuerStatus) {
	c.Status = status
}

// Copy deep copies the issuer
func (c *CAIssuer) Copy() GenericIssuer {
	return c.DeepCopy()
}
