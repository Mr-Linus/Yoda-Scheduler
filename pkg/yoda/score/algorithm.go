package score

import (
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/collection"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/filter"
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
	"Memory":      1}

func CalculateSumValue(value string, state *framework.CycleState) int64 {
	d, err := state.Read(framework.StateKey("Max" + value))
	if err != nil {
		klog.V(3).Infof("Error Get CycleState Info: %v", err)
	}
	return d.(*collection.Data).Value
}

func CalculateValueScore(value string, state *framework.CycleState, node *nodeinfo.NodeInfo) (int64, error) {
	d, err := state.Read(framework.StateKey("Max" + value))
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
