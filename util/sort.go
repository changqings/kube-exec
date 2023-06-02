package util

type PodContainer struct {
	DeploymentName string
	NameSpace      string
	PodName        string
	ContainerName  string
}
type PcSlice []PodContainer

func (pcs PcSlice) Len() int {
	return len(pcs)
}

func (pcs PcSlice) Less(i, j int) bool {
	return pcs[i].NameSpace < pcs[j].NameSpace
}

func (pcs PcSlice) Swap(i, j int) {
	pcs[i], pcs[j] = pcs[j], pcs[i]
}
