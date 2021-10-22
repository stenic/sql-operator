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

type SqlHostRef struct {
	// Name of the SqlHost
	Name string `json:"name,omitempty"`
	// Namespace of the SqlHost
	Namespace string `json:"namespace,omitempty"`
}
