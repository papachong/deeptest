package service

import (
	"github.com/aaronchen2k/deeptest/internal/agent/exec"
	execDomain "github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"sync"
)

var (
	breakMap sync.Map
)

type ScenarioExecService struct {
	ScenarioRepo       *repo.ScenarioRepo       `inject:""`
	ScenarioNodeRepo   *repo.ScenarioNodeRepo   `inject:""`
	ScenarioReportRepo *repo.ScenarioReportRepo `inject:""`
	TestLogRepo        *repo.LogRepo            `inject:""`
	EnvironmentRepo    *repo.EnvironmentRepo    `inject:""`

	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	EndpointRepo          *repo.EndpointRepo          `inject:""`
	ServeServerRepo       *repo.ServeServerRepo       `inject:""`

	SceneService          *SceneService          `inject:""`
	EnvironmentService    *EnvironmentService    `inject:""`
	DatapoolService       *DatapoolService       `inject:""`
	ScenarioNodeService   *ScenarioNodeService   `inject:""`
	ScenarioReportService *ScenarioReportService `inject:""`
}

func (s *ScenarioExecService) LoadExecResult(scenarioId int) (result domain.Report, err error) {
	scenario, err := s.ScenarioRepo.Get(uint(scenarioId))
	if err != nil {
		return
	}

	result.Name = scenario.Name

	return
}

func (s *ScenarioExecService) LoadExecData(scenarioId, environmentId uint) (ret agentExec.ScenarioExecObj, err error) {
	scenario, err := s.ScenarioRepo.Get(scenarioId)
	if err != nil {
		return
	}

	// get processor tree
	ret.ScenarioId = scenarioId
	ret.Name = scenario.Name
	ret.RootProcessor, _ = s.ScenarioNodeService.GetTree(scenario, true)
	ret.RootProcessor.Session = agentExec.Session{}

	// get variables
	s.SceneService.LoadEnvVarMapByScenario(&ret.ExecScene, scenarioId, environmentId)
	s.SceneService.LoadProjectSettings(&ret.ExecScene, scenario.ProjectId)

	return
}

func (s *ScenarioExecService) SaveReport(scenarioId int, userId uint, rootResult execDomain.ScenarioExecResult) (report model.ScenarioReport, err error) {
	scenario, _ := s.ScenarioRepo.Get(uint(scenarioId))
	rootResult.Name = scenario.Name

	report = model.ScenarioReport{
		Name:      scenario.Name,
		StartTime: rootResult.StartTime,
		EndTime:   rootResult.EndTime,
		Duration:  rootResult.EndTime.UnixMilli() - rootResult.StartTime.UnixMilli(),

		ProgressStatus: rootResult.ProgressStatus,
		ResultStatus:   rootResult.ResultStatus,

		ScenarioId:   scenario.ID,
		ProjectId:    scenario.ProjectId,
		CreateUserId: userId,
		ExecEnvId:    rootResult.EnvironmentId,
	}

	s.countRequest(rootResult, &report)
	s.summarizeInterface(&report)

	s.ScenarioReportRepo.Create(&report)
	s.TestLogRepo.CreateLogs(rootResult, &report)
	logs, _ := s.ScenarioReportService.GetById(report.ID)
	report.Logs = logs.Logs
	report.Priority = scenario.Priority

	return
}

func (s *ScenarioExecService) GenerateReport(scenarioId int, userId uint, rootResult execDomain.ScenarioExecResult) (report model.ScenarioReport, err error) {
	scenario, _ := s.ScenarioRepo.Get(uint(scenarioId))
	rootResult.Name = scenario.Name

	report = model.ScenarioReport{
		Name:      scenario.Name,
		StartTime: rootResult.StartTime,
		EndTime:   rootResult.EndTime,
		Duration:  rootResult.EndTime.Unix() - rootResult.StartTime.Unix(),

		ProgressStatus: rootResult.ProgressStatus,
		ResultStatus:   rootResult.ResultStatus,

		ScenarioId:   scenario.ID,
		ProjectId:    scenario.ProjectId,
		CreateUserId: userId,
	}

	s.countRequest(rootResult, &report)
	s.summarizeInterface(&report)

	//s.ScenarioReportRepo.Create(&report)
	//s.TestLogRepo.CreateLogs(rootResult, &report)

	return
}

func (s *ScenarioExecService) countRequest(result execDomain.ScenarioExecResult, report *model.ScenarioReport) {
	report.TotalProcessorNum++
	report.FinishProcessorNum++
	if result.ProcessorType == consts.ProcessorInterfaceDefault {
		s.countInterface(result.InterfaceId, result.ResultStatus, report)

		report.TotalRequestNum++
		report.Duration += result.Cost

		switch result.ResultStatus {
		case consts.Pass:
			report.PassRequestNum++

		case consts.Fail:
			report.FailRequestNum++
			report.ResultStatus = consts.Fail

		default:
		}

	} else if result.ProcessorType == consts.ProcessorAssertionDefault {
		report.TotalAssertionNum++
		switch result.ResultStatus {
		case consts.Pass:
			report.PassAssertionNum++

		case consts.Fail:
			report.FailAssertionNum++
			report.ResultStatus = consts.Fail

		default:
		}
	}

	if result.Children == nil {
		return
	}

	for _, log := range result.Children {
		s.countRequest(*log, report)
	}
}

func (s *ScenarioExecService) countInterface(interfaceId uint, status consts.ResultStatus, report *model.ScenarioReport) {
	if report.InterfaceStatusMap == nil {
		report.InterfaceStatusMap = map[uint]map[consts.ResultStatus]int{}
	}

	if report.InterfaceStatusMap[interfaceId] == nil {
		report.InterfaceStatusMap[interfaceId] = map[consts.ResultStatus]int{}
		report.InterfaceStatusMap[interfaceId][consts.Pass] = 0
		report.InterfaceStatusMap[interfaceId][consts.Fail] = 0
	}

	switch status {
	case consts.Pass:
		report.InterfaceStatusMap[interfaceId][consts.Pass]++

	case consts.Fail:
		report.InterfaceStatusMap[interfaceId][consts.Fail]++

	default:
	}
}

func (s *ScenarioExecService) summarizeInterface(report *model.ScenarioReport) {
	for _, val := range report.InterfaceStatusMap {
		if val[consts.Fail] > 0 {
			report.FailInterfaceNum++
		} else {
			report.PassInterfaceNum++
		}

		report.TotalInterfaceNum++
	}
}

func (s *ScenarioExecService) GetScenarioNormalData(id, environmentId uint) (ret execDomain.Report, err error) {

	ret.ScenarioId = id

	environment, err := s.EnvironmentRepo.Get(environmentId)
	if err != nil {
		return
	}
	ret.ExecEnv = environment.Name

	scenarioIds := []uint{id}
	interfaceNum, err := s.ScenarioNodeRepo.GetNumberByScenariosAndEntityCategory(scenarioIds, "processor_interface")
	if err != nil {
		return ret, err
	}
	assertionNum, err := s.ScenarioNodeRepo.GetNumberByScenariosAndEntityCategory(scenarioIds, "processor_assertion")
	if err != nil {
		return ret, err
	}
	processorNum, err := s.ScenarioNodeRepo.GetNumberByScenariosAndEntityCategory(scenarioIds, "")
	if err != nil {
		return ret, err
	}

	ret.TotalInterfaceNum = int(interfaceNum)
	ret.TotalAssertionNum = int(assertionNum)
	ret.TotalProcessorNum = int(processorNum)

	return

}
