package config

import (
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
)

var (
	config *clientcmdapi.Config
)

func init() {
	var err error
	pathOptions := clientcmd.NewDefaultPathOptions()
	config, err = pathOptions.GetStartingConfig()
	if err != nil {
		panic(err)
	}
}

func CurrentContext() string {
	return config.CurrentContext
}

func SetCurrentContext(context string) {
	config.CurrentContext = context
}

func ContextNamespace() string {
	ctx, ok := config.Contexts[config.CurrentContext]
	if !ok {
		return ""
	}
	ns := ctx.Namespace
	return ns
}

func ListContexts() []string {
	contexts := make([]string, 0)
	for name := range config.Contexts {
		contexts = append(contexts, name)
	}
	return contexts
}
