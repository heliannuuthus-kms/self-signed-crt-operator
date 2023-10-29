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
	GetSpec() *SelfSignedIssuerSpec
	GetStatus() *SelfSignedIssuerStatus
}

var _ GenericIssuer = &SelfSignedIssuer{}
var _ GenericIssuer = &ClusterSelfSignedIssuer{}

// GetObjectMeta returns the k8s object metadata
func (c *ClusterSelfSignedIssuer) GetObjectMeta() *metav1.ObjectMeta {
	return &c.ObjectMeta
}

// GetSpec returns the issuer spec
func (c *ClusterSelfSignedIssuer) GetSpec() *SelfSignedIssuerSpec {
	return &c.Spec
}

// GetStatus returns the issuer status
func (c *ClusterSelfSignedIssuer) GetStatus() *SelfSignedIssuerStatus {
	return &c.Status
}

// SetSpec sets the issuer spec
func (c *ClusterSelfSignedIssuer) SetSpec(spec SelfSignedIssuerSpec) {
	c.Spec = spec
}

// SetStatus sets the issuer status
func (c *ClusterSelfSignedIssuer) SetStatus(status SelfSignedIssuerStatus) {
	c.Status = status
}

// Copy deep copies the issuer
func (c *ClusterSelfSignedIssuer) Copy() GenericIssuer {
	return c.DeepCopy()
}

// GetObjectMeta returns the k8s object metadata
func (c *SelfSignedIssuer) GetObjectMeta() *metav1.ObjectMeta {
	return &c.ObjectMeta
}

// GetSpec returns the issuer spec
func (c *SelfSignedIssuer) GetSpec() *SelfSignedIssuerSpec {
	return &c.Spec
}

// GetStatus returns the issuer status
func (c *SelfSignedIssuer) GetStatus() *SelfSignedIssuerStatus {
	return &c.Status
}

// SetSpec sets the issuer spec
func (c *SelfSignedIssuer) SetSpec(spec SelfSignedIssuerSpec) {
	c.Spec = spec
}

// SetStatus sets the issuer status
func (c *SelfSignedIssuer) SetStatus(status SelfSignedIssuerStatus) {
	c.Status = status
}

// Copy deep copies the issuer
func (c *SelfSignedIssuer) Copy() GenericIssuer {
	return c.DeepCopy()
}
