package config

var (
	config *Config
)

const (
	// DefaultCRDNamespace is the default namespace where we create CRD instances.
	DefaultCRDNamespace string = "crd"
	// DefaultKubeConfigPath is the default local path of kubeconfig file.
	DefaultKubeConfigPath string = "/.kube/config"
)

type Config struct {
	crdNamespace string
	kubeconfig   string
}

func (c *Config) GetCRDNamespace() string {
	return c.crdNamespace
}

func (c *Config) GetKubeconfig() string {
	return c.kubeconfig
}
