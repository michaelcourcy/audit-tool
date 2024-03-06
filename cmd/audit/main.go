package main

import (
	"context"
	"fmt"
	"os"
	"time"

	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/rest"

	v1alpha1 "github.com/michaelcourcy/audit-tool/pkg/action"
	"github.com/michaelcourcy/audit-tool/pkg/client"
	v1 "k8s.io/api/core/v1"
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
	fmt.Printf("Kasten is installed in %s under the release %s \n", kastenNamespace, kastenRelease)

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

	namespacesWithPVCs := rpoForNamespaceWithPVC(corev1Client, actionClient, kastenNamespace)

	rpoForNamespaceWithoutPVC(corev1Client, actionClient, namespacesWithPVCs)

	clusterInfo(corev1Client)

}

func clusterInfo(corev1Client *kubernetes.Clientset) {
	fmt.Println("")
	fmt.Println("=============================")
	fmt.Println("Information about the cluster")
	fmt.Println("=============================")
	nodes, err := corev1Client.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
	}
	numNodes := len(nodes.Items)
	fmt.Printf("There are %d nodes in this cluster, looking for nodes in error\n", numNodes)
	nodeInError := false
	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == v1.NodeReady && condition.Status != v1.ConditionTrue {
				nodeInError = true
				fmt.Printf("NodName: %s, Condition type: %s %s", node.Name, condition.Type, condition.Status)
			}
		}
	}
	if !nodeInError {
		fmt.Printf("No nodes are in error\n")
	}
}

func rpoForNamespaceWithoutPVC(corev1Client *kubernetes.Clientset, actionClient *rest.RESTClient, namespacesWithPVCs map[string][]string) {
	var namespacesWithoutPVCs []string
	namespaces, err := corev1Client.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
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
}

func rpoForNamespaceWithPVC(corev1Client *kubernetes.Clientset, actionClient *rest.RESTClient, kastenNamespace string) map[string][]string {
	//list all namespaces that has pvc
	pvcs, err := corev1Client.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err.Error())
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
		result := v1alpha1.BackupActionList{}
		err = actionClient.
			Get().
			Resource("backupactions").Namespace(namespace).
			Do(context.TODO()).
			Into(&result)
		if err != nil {
			fmt.Println(err)
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
	}
	return namespacesWithPVCs
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
