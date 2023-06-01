package main

import (
	"context"
	"kube-exec/command"
	"log"
	"strings"

	k8sCrdClient "github.com/Tsingshen/k8scrd/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

type PodContainer struct {
	DeploymentName string
	NameSpace      string
	PodName        string
	ContainerName  []string
}

func main() {
	cs := k8sCrdClient.GetClient()
	config := k8sCrdClient.GetRestConfig()

	podsList, err := cs.CoreV1().Pods(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
	}
	m := getPodDeployWithSplit(cs, podsList.Items)
	command.GetRelease(config, cs, m)

}

func getPodDeployWithSplit(cs *kubernetes.Clientset, pods []corev1.Pod) map[string][]string {

	nsDeployPod := make(map[string][]string)
	containerName := "app"

	for _, p := range pods {
		if len(p.GetOwnerReferences()) > 0 {
			if p.OwnerReferences[0].Kind == "ReplicaSet" {
				hash := p.GetLabels()["pod-template-hash"]
				dn := strings.Split(p.GenerateName, "-"+hash+"-")[0]
				nsPod := []string{p.Namespace, p.Name}
				if checkContainer(&p, containerName) {
					nsDeployPod[dn] = nsPod
				}
			}
		}
	}

	return nsDeployPod
}

func checkContainer(p *corev1.Pod, name string) bool {
	for _, v := range p.Spec.Containers {
		if v.Name == name {
			return true
		}
	}

	return false
}
