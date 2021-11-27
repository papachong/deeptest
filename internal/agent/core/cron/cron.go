package agentCron

import (
	"fmt"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/cron"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/date"
	"github.com/aaronchen2k/deeptest/internal/pkg/lib/log"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/kataras/iris/v12"
	"sync"
	"time"
)

type AgentCron struct {
	syncMap sync.Map
}

func NewAgentCron() *AgentCron {
	inst := &AgentCron{}
	return inst
}

func (s *AgentCron) Init() {
	s.syncMap.Store("isRunning", false)
	s.syncMap.Store("lastCompletedTime", int64(0))

	cronUtils.AddTask(
		"check",
		fmt.Sprintf("@every %ds", serverConsts.WebCheckInterval),
		func() {
			isRunning, _ := s.syncMap.Load("isRunning")
			lastCompletedTime, _ := s.syncMap.Load("lastCompletedTime")

			if isRunning.(bool) || time.Now().Unix()-lastCompletedTime.(int64) < serverConsts.WebCheckInterval {
				logUtils.Infof("skip this iteration " + dateUtils.DateTimeStr(time.Now()))
				return
			}

			s.syncMap.Store("isRunning", true)

			// do somethings

			s.syncMap.Store("isRunning", false)
			s.syncMap.Store("lastCompletedTime", time.Now().Unix())
		},
	)

	iris.RegisterOnInterrupt(func() {
		cronUtils.Stop()
	})
}
