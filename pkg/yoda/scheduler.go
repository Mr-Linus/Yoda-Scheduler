package yoda

import (
	"context"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/collection"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/filter"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/sort"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

const (
	Name  = "yoda"
)

var (
	_ framework.FilterPlugin = &Yoda{}
	_ framework.PostFilterPlugin = &Yoda{}
	_ framework.QueueSortPlugin = &Yoda{}
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
		args: args,
		handle: f,
	}, nil
}

func (y *Yoda) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, node *nodeinfo.NodeInfo) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v", pod.Name, node.Node().Name)
	if ok,msg := filter.CheckGPUHealth(node);ok{
		if !filter.PodFitsLevel(pod,node){
				return framework.NewStatus(framework.Unschedulable, "Node:"+node.Node().Name+" GPU Level Not Fit")
		}
		if !filter.PodFitsMemory(pod,node){
			return framework.NewStatus(framework.Unschedulable, "Node:"+node.Node().Name+" GPU Memory Not Fit")
		}
		if !filter.PodFitsNumber(pod,node){
			return framework.NewStatus(framework.Unschedulable,"Node:"+node.Node().Name+" GPU Number Not Fit")
		}
		return framework.NewStatus(framework.Success, "")
	}else {
		return framework.NewStatus(framework.Unschedulable, "Node:"+node.Node().Name+msg)
	}

}

func (y *Yoda) PostFilter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodes []*v1.Node, filteredNodesStatuses framework.NodeToStatusMap) *framework.Status {
	klog.V(3).Infof("collect info for scheduling  pod: %v", pod.Name)
	return collection.ParallelCollection(state,nodes,filteredNodesStatuses)
}

func (y *Yoda) Less(podInfo1, podInfo2 *framework.PodInfo) bool{
	return sort.Less(podInfo1,podInfo2)
}

