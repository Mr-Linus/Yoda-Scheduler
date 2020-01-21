package yoda

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
	"strconv"
)

func NodeHasGPU(node *nodeinfo.NodeInfo) bool{
	if _,ok := node.Node().Labels["scv/Gpu"];ok{
		if node.Node().Labels["scv/Gpu"] == "True" {
			return true
		}
	}
	return false
}

func NodeGPUHealth(node *nodeinfo.NodeInfo) bool{
	if node.Node().Labels["scv/Health"] == "Healthy" {
		return true
	}
	return false
}

func PodNeedLevel(pod *v1.Pod) bool{
	if _,ok := pod.Labels["scv/Level"];ok {
		return true
	}
	return false
}

func PodNeedMemory(pod *v1.Pod) bool{
	if _,ok := pod.Labels["scv/FreeMemory"];ok {
		return true
	}
	return false
}

func PodFitsMemory(pod *v1.Pod,node *nodeinfo.NodeInfo) bool {
	if StrToUInt(node.Node().Labels["scv/FreeMemory"]) >= StrToUInt(pod.Labels["scv/FreeMemory"]){
		return true
	}
	return false
}

func PodFitsLevel(pod *v1.Pod,node *nodeinfo.NodeInfo) bool{
	var (
		podLevel = 0
		nodeLevel = 0
	)
	switch pod.Labels["scv/Level"] {
		case "High": podLevel = 3
		case "Medium": podLevel = 2
		case "Low": podLevel = 1
	}
	switch node.Node().Labels["scv/Level"] {
		case "High": nodeLevel = 3
		case "Medium": nodeLevel = 2
		case "Low": nodeLevel = 1
	}
	if nodeLevel >= podLevel{
		return true
	}
	return false
}

func StrToUInt(str string) uint {
	i, e := strconv.Atoi(str)
	if e != nil {
		return 0
	}
	return uint(i)
}