package sort

import (
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"strconv"
)

func Less(podInfo1, podInfo2 *framework.PodInfo) bool {
	return GetPodPriority(podInfo1) > GetPodPriority(podInfo2)
}

func GetPodPriority(podInfo *framework.PodInfo) int {
	if p, ok := podInfo.Pod.Labels["scv/Priority"]; ok {
		pri, _ := strconv.Atoi(p)
		return pri
	}
	return 0
}
