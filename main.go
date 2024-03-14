package main

import (
	"context"
	"kube-exec/command"
	"kube-exec/util"
	"log/slog"
	"os"

	k8sCrdClient "github.com/changqings/k8scrd/client"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func main() {
	cs := k8sCrdClient.GetClient()
	config := k8sCrdClient.GetRestConfig()

	// list all pods
	podsList, err := cs.CoreV1().Pods(corev1.NamespaceAll).List(context.Background(), metav1.ListOptions{})
	if err != nil {
		slog.Error("list all pods err", "msg", err)
		os.Exit(1)
	}

	// get deployment:pod 1:1
	pcs := util.GetPodDeployWithSplit(podsList.Items)

	// get pod /etc/os-release
	command.GetRelease(config, cs, pcs)
}
