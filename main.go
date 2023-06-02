package main

import (
	"context"
	"kube-exec/command"
	"kube-exec/util"
	"log"
	"strings"

	k8sCrdClient "github.com/Tsingshen/k8scrd/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	cs := k8sCrdClient.GetClient()
	config := k8sCrdClient.GetRestConfig()

	podsList, err := cs.CoreV1().Pods(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		log.Println(err)
	}
	pcs := getPodDeployWithSplit(cs, podsList.Items)
	command.GetRelease(config, cs, pcs)

}

func getPodDeployWithSplit(cs *kubernetes.Clientset, pods []corev1.Pod) util.PcSlice {

	pcs := util.PcSlice{}
	containerName := "app"

	for _, p := range pods {
		if len(p.GetOwnerReferences()) > 0 {
			if p.OwnerReferences[0].Kind == "ReplicaSet" {
				hash := p.GetLabels()["pod-template-hash"]
				dpname := strings.Split(p.GenerateName, "-"+hash+"-")[0]
				if checkContainerNameAndStat(&p, containerName) {
					pc := util.PodContainer{
						ContainerName:  containerName,
						NameSpace:      p.Namespace,
						PodName:        p.Name,
						DeploymentName: dpname,
					}

					pcs = append(pcs, pc)
				}
			}
		}
	}

	return pcs
}

func checkContainerNameAndStat(p *corev1.Pod, name string) bool {
	for _, v := range p.Spec.Containers {
		if v.Name == name {
			for _, sc := range p.Status.Conditions {
				if sc.Type == corev1.PodReady {
					if sc.Status == corev1.ConditionTrue {
						return true
					}

				}
			}
		}
	}

	return false
}
