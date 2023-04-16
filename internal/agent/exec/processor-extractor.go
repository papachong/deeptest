package agentExec

import (
	"encoding/json"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/domain"
	"github.com/aaronchen2k/deeptest/internal/agent/exec/utils/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	extractorHelper "github.com/aaronchen2k/deeptest/internal/pkg/helper/extractor"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"time"
)

type ProcessorExtractor struct {
	ID uint `json:"id" yaml:"id"`
	ProcessorEntityBase

	domain.ExtractorBase

	InterfaceID uint `json:"interfaceID"`

	Disabled bool `json:"disabled"`
}

func (entity ProcessorExtractor) Run(processor *Processor, session *Session) (err error) {
	logUtils.Infof("extractor entity")

	startTime := time.Now()
	processor.Result = &agentDomain.ScenarioExecResult{
		ID:                int(entity.ProcessorID),
		Name:              entity.Name,
		ProcessorCategory: entity.ProcessorCategory,
		ProcessorType:     entity.ProcessorType,
		StartTime:         &startTime,
		ParentId:          int(entity.ParentID),
	}

	brother, ok := getPreviousBrother(*processor)
	if !ok || brother.EntityType != consts.ProcessorInterfaceDefault {
		processor.Result.Summary = fmt.Sprintf("先前节点不是接口，无法应用提取器。")
		processor.AddResultToParent()
		execUtils.SendExecMsg(*processor.Result, session.WsMsg)
		return
	}

	resp := v1.DebugResponse{}
	json.Unmarshal([]byte(brother.Result.RespContent), &resp)

	entity.Src = consts.Body
	entity.Type = getExtractorTypeForProcessor(entity.ProcessorType)

	entity.Result, err = extractorHelper.Extract(entity.ExtractorBase, resp)
	if err != nil {
		processor.Result.Summary = fmt.Sprintf("%s提取器解析错误 %s。", entity.ProcessorType, err.Error())
		processor.AddResultToParent()
		execUtils.SendExecMsg(*processor.Result, session.WsMsg)
		return
	}

	SetVariable(processor.ParentId, entity.Variable, entity.Result, consts.Public) // set in parent scope

	processor.Result.Summary = fmt.Sprintf("将结果\"%v\"赋予变量\"%s\"。", entity.Result, entity.Variable)
	processor.AddResultToParent()
	execUtils.SendExecMsg(*processor.Result, session.WsMsg)

	endTime := time.Now()
	processor.Result.EndTime = &endTime

	return
}

func getExtractorTypeForProcessor(processorType consts.ProcessorType) (ret consts.ExtractorType) {
	if processorType == consts.ProcessorExtractorBoundary {
		ret = consts.Boundary
	} else if processorType == consts.ProcessorExtractorJsonQuery {
		ret = consts.JsonQuery
	} else if processorType == consts.ProcessorExtractorHtmlQuery {
		ret = consts.HtmlQuery
	} else if processorType == consts.ProcessorExtractorXmlQuery {
		ret = consts.XmlQuery
	}

	return
}
