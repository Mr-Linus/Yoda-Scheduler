package score

import (
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/collection"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/filter"
	v1 "k8s.io/api/core/v1"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// Sum is from collection/collection.go
// var Sum = []string{"Cores","FreeMemory","Bandwidth","MemoryClock","MemorySum","Number","Memory"}

var Weights = map[string]int64{
	"Cores":       2,
	"FreeMemory":  9,
	"Bandwidth":   3,
	"MemoryClock": 2,
	"MemorySum":   1,
	"Number":      1,
	"Memory":      1,
}

func CalculateValueScore(value string, state *framework.CycleState, node *nodeinfo.NodeInfo) (int64, error) {
	state.RLock()
	d, err := state.Read(framework.StateKey("Max" + value))
	state.RUnlock()
	if err != nil {
		klog.V(3).Infof("Error Get CycleState Info: %v", err)
		return 0, err
	}
	return filter.StrToInt64(node.Node().Labels["scv/"+value]) * 100 / d.(*collection.Data).Value, nil
}

func CalculateCollectScore(state *framework.CycleState, node *nodeinfo.NodeInfo) (int64, error) {
	var score int64 = 0
	for v, w := range Weights {
		s, err := CalculateValueScore(v, state, node)
		if err != nil {
			return 0, err
		}
		score += s * w
	}
	return score, nil
}

func CalculatePodUseScore(node *nodeinfo.NodeInfo) int64 {
	var score  = filter.StrToInt64(node.Node().GetLabels()["scv/Memory"])
	var memSum int64 = 0
	for _, pod := range node.Pods(){
		if mem,ok := pod.GetLabels()["scv/FreeMemory"];ok{
			if pod.Status.Phase != v1.PodSucceeded{
				memSum += filter.StrToInt64(mem)
			}
		}
	}
	score -= memSum
	return score
}