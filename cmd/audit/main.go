package main

import (
	"context"
	"fmt"
	"os"
	"time"

	helm "github.com/mittwald/go-helm-client"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"

	"github.com/michaelcourcy/audit-tool/pkg/action"
	"github.com/michaelcourcy/audit-tool/pkg/client"
	"github.com/michaelcourcy/audit-tool/pkg/profile"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	log.SetLevel(log.WarnLevel)
}

func main() {
	kastenNamespace, kastenRelease := getKastenNamespaceAndRelease()

	// create the necessary client
	config, err := client.Config()
	if err != nil {
		panic(err.Error())
	}
	corev1Client, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	actionClient, err := client.ActionClient(config)
	if err != nil {
		panic(err)
	}

	profileClient, err := client.ProfileClient(config)
	if err != nil {
		panic(err)
	}

	discoveryClient, err := client.DiscoveryClient(config)
	if err != nil {
		panic(err)
	}

	helmClient, err := client.HelmClient(config, kastenNamespace)
	if err != nil {
		panic(err)
	}

	err = clusterInfo(corev1Client, discoveryClient)
	if err != nil {
		panic(err)
	}

	// audits
	err = checkKastenInstall(corev1Client, kastenNamespace, kastenRelease, helmClient)
	if err != nil {
		panic(err)
	}

	err = profilesAudit(profileClient, kastenNamespace)
	if err != nil {
		panic(err)
	}

	namespacesWithPVCs, err := rpoForNamespaceWithPVC(corev1Client, actionClient, kastenNamespace)
	if err != nil {
		panic(err)
	}

	err = rpoForNamespaceWithoutPVC(corev1Client, actionClient, namespacesWithPVCs)
	if err != nil {
		panic(err)
	}

}

func checkKastenInstall(corev1Client *kubernetes.Clientset, kastenNamespace string, kastenRelease string, helmClient helm.Client) error {
	fmt.Println("")
	fmt.Println("=======================")
	fmt.Println("Checking Kasten install")
	fmt.Println("=======================")
	fmt.Printf("  checking if Kasten is installed in namespace %s under the release %s \n", kastenNamespace, kastenRelease)
	//ns kasten exist
	_, err := corev1Client.CoreV1().Namespaces().Get(context.TODO(), kastenNamespace, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		fmt.Printf("  --> WARNING !! %s namespace not found, kasten is maybe installed in another namespace \n", kastenNamespace)
		return fmt.Errorf("%s namespace not found, kasten is maybe installed in another namespace", kastenNamespace)
	}
	//release exist
	release, err := helmClient.GetRelease(kastenRelease)
	if err != nil {
		return err
	}
	fmt.Printf("  The version of kasten is %s \n", release.Chart.Metadata.AppVersion)
	//all pods healthy
	results, err := corev1Client.CoreV1().Pods(kastenNamespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	podsInError := false
	for _, pod := range results.Items {
		for _, c := range pod.Status.ContainerStatuses {
			if c.State.Waiting != nil && c.State.Waiting.Reason == "CrashLoopBackOff" {
				podsInError = true
				fmt.Printf("  pod %s is in error: %s\n", pod.ObjectMeta.Name, pod.Status.Phase)
			}
		}
	}
	if podsInError {
		fmt.Printf("  --> WARNING !! some pods in the kasten namespace %s are on error \n", kastenNamespace)
	} else {
		fmt.Printf("  --> No pods in the kasten namespace %s are on error \n", kastenNamespace)
	}
	return nil
}

func profilesAudit(profileClient *rest.RESTClient, kastenNamespace string) error {

	fmt.Println("")
	fmt.Println("=================")
	fmt.Println("Auditing profiles")
	fmt.Println("=================")

	result := profile.ProfileList{}
	err := profileClient.
		Get().
		Resource("profiles").Namespace(kastenNamespace).
		Do(context.TODO()).
		Into(&result)
	if err != nil {
		return err
	} else {
		if len(result.Items) == 0 {
			fmt.Println("  -> WARNING !! there is no profile at all, you don't have real backup ")
			return nil
		}
		foundLocationProfile := false
		foundImmutable := false
		for _, profileInKasten := range result.Items {
			if profileInKasten.Spec.Type == "Location" {
				foundLocationProfile = true
				if profileInKasten.Spec.LocationSpec.Location.LocationType == "ObjectStore" {
					if profileInKasten.Spec.LocationSpec.Location.ObjectStore.ProtectionPeriod != "" {
						foundImmutable = true
					}
				}
			}
			if profileInKasten.Status.Validation != "Success" {
				fmt.Printf("  --> WARNING !! found profile %s which is not valid\n", profileInKasten.Name)
			}
		}
		if !foundLocationProfile {
			fmt.Println("  --> WARNING !! there is no location profile at all, you don't have real backup ")
		} else {
			fmt.Println("  At least one location profile was found")
		}
		if foundLocationProfile && !foundImmutable {
			fmt.Println("  --> WARNING !! there is no immutable profile your are not protected against Ransomware ")
		}
	}
	return nil
}

func clusterInfo(corev1Client *kubernetes.Clientset, discoveryClient *discovery.DiscoveryClient) error {
	fmt.Println("")
	fmt.Println("=============================")
	fmt.Println("Information about the cluster")
	fmt.Println("=============================")
	information, err := discoveryClient.ServerVersion()
	if err != nil {
		return err
	}
	fmt.Printf("  Kubernetes version : %s.%s \n", information.Major, information.Minor)
	fmt.Printf("  Platform : %s \n", information.Platform)

	nodes, err := corev1Client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	numNodes := len(nodes.Items)
	fmt.Printf("  There are %d nodes in this cluster, looking for nodes in error\n", numNodes)
	nodeInError := false
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady && condition.Status != v1.ConditionTrue {
				nodeInError = true
				fmt.Printf("  NodName: %s, Condition type: %s %s", node.Name, condition.Type, condition.Status)
			}
		}
	}
	if !nodeInError {
		fmt.Printf("  --> No nodes are in error\n")
	} else {
		fmt.Printf("  --> Some nodes are in error\n")
	}

	return nil
}

func rpoForNamespaceWithoutPVC(corev1Client *kubernetes.Clientset, actionClient *rest.RESTClient, namespacesWithPVCs map[string][]string) error {
	var namespacesWithoutPVCs []string
	namespaces, err := corev1Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}
	for _, namespace := range namespaces.Items {
		_, ok := namespacesWithPVCs[namespace.Name]
		if !ok {
			namespacesWithoutPVCs = append(namespacesWithoutPVCs, namespace.Name)
		}
	}

	log.WithFields(log.Fields{
		"namespacesWithoutPVCs": namespacesWithoutPVCs,
	}).Info("namespace without pvcs")

	fmt.Println("")
	fmt.Println("======================")
	fmt.Println("Namespaces without PVC")
	fmt.Println("======================")

	for _, namespace := range namespacesWithoutPVCs {
		fmt.Printf("%s has no PVC\n", namespace)
		rpo(namespace, actionClient)
	}
	return nil
}

func rpoForNamespaceWithPVC(corev1Client *kubernetes.Clientset, actionClient *rest.RESTClient, kastenNamespace string) (map[string][]string, error) {
	//list all namespaces that has pvc
	pvcs, err := corev1Client.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return map[string][]string{}, err
	}
	namespacesWithPVCs := make(map[string][]string)
	for _, pvc := range pvcs.Items {
		pvcsInNamespace, ok := namespacesWithPVCs[pvc.Namespace]
		if ok {
			namespacesWithPVCs[pvc.Namespace] = append(pvcsInNamespace, pvc.Name)
		} else {
			var pvcsInNamespace []string
			pvcsInNamespace = append(pvcsInNamespace, pvc.Name)
			namespacesWithPVCs[pvc.Namespace] = pvcsInNamespace
		}
	}

	log.WithFields(log.Fields{
		"namespacesWithPVCs": namespacesWithPVCs,
	}).Info("namespace with pvcs")

	fmt.Println("")
	fmt.Println("===================")
	fmt.Println("Namespaces with PVC")
	fmt.Println("===================")
	for namespace, pvcs := range namespacesWithPVCs {
		if namespace == kastenNamespace {
			continue
		}
		fmt.Printf("%s has %d PVCs\n", namespace, len(pvcs))
		rpo(namespace, actionClient)
	}
	return namespacesWithPVCs, nil
}

func rpo(namespace string, actionClient *rest.RESTClient) error {
	result := action.BackupActionList{}
	err := actionClient.
		Get().
		Resource("backupactions").Namespace(namespace).
		Do(context.TODO()).
		Into(&result)
	if err != nil {
		return err
	} else {
		if len(result.Items) == 0 {
			fmt.Printf("  --> No backupactions in namespace %s \n", namespace)
		} else {
			padding := "  %-30s %-10s %-30s %-30s \n"
			first := true
			completeBackupAction := false
			var rpoMessage string
			fmt.Printf(padding, "BACKUPACTION", "STATE", "START", "STOP")
			for _, backupAction := range result.Items {
				if first && backupAction.Status.State == "Complete" {
					first = false
					completeBackupAction = true
					now := time.Now()
					rpoDuration := now.Sub(backupAction.Status.EndTime)
					days := int64(rpoDuration.Hours() / 24)
					hours := int64(rpoDuration.Hours()) % 24
					rpoMessage = fmt.Sprintf("  --> The last RPO is %d days and %d hours", days, hours)
				}
				fmt.Printf(padding, backupAction.Name, backupAction.Status.State, backupAction.CreationTimestamp, backupAction.Status.EndTime)
			}
			if !completeBackupAction {
				fmt.Println("  --> WARNING !! It seems that no backupaction were successful")
			} else {
				fmt.Println(rpoMessage)
			}

		}
	}
	return nil
}

func getKastenNamespaceAndRelease() (string, string) {
	namespace := os.Getenv("KASTEN_NAMESPACE")
	if namespace == "" {
		namespace = "kasten-io"
	}
	release := os.Getenv("KASTEN_RELEASE")
	if release == "" {
		release = "k10"
	}
	return namespace, release
}
