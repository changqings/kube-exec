package util

import (
	"sort"
	"strings"

	corev1 "k8s.io/api/core/v1"
)

type deployNs struct {
	DeploymentName string
	NameSpace      string
}

func DeploymentWithOnePod(pcs PcSlice) PcSlice {

	m := make(map[deployNs]PodContainer)

	for _, pc := range pcs {
		dn := deployNs{pc.DeploymentName, pc.NameSpace}
		if _, ok := m[dn]; !ok {
			m[dn] = pc
		}
	}

	result := make(PcSlice, 0, len(m))
	for _, pc := range m {
		result = append(result, pc)
	}

	sort.Sort(result)

	return result
}

func GetPodDeployWithSplit(pods []corev1.Pod) PcSlice {

	pcs := PcSlice{}
	containerName := "app"

	for _, p := range pods {
		if len(p.GetOwnerReferences()) > 0 {
			if p.OwnerReferences[0].Kind == "ReplicaSet" {
				hash := p.GetLabels()["pod-template-hash"]
				dpname := strings.Split(p.GenerateName, "-"+hash+"-")[0]
				if checkContainerNameAndStat(&p, containerName) {
					pc := PodContainer{
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
