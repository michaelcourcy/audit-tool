package client

import (
	"fmt"

	action "github.com/michaelcourcy/audit-tool/pkg/action"
	"github.com/michaelcourcy/audit-tool/pkg/profile"
	helm "github.com/mittwald/go-helm-client"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/discovery"

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
	action.AddToScheme(scheme.Scheme)
	apiConfig := *config
	apiConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: action.GroupName, Version: action.GroupVersion}
	apiConfig.APIPath = "/apis"
	apiConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	apiConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	return rest.UnversionedRESTClientFor(&apiConfig)
}

func ProfileClient(config *rest.Config) (*rest.RESTClient, error) {
	profile.AddToScheme(scheme.Scheme)
	apiConfig := *config
	apiConfig.ContentConfig.GroupVersion = &schema.GroupVersion{Group: profile.GroupName, Version: profile.GroupVersion}
	apiConfig.APIPath = "/apis"
	apiConfig.NegotiatedSerializer = serializer.NewCodecFactory(scheme.Scheme)
	apiConfig.UserAgent = rest.DefaultKubernetesUserAgent()

	return rest.UnversionedRESTClientFor(&apiConfig)
}

func HelmClient(config *rest.Config, kastenNamespace string) (helm.Client, error) {
	opt := &helm.RestConfClientOptions{
		Options: &helm.Options{
			Namespace:        kastenNamespace,
			RepositoryCache:  "/tmp/.helmcache",
			RepositoryConfig: "/tmp/.helmrepo",
			Debug:            false,
			Linting:          false,
			DebugLog: func(format string, v ...interface{}) {
				log.Debugf(format, v...)
			},
		},
		RestConfig: config,
	}
	return helm.NewClientFromRestConf(opt)
}

func DiscoveryClient(config *rest.Config) (*discovery.DiscoveryClient, error) {
	return discovery.NewDiscoveryClientForConfig(config)
}
