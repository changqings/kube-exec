package command

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type PodContainer struct {
	DeploymentName string
	NameSpace      string
	PodName        string
	ContainerName  string
}

type deployNs struct {
	DeploymentName string
	NameSpace      string
}

func GetRelease(config *rest.Config, cs *kubernetes.Clientset, pcs []PodContainer) {

	pcsUinq := deploymentWithOnePod(pcs)

	for _, pc := range pcsUinq {

		req := cs.CoreV1().RESTClient().Post().
			Resource("pods").
			Namespace(pc.NameSpace).
			Name(pc.PodName).
			SubResource("exec").
			VersionedParams(&corev1.PodExecOptions{
				Container: pc.ContainerName,
				Command:   []string{"cat", "/etc/os-release"},
				Stdin:     true,
				Stdout:    true,
				TTY:       false,
			}, scheme.ParameterCodec)

		exec, err := remotecommand.NewSPDYExecutor(config, http.MethodPost, req.URL())
		if err != nil {
			fmt.Println("exec err:", err)
		}

		var stdout, stderr bytes.Buffer

		err = exec.StreamWithContext(context.Background(), remotecommand.StreamOptions{
			Stdin:  nil,
			Stdout: &stdout,
			Stderr: &stderr,
			Tty:    false,
		})

		if err != nil {
			fmt.Println("exec stream err:", err)
		}
		scan := bufio.NewScanner(&stdout)
		for scan.Scan() {
			if strings.HasPrefix(scan.Text(), "ID=") {
				fmt.Printf("pod %s.%s /etc/os-release info: os=%s\n", pc.NameSpace, pc.PodName, strings.Split(scan.Text(), "ID=")[1])
				break
			}
		}
	}
}

func deploymentWithOnePod(pcs []PodContainer) []PodContainer {

	m := make(map[deployNs]PodContainer)

	for _, pc := range pcs {
		dn := deployNs{pc.DeploymentName, pc.NameSpace}
		if _, ok := m[dn]; !ok {
			m[dn] = pc
		}
	}

	result := make([]PodContainer, 0, len(m))
	for _, pc := range m {
		result = append(result, pc)
	}

	return result
}
