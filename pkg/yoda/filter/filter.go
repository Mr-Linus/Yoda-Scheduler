package filter

import (
	"strconv"

	v1 "k8s.io/api/core/v1"

	scv "github.com/NJUPT-ISL/SCV/api/v1"
)

func PodFitsNumber(pod *v1.Pod, scv *scv.Scv) (bool, uint) {
	if number, ok := pod.GetLabels()["scv/number"]; ok {
		return strToUint(number) <= scv.Status.CardNumber, strToUint(number)
	}
	return scv.Status.CardNumber > 0, 1
}

func PodFitsMemory(number uint, pod *v1.Pod, scv *scv.Scv) (bool, uint64) {
	if memory, ok := pod.GetLabels()["scv/memory"]; ok {
		fitsCard := uint(0)
		m := StrToUint64(memory)
		for _, card := range scv.Status.CardList {
			if CardFitsMemory(m, card) {
				fitsCard++
			}
		}
		if fitsCard >= number {
			return true, m
		}
		return false, m
	}
	return true, 0
}

func PodFitsClock(number uint, pod *v1.Pod, scv *scv.Scv) (bool, uint) {
	if clock, ok := pod.GetLabels()["scv/clock"]; ok {
		fitsCard := uint(0)
		c := strToUint(clock)
		for _, card := range scv.Status.CardList {
			if CardFitsClock(c, card) {
				fitsCard++
			}
		}
		if fitsCard >= number {
			return true, c
		}
		return false, c
	}
	return true, 0
}

func CardFitsMemory(memory uint64, card scv.Card) bool {
	return card.Health == "Healthy" && card.FreeMemory >= memory
}

func CardFitsClock(clock uint, card scv.Card) bool {
	return card.Health == "Healthy" && card.Clock == clock
}

func strToUint(str string) uint {
	if i, e := strconv.Atoi(str); e != nil {
		return 0
	} else {
		return uint(i)
	}
}

func StrToUint64(str string) uint64 {
	if i, e := strconv.Atoi(str); e != nil {
		return 0
	} else {
		return uint64(i)
	}
}

func StrToInt64(str string) int64 {
	if i, e := strconv.Atoi(str); e != nil {
		return 0
	} else {
		return int64(i)
	}
}

func Uint64ToInt64(intNum uint64) int64 {
	return StrToInt64(strconv.FormatUint(intNum, 10))
}
