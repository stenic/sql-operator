package v1alpha1

type SqlHostEndpoint struct {
	// Host of the sql server. This can be either dns or ip
	Host string `json:"host,omitempty"`

	// Port to use when connecting to the sql server
	Port int `json:"port,omitempty"`
}

type SqlCredentials struct {
	// Username for the sql user
	Username string `json:"username,omitempty"`

	// Password for the sql user
	Password string `json:"password,omitempty"`
}

type SqlObjectRef struct {
	// Name of the SqlObject
	Name string `json:"name,omitempty"`
	// Namespace of the SqlObject
	Namespace string `json:"namespace,omitempty"`
}

// CleanupPolicy describes how the resource will be handled on delete.
// +kubebuilder:validation:Enum=Retain;Delete
type CleanupPolicy string

const (
	// Keep
	CleanupPolicyRetain CleanupPolicy = "Retain"

	// ForbidConcurrent forbids concurrent runs, skipping next run if previous
	// hasn't finished yet.
	CleanupPolicyDelete CleanupPolicy = "Delete"
)
