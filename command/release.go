package command

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

func GetRelease(config *rest.Config, cs *kubernetes.Clientset, m map[string][]string) {

	for _, v := range m {

		ns := v[0]
		pod := v[1]
		fmt.Printf("pod %s.%s /etc/os-release info:\n", ns, pod)

		fmt.Printf("ns=%s,pod=%s\n", ns, pod)
		req := cs.CoreV1().RESTClient().Post().
			Resource("pods").
			Namespace(ns).
			Name(pod).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Container: "app",
				Command:   []string{"cat", "/etc/os-release"},
				Stdin:     true,
				Stdout:    true,
				TTY:       false,
			}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
		if err != nil {
			fmt.Println("exec err:", err)
		}

		err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: io.Writer(os.Stdout),
			Stderr: io.Writer(os.Stderr),
			Tty:    false,
		})

		if err != nil {
			fmt.Println("exec stream err:", err)
		}

		fmt.Println("------")

	}
}
