package command

import (
	"context"
	"fmt"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apiextensions-apiserver/pkg/client/clientset/clientset/scheme"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func GetRelease(config *rest.Config, cs *kubernetes.Clientset, m map[string][]string) {

	// var stdout, stderr bytes.Buffer
	// for _, v := range m {
	if true {

		ns := "default"
		pod := "go-fiber-675455d847-6z785"

		fmt.Printf("ns=%s,pod=%s\n", ns, pod)
		req := cs.CoreV1().RESTClient().Post().
			Resource("pods").
			Namespace(ns).
			Name(pod).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Container: "app",
				Command:   []string{"/bin/sh", "-c", "cat /etc/os-release"},
				Stdout:    true,
				Stderr:    true,
			}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
		if err != nil {
			fmt.Println("exec err:", err)
		}

		err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
			Stdout: os.Stdout,
			Stderr: os.Stderr,
		})

		if err != nil {
			fmt.Println("exec stream err:", err)
		}

	}
}
