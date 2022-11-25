package v1alpha1

import (
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// MicroservicesSpec defines the desired state of Microservices.
type MicroservicesSpec struct {
	// NOTE: currently this field is not considered.
	// Components defines requierements for a set of tempo components.
	//
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Tempo Components"
	Components TempoComponentsSpec `json:"template,omitempty"`

	// NOTE: currently this field is not considered.
	// The resources are split in between components.
	// Tempo operator knows how to split them appropriately based on grafana/tempo/issues/1540.
	//
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Resource Requirements"
	Resources *corev1.ResourceRequirements `json:"resources,omitempty"`

	// NOTE: currently this field is not considered.
	// Storage defines S3 compatible object storage configuration.
	// User is required to create secret and supply it.
	//
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Object Storage"
	Storage ObjectStorageSpec `json:"storage,omitempty"`

	// NOTE: currently this field is not considered.
	// Retention period defined by dataset.
	// User can specify how long data should be stored.
	//
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Retention Period"
	Retention RetentionSpec `json:"retention,omitempty"`

	// StorageClassName for PVCs used by ingester. Defaults to nil (default storage class in the cluster).
	//
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="StorageClassName for PVCs"
	StorageClassName *string `json:"storageClassName,omitempty"`

	// NOTE: currently this field is not considered.
	// LimitSpec is used to limit ingestion and querying rates.
	//
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Ingestion and Querying Ratelimiting"
	LimitSpec LimitSpec `json:"limits,omitempty"`

	// StorageSize for PVCs used by ingester. Defaults to 10Gi.
	//
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Storage size for PVCs"
	StorageSize resource.Quantity `json:"storageSize,omitempty"`

	// NOTE: currently this field is not considered.
	// ReplicationFactor is used to define how many component replicas should exist.
	//
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Replication Factor"
	ReplicationFactor int `json:"replicationFactor,omitempty"`

	// Tenants defines the per-tenant authentication and authorization spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Tenants Configuration"
	Tenants *TenantsSpec `json:"tenants,omitempty"`
}

// MicroservicesStatus defines the observed state of Microservices.
type MicroservicesStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Microservices is the Schema for the microservices API.
type Microservices struct {
	Status            MicroservicesStatus `json:"status,omitempty"`
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              MicroservicesSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// MicroservicesList contains a list of Microservices.
type MicroservicesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Microservices `json:"items"`
}

// ModeType is the authentication/authorization mode in which LokiStack Gateway will be configured.
//
// +kubebuilder:validation:Enum=openshift
type ModeType string

const (
	// Static mode asserts the Authorization Spec's Roles and RoleBindings
	// using an in-process OpenPolicyAgent Rego authorizer.
	OpenShift ModeType = "openshift"
)

type TenantsSpec struct {
	// Mode defines the multitenancy mode.
	//
	// +required
	// +kubebuilder:validation:Required
	// +kubebuilder:default:=openshift
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors={"urn:alm:descriptor:com.tectonic.ui:select:static","urn:alm:descriptor:com.tectonic.ui:select:dynamic","urn:alm:descriptor:com.tectonic.ui:select:openshif","urn:alm:descriptor:com.tectonic.ui:select:openshift"},displayName="Mode"
	Mode ModeType `json:"mode"`

	// Authentication defines the lokistack-gateway component authentication configuration spec per tenant.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Authentication"
	Authentication []AuthenticationSpec `json:"authentication,omitempty"`
	// Authorization defines the lokistack-gateway component authorization configuration spec per tenant.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Authorization"
	Authorization *AuthorizationSpec `json:"authorization,omitempty"`
}

// PermissionType is a LokiStack Gateway RBAC permission.
//
// +kubebuilder:validation:Enum=read;write
type PermissionType string

const (
	// Write gives access to write data to a tenant.
	Write PermissionType = "write"
	// Read gives access to read data from a tenant.
	Read PermissionType = "read"
)

// SubjectKind is a kind of LokiStack Gateway RBAC subject.
//
// +kubebuilder:validation:Enum=user;group
type SubjectKind string

const (
	// User represents a subject that is a user.
	User SubjectKind = "user"
	// Group represents a subject that is a group.
	Group SubjectKind = "group"
)

// Subject represents a subject that has been bound to a role.
type Subject struct {
	Name string      `json:"name"`
	Kind SubjectKind `json:"kind"`
}

// RoleBindingsSpec binds a set of roles to a set of subjects.
type RoleBindingsSpec struct {
	Name     string    `json:"name"`
	Subjects []Subject `json:"subjects"`
	Roles    []string  `json:"roles"`
}

// AuthorizationSpec defines the opa, role bindings and roles
// configuration per tenant for lokiStack Gateway component.
type AuthorizationSpec struct {
	// Roles defines a set of permissions to interact with a tenant.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Static Roles"
	Roles []RoleSpec `json:"roles"`
	// RoleBindings defines configuration to bind a set of roles to a set of subjects.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Static Role Bindings"
	RoleBindings []RoleBindingsSpec `json:"roleBindings"`
}

// RoleSpec describes a set of permissions to interact with a tenant.
type RoleSpec struct {
	Name        string           `json:"name"`
	Resources   []string         `json:"resources"`
	Tenants     []string         `json:"tenants"`
	Permissions []PermissionType `json:"permissions"`
}

// AuthenticationSpec defines the oidc configuration per tenant for lokiStack Gateway component.
type AuthenticationSpec struct {
	// TenantName defines the name of the tenant.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Tenant Name"
	TenantName string `json:"tenantName"`
	// TenantID defines the id of the tenant.
	//
	// +required
	// +kubebuilder:validation:Required
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Tenant ID"
	TenantID string `json:"tenantId"`
}

// ObjectStorageSpec defines the requirements to access the object
// storage bucket to persist traces by the ingester component.
type ObjectStorageSpec struct {
	// TLS configuration for reaching the object storage endpoint.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="TLS Config"
	TLS *ObjectStorageTLSSpec `json:"tls,omitempty"`

	// Secret for object storage authentication.
	// Name of a secret in the same namespace as the tempo Microservices custom resource.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Object Storage Secret"
	Secret string `json:"secret,omitempty"`
}

// ObjectStorageTLSSpec is the TLS configuration for reaching the object storage endpoint.
type ObjectStorageTLSSpec struct {
	// CA is the name of a ConfigMap containing a CA certificate.
	// It needs to be in the same namespace as the LokiStack custom resource.
	//
	// +optional
	// +kubebuilder:validation:optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:io.kubernetes:ConfigMap",displayName="CA ConfigMap Name"
	CA string `json:"caName,omitempty"`
}

// TempoComponentsSpec defines the template of all requirements to configure
// scheduling of all Tempo components to be deployed.
type TempoComponentsSpec struct {
	// Distributor defines the distributor component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Distributor pods"
	Distributor *TempoComponentSpec `json:"distributor,omitempty"`

	// Ingester defines the ingester component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Ingester pods"
	Ingester *TempoComponentSpec `json:"ingester,omitempty"`

	// Compactor defines the lokistack compactor component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Compactor pods"
	Compactor *TempoComponentSpec `json:"compactor,omitempty"`

	// Querier defines the querier component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Querier pods"
	Querier *TempoComponentSpec `json:"querier,omitempty"`

	// TempoQueryFrontendSpec defines the query frontend spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Query Frontend pods"
	QueryFrontend *TempoQueryFrontendSpec `json:"queryFrontend,omitempty"`

	// Gateway defines the gateway component spec.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Gateway pods"
	Gateway *TempoComponentSpec `json:"gateway,omitempty"`
}

// TempoComponentSpec defines specific schedule settings for tempo components.
type TempoComponentSpec struct {
	// Replicas represents the number of replicas to create for this component.
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Component Replicas"
	Replicas *int32 `json:"replicas,omitempty"`

	// NodeSelector is the simplest recommended form of node selection constraint.
	//
	// +optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Node Selector"
	NodeSelector map[string]string `json:"nodeSelector,omitempty"`

	// Tolerations defines component specific pod tolerations.
	//
	// +optional
	// +listType=atomic
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Tolerations"
	Tolerations []corev1.Toleration `json:"tolerations,omitempty"`
}

// TempoQueryFrontendSpec extends TempoComponentSpec with frontend specific parameters.
type TempoQueryFrontendSpec struct {
	// TempoComponentSpec is embedded to extend this definition with further options.
	//
	// Currently there is no way to inline this field.
	// See: https://github.com/golang/go/issues/6213
	//
	// +optional
	// +kubebuilder:validation:Optional
	TempoComponentSpec `json:"component,omitempty"`

	// JaegerQuerySpec defines Jaeger Query spefic options.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Jaeger Query Settings"
	JaegerQuery JaegerQuerySpec `json:"jaegerQuery"`
}

// JaegerQuerySpec defines Jaeger Query options.
type JaegerQuerySpec struct {
	// Enabled is used to define if Jaeger Query component should be created.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Jaeger Query Enabled"
	Enabled bool `json:"enabled"`
}

// LimitSpec defines Global and PerTenant rate limits.
type LimitSpec struct {
	// PerTenant is used to define rate limits per tenant.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Tenant Limits"
	PerTenant map[string]RateLimitSpec `json:"perTenant,omitempty"`

	// Global is used to define global rate limits.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Global Limit"
	Global RateLimitSpec `json:"global"`
}

// RateLimitSpec defines rate limits for Ingestion and Query components.
type RateLimitSpec struct {
	// Ingestion is used to define ingestion rate limits.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Ingestion Limit"
	Ingestion IngestionLimitSpec `json:"ingestion"`

	// Query is used to define query rate limits.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Query Limit"
	Query QueryLimit `json:"query"`
}

// IngestionLimitSpec defines the limits applied at the ingestion path.
type IngestionLimitSpec struct {
	// IngestionBurstSizeBytes defines the burst size (bytes) used in ingestion.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Ingestion Burst Size in Bytes"
	IngestionBurstSizeBytes *int `json:"ingestionBurstSizeBytes,omitempty"`

	// IngestionRateLimitBytes defines the Per-user ingestion rate limit (bytes) used in ingestion.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Ingestion Rate Limit in Bytes"
	IngestionRateLimitBytes *int `json:"ingestionRateLimitBytes,omitempty"`

	// MaxBytesPerTrace defines the maximum number of bytes of an acceptable trace.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Bytes per Trace"
	MaxBytesPerTrace *int `json:"maxBytesPerTrace,omitempty"`

	// MaxTracesPerUser defines the maximum number of traces a user can send.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Traces per User"
	MaxTracesPerUser *int `json:"maxTracesPerUser,omitempty"`
}

// QueryLimit defines query limits.
type QueryLimit struct {
	// MaxBytesPerTagValues defines the maximum size in bytes of a tag-values query.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Tags per User"
	MaxBytesPerTagValues *int `json:"maxBytesPerTagValues,omitempty"`
	// MaxSearchBytesPerTrace defines the maximum size of search data for a single
	// trace in bytes.
	// default: `0` to disable.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:number",displayName="Max Traces per User"
	MaxSearchBytesPerTrace *int `json:"maxSearchBytesPerTrace"`
}

// RetentionSpec defines global and per tenant retention configurations.
type RetentionSpec struct {
	// PerTenant is used to configure retention per tenant.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="PerTenant Retention"
	PerTenant map[string]RetentionConfig `json:"perTenant,omitempty"`
	// Global is used to configure global retention.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,displayName="Global Retention"
	Global RetentionConfig `json:"global"`
}

// RetentionConfig defines how long data should be provided.
type RetentionConfig struct {
	// Traces defines retention period. Supported parameter suffixes are “s”, “m” and “h”.
	// example: 336h
	// default: value is 48h.
	//
	// +optional
	// +kubebuilder:validation:Optional
	// +operator-sdk:csv:customresourcedefinitions:type=spec,xDescriptors="urn:alm:descriptor:com.tectonic.ui:text",displayName="Trace Retention Period"
	Traces time.Duration `json:"traces"`
}

func init() {
	SchemeBuilder.Register(&Microservices{}, &MicroservicesList{})
}
