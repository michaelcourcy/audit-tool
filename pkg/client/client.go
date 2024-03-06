package client

import (
	"fmt"

	v1alpha1 "github.com/michaelcourcy/audit-tool/pkg/action"

	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/runtime/serializer"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func Config() (*rest.Config, error) {
	fmt.Println("")
	fmt.Println("======================")
	fmt.Println("Runtime for audit tool")
	fmt.Println("======================")
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	// or load it from local kubeconfig
	if err != nil {
		fmt.Println("Audit tool is not executing in pod")
		// if you want to change the loading rules (which files in which order), you can do so here
		loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
		// if you want to change override values or bind them to flags, there are methods to help you
		configOverrides := &clientcmd.ConfigOverrides{}
		kubeConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, configOverrides)
		config, err = kubeConfig.ClientConfig()
		if err != nil {
			return nil, err
		} else {
			return config, nil
		}
	}
	fmt.Println("Audit tool is executing in pod")
	return config, nil
}

func ActionClient(config *rest.Config) (*rest.RESTClient, error) {
	v1alpha1.AddToScheme(scheme.Scheme)
	apiConfig := *config
	apiConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: v1alpha1.GroupName, Version: v1alpha1.GroupVersion}
	apiConfig.APIPath = "/apis"
	apiConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	apiConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	return rest.UnversionedRESTClientFor(&apiConfig)
}
