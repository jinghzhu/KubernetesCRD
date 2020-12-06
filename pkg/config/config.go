package config

import "os"

func init() {
	initConfig()
}

func initConfig() {
	config = &Config{}
	config.crdNamespace = os.Getenv("CRD_NAMESPACE")
	if config.crdNamespace == "" {
		config.crdNamespace = DefaultCRDNamespace
	}

	config.kubeconfig = os.Getenv("KUBECONFIG")
	if config.kubeconfig == "" {
		config.kubeconfig = DefaultKubeConfigPath
	}
}

func GetConfig() *Config {
	return config
}
