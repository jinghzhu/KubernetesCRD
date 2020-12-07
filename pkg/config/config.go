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

	config.kubeconfigPath = os.Getenv("CRD_KUBECONFIG")
	if config.kubeconfigPath == "" {
		config.kubeconfigPath = DefaultKubeconfigPath
	}
}

func GetConfig() *Config {
	return config
}
