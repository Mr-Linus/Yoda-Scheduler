package yoda

import (
	"context"
	"fmt"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/collection"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/filter"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/score"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/sort"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

const (
	Name = "yoda"
)

var (
	_ framework.QueueSortPlugin  = &Yoda{}
	_ framework.FilterPlugin     = &Yoda{}
	_ framework.PostFilterPlugin = &Yoda{}
	_ framework.ScorePlugin      = &Yoda{}
	_ framework.ScoreExtensions  = &Yoda{}
)

type Args struct {
	KubeConfig string `json:"kubeconfig,omitempty"`
	Master     string `json:"master,omitempty"`
}

type Yoda struct {
	args   *Args
	handle framework.FrameworkHandle
}

func (y *Yoda) Name() string {
	return Name
}

func New(configuration *runtime.Unknown, f framework.FrameworkHandle) (framework.Plugin, error) {
	args := &Args{}
	if err := framework.DecodeInto(configuration, args); err != nil {
		return nil, err
	}
	klog.V(3).Infof("get plugin config args: %+v", args)
	return &Yoda{
		args:   args,
		handle: f,
	}, nil
}

func (y *Yoda) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, node *nodeinfo.NodeInfo) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v", pod.Name, node.Node().Name)
	if ok, msg := filter.CheckGPUHealth(node); ok {
		if !filter.PodFitsLevel(pod, node) {
			return framework.NewStatus(framework.Unschedulable, "Node:"+node.Node().Name+" GPU Level Not Fit")
		}
		if !filter.PodFitsMemory(pod, node) {
			return framework.NewStatus(framework.Unschedulable, "Node:"+node.Node().Name+" GPU Memory Not Fit")
		}
		if !filter.PodFitsNumber(pod, node) {
			return framework.NewStatus(framework.Unschedulable, "Node:"+node.Node().Name+" GPU Number Not Fit")
		}
		return framework.NewStatus(framework.Success, "")
	} else {
		return framework.NewStatus(framework.Unschedulable, "Node:"+node.Node().Name+msg)
	}

}

func (y *Yoda) PostFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodes []*v1.Node, filteredNodesStatuses framework.NodeToStatusMap) *framework.Status {
	klog.V(3).Infof("collect info for scheduling  pod: %v", pod.Name)
	return collection.ParallelCollection(collection.Workers, state, nodes, filteredNodesStatuses)
}

func (y *Yoda) Less(podInfo1, podInfo2 *framework.PodInfo) bool {
	return sort.Less(podInfo1, podInfo2)
}

func (y *Yoda) Score(ctx context.Context, state *framework.CycleState, p *v1.Pod, nodeName string) (int64, *framework.Status) {
	nodeInfo, err := y.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
	}
	s, err := score.Score(state, nodeInfo)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("Score Node Error: %v", err))
	}
	klog.V(3).Infof("node : %v yoda-score: %v",nodeName,s)
	return s, framework.NewStatus(framework.Success, "")
}

func (y *Yoda) NormalizeScore(ctx context.Context, state *framework.CycleState, p *v1.Pod, scores framework.NodeScoreList) *framework.Status {
	var (
		highest int64 = 0
		lowest = scores[0].Score
	)
	for _, nodeScore := range scores {
		if nodeScore.Score < lowest {
			lowest = nodeScore.Score
		}
	}
	if lowest < 0 {
		for i := range scores {
			scores[i].Score -= lowest
		}
	}
	for _, nodeScore := range scores {
		if nodeScore.Score > highest {
			highest = nodeScore.Score
		}
	}
	// Set Range to [0-100]
	for i, nodeScore := range scores {
		scores[i].Score = nodeScore.Score * framework.MaxNodeScore / highest
	}
	return framework.NewStatus(framework.Success, "")
}

func (y *Yoda) ScoreExtensions() framework.ScoreExtensions {
	return y
}
