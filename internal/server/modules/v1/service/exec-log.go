package service

import (
	"encoding/json"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	serverDomain "github.com/aaronchen2k/deeptest/internal/server/modules/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/v1/repo"
	"strings"
)

type ExecLogService struct {
	ScenarioProcessorRepo *repo.ScenarioProcessorRepo `inject:""`
	ScenarioRepo          *repo.ScenarioRepo          `inject:""`
	TestResultRepo        *repo.ReportRepo            `inject:""`
	TestLogRepo           *repo.LogRepo               `inject:""`
	InterfaceRepo         *repo.InterfaceRepo         `inject:""`
	InterfaceService      *InterfaceService           `inject:""`
}

func (s *ExecLogService) CreateProcessorLog(processor *model.Processor, log *domain.ExecLog, parentPersistentId uint) (po model.ExecLogProcessor, err error) {
	po = model.ExecLogProcessor{
		Name:              processor.Name,
		ProcessorCategory: processor.EntityCategory,
		ProcessorType:     processor.EntityType,
		ProcessorId:       processor.ID,

		ParentId: parentPersistentId,
		ReportId: log.ReportId,
	}

	po.Summary = strings.Join(log.Summary, "; ")

	outputBytes, _ := json.Marshal(log.Output)
	po.Output = string(outputBytes)

	err = s.TestLogRepo.Save(&po)
	log.Id = po.ID
	log.PersistentId = po.ID

	return
}

func (s *ExecLogService) CreateInterfaceLog(req serverDomain.InvocationRequest, resp serverDomain.InvocationResponse, parentLog *domain.ExecLog) (
	po model.ExecLogProcessor, err error) {
	po = model.ExecLogProcessor{
		Name:              req.Name,
		ProcessorCategory: consts.ProcessorInterface,
		ProcessorType:     consts.ProcessorInterfaceDefault,
		ResultStatus:      consts.Pass, // TODO:
		InterfaceId:       req.Id,

		ParentId: parentLog.PersistentId,
		ReportId: parentLog.ReportId,
	}

	bytesReq, _ := json.Marshal(req)
	po.ReqContent = string(bytesReq)

	bytesReps, _ := json.Marshal(resp)
	po.RespContent = string(bytesReps)

	err = s.TestLogRepo.Save(&po)

	return
}
