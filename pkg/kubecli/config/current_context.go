package config

import "k8s.io/client-go/tools/clientcmd"

func CurrentContext() (string, error) {
	pathOptions := clientcmd.NewDefaultPathOptions()
	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		return "", err
	}

	if config.CurrentContext == "" {
		return "current-context is not set", nil
	}
	return config.CurrentContext, nil
}
