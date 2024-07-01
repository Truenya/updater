package updater

import (
	"fmt"
	"time"

	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/statistic"
	"gitlab-dev.ispsystem.net/team/vm/vm_custom_updater/util"
)

type Reporter struct {
	curStage     int
	startTime    time.Time
	lastStageDur time.Duration
	progress     chan util.Progress
}

func (r *Reporter) setUp(progress chan util.Progress) error {
	if progress == nil {
		return fmt.Errorf("nil chan passed to setup method")
	}

	r.progress = progress
	r.startTime = time.Now()
	r.curStage = 1

	return nil
}

func (r *Reporter) report(message, service string, uType util.UpdaterType) {
	LastStageKey := GetStatisticKey(r.curStage, uType, message)

	cur, ok := statistic.Get(service, LastStageKey)
	if !ok {
		cur = r.lastStageDur
	}
	r.progress <- util.Progress{
		Message:      message,
		LastStageDur: cur,
		Elapsed:      time.Since(r.startTime),
		Est:          statistic.GetEstDuration(service, r.curStage, uType),
	}
}

const MovingAverageWindow = 2

func Type2String() func(util.UpdaterType) string {
	uType2String := map[util.UpdaterType]string{
		util.UpdateLocal:   "local",
		util.UpdateRemote:  "container",
		util.UpdateInPlace: "inplace",
	}

	return func(uType util.UpdaterType) string {
		return uType2String[uType]
	}
}

func GetStatisticKey(curStage int, uType util.UpdaterType, msg string) string {
	return fmt.Sprintf("%d_%s_%s", curStage, Type2String()(uType), msg)
}

func (r *Reporter) updateMeanStage(cur time.Duration, service string, uType util.UpdaterType, msg string) {
	r.lastStageDur = cur

	stageStatisticKey := GetStatisticKey(r.curStage, uType, msg)

	prev, ok := statistic.Get(service, stageStatisticKey)
	if ok {
		// Скользящее среднее с шагом в 2
		cur := (prev + cur) / MovingAverageWindow
		statistic.Set(service, stageStatisticKey, cur)
	} else {
		statistic.Set(service, stageStatisticKey, cur)
	}
	r.curStage++
}
