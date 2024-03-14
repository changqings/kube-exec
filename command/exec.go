package command

import (
	"bufio"
	"bytes"
	"context"
	"kube-exec/util"
	"log/slog"
	"net/http"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/remotecommand"
)

type osVersion struct {
	id      string
	version string
}

func GetRelease(config *rest.Config, cs *kubernetes.Clientset, pcs util.PcSlice) {

	pcsUinq := util.DeploymentWithOnePod(pcs)

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
			slog.Error("exec NewSPDYExecutor()", "msg", err)
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
			slog.Error("exec StreamWithContext()", "msg", err)
			continue
		}

		os := osVersion{}

		scan := bufio.NewScanner(&stdout)
		for scan.Scan() {
			if strings.HasPrefix(scan.Text(), "ID=") {
				os.id = strings.Split(scan.Text(), "ID=")[1]
			}
			if strings.HasPrefix(scan.Text(), "VERSION_ID=") {
				v := strings.Split(scan.Text(), "VERSION_ID=")[1]
				os.version = strings.Trim(v, "\"")
			}
		}
		slog.Info("Get pod os_version:", "ns", pc.NameSpace, "deployment", pc.DeploymentName, "pod", pc.PodName, "os", os.id, "version", os.version)
	}
}
