package collection

import (
	scv "github.com/NJUPT-ISL/SCV/api/v1"
	"github.com/NJUPT-ISL/Yoda-Scheduler/pkg/yoda/filter"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/scheduler/framework"
)

type Data struct {
	Value MaxValue
}

type MaxValue struct {
	MaxBandwidth   uint
	MaxClock       uint
	MaxCore        uint
	MaxFreeMemory  uint64
	MaxPower       uint
	MaxTotalMemory uint64
}

func (s *Data) Clone() framework.StateData {
	c := &Data{
		Value: s.Value,
	}
	return c
}

func CollectMaxValues(state *framework.CycleState, pod *v1.Pod, scvList scv.ScvList) *framework.Status {
	data := Data{Value: MaxValue{
		MaxBandwidth:   1,
		MaxClock:       1,
		MaxCore:        1,
		MaxFreeMemory:  1,
		MaxPower:       1,
		MaxTotalMemory: 1,
	}}
	for _, item := range scvList.Items {
		s := item.DeepCopy()
		if ok, number := filter.PodFitsNumber(pod, s); ok {
			isFitsMemory, memory := filter.PodFitsMemory(number, pod, s)
			isFitsClock, clock := filter.PodFitsClock(number, pod, s)
			if isFitsClock && isFitsMemory {
				for _, card := range s.Status.CardList {
					if card.FreeMemory >= memory && card.Clock >= clock {
						ProcessMaxValueWithCard(card, &data)
					}
				}
			}
		}
	}
	state.Write("Max", &data)
	return framework.NewStatus(framework.Success, "")
}

func ProcessMaxValueWithCard(card scv.Card, data *Data) {
	if card.FreeMemory > data.Value.MaxFreeMemory {
		data.Value.MaxFreeMemory = card.FreeMemory
	}
	if card.Clock > data.Value.MaxClock {
		data.Value.MaxClock = card.Clock
	}
	if card.TotalMemory > data.Value.MaxTotalMemory {
		data.Value.MaxTotalMemory = card.TotalMemory
	}
	if card.Bandwidth > data.Value.MaxBandwidth {
		data.Value.MaxBandwidth = card.Bandwidth
	}
	if card.Core > data.Value.MaxCore {
		data.Value.MaxCore = card.Core
	}
	if card.Power > data.Value.MaxPower {
		data.Value.MaxPower = card.Power
	}
}
