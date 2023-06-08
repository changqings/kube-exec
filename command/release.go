package command

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"kube-exec/util"
	"net/http"
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type deployNs struct {
	DeploymentName string
	NameSpace      string
}

type osVersion struct {
	id      string
	version string
}

func GetRelease(config *rest.Config, cs *kubernetes.Clientset, pcs util.PcSlice) {

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
			continue
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
			continue
		}

		os := osVersion{}

		scan := bufio.NewScanner(&stdout)
		for scan.Scan() {
			if strings.HasPrefix(scan.Text(), "ID=") {
				os.id = strings.Split(scan.Text(), "ID=")[1]
			}
			if strings.HasPrefix(scan.Text(), "VERSION_ID=") {
				os.version = strings.Split(scan.Text(), "VERSION_ID=")[1]
			}
		}
		fmt.Printf("ns=%s deployment=%s pod=%s /etc/os-release info: os=%s version=%s\n", pc.NameSpace, pc.DeploymentName, pc.PodName, os.id, os.version)
	}
}

func deploymentWithOnePod(pcs util.PcSlice) util.PcSlice {

	m := make(map[deployNs]util.PodContainer)

	for _, pc := range pcs {
		dn := deployNs{pc.DeploymentName, pc.NameSpace}
		if _, ok := m[dn]; !ok {
			m[dn] = pc
		}
	}

	result := make(util.PcSlice, 0, len(m))
	for _, pc := range m {
		result = append(result, pc)
	}

	sort.Sort(result)

	return result
}
