package filter

import (
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
	"strconv"
)

//TODO: PodFitsGPUNumber(NodeHasGPUNumber+PodHasGPUNumber)

func CheckGPUHealth(node *nodeinfo.NodeInfo) (bool, string){
	var msg = ""
	if NodeHasGPU(node){
		if NodeGPUHealth(node){
			return true, msg
		}
		return false, "GPU Unhealthy"
	}
	return false, "No GPU"
}

func NodeHasGPUNumber(node *nodeinfo.NodeInfo) bool {
	_,ok := node.Node().Labels["scv/Number"]
	return ok
}

func NodeHasGPU(node *nodeinfo.NodeInfo) bool{
	if _,ok := node.Node().Labels["scv/Gpu"];ok{
		if node.Node().Labels["scv/Gpu"] == "True" {
			return true
		}
	}
	return false
}

func NodeHasLevel(node *nodeinfo.NodeInfo) bool{
	_,ok := node.Node().Labels["scv/Level"]
	return ok
}

func NodeHasFreeMemory(node *nodeinfo.NodeInfo) bool{
	_,ok := node.Node().Labels["scv/FreeMemory"]
	return ok
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
	_,ok := pod.Labels["scv/FreeMemory"]
	return ok
}

func PodNeedGPUNumber(pod *v1.Pod) bool {
	_,ok := pod.Labels["scv/Number"]
	return ok
}

func PodFitsMemory(pod *v1.Pod,node *nodeinfo.NodeInfo) bool {
	if PodNeedMemory(pod){
		if NodeHasFreeMemory(node){
			return StrToUInt(node.Node().Labels["scv/FreeMemory"]) >= StrToUInt(pod.Labels["scv/FreeMemory"])
		}
		return false
	}
	return true
}

func PodFitsLevel(pod *v1.Pod,node *nodeinfo.NodeInfo) bool{
	if PodNeedLevel(pod){
		if NodeHasLevel(node){
			return GetLevel(node.Node().Labels["scv/Level"]) >= GetLevel(pod.Labels["scv/Level"])
		}
		return false
	}
	return true
}

func PodFitsNumber(pod *v1.Pod,node *nodeinfo.NodeInfo) bool {
	if PodNeedGPUNumber(pod){
		if NodeHasGPUNumber(node){
			return 	StrToUInt(node.Node().Labels["scv/Number"]) >= StrToUInt(pod.Labels["scv/Number"])
		}
		return false
	}
	return true
}

func GetLevel(label string) int{
	var level = 0
	switch label {
		case "High": level = 3
		case "Medium": level = 2
		case "Low": level = 1
	}
	return level
}

func StrToUInt(str string) uint {
	if i, e := strconv.Atoi(str);e != nil {
		return 0
	}else {
		return uint(i)
	}
}

func StrToInt64(str string) int64 {
	if i, e := strconv.Atoi(str);e != nil {
		return 0
	}else {
		return int64(i)
	}
}