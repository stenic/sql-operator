package v1alpha1

type SqlHostEndpoint struct {
	Host string `json:"host,omitempty"`
	Port int    `json:"port,omitempty"`
}

type SqlCredentials struct {
	Username string `json:"username,omitempty"`
	Password string `json:"password,omitempty"`
}

type SqlHostRef struct {
	Name      string `json:"name,omitempty"`
	Namespace string `json:"namespace,omitempty"`
}
