package config

var (
	config *Config
)

const (
	// DefaultCRDNamespace is the default namespace where we create CRD instances.
	DefaultCRDNamespace string = "crd"
	// DefaultKubeconfigPath is the default local path of kubeconfig file.
	DefaultKubeconfigPath string = "/.kube/config"
)

type Config struct {
	crdNamespace   string
	kubeconfigPath string
}

func (c *Config) GetCRDNamespace() string {
	return c.crdNamespace
}

func (c *Config) GetKubeconfigPath() string {
	return c.kubeconfigPath
}
