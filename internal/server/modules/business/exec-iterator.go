package business

import (
	"container/list"
	"fmt"
	agentDomain "github.com/aaronchen2k/deeptest/internal/agent/exec"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/helper/exec"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	_logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
)

var (
	IteratorContextValueStack *list.List
	IteratorContextIndexStack *list.List
)

type ExecIterator struct {
	ScenarioProcessorRepo *repo.ScenarioProcessorRepo `inject:""`
}

func (s *ExecIterator) InitIteratorContext() {
	IteratorContextValueStack = list.New()
	IteratorContextIndexStack = list.New()
}

func (s *ExecIterator) Push(iterator domain.ExecIterator) {
	IteratorContextValueStack.PushFront(iterator)
	IteratorContextIndexStack.PushFront(0)

	return
}

func (s *ExecIterator) Pop() {
	IteratorContextValueStack.Remove(IteratorContextValueStack.Front())
	IteratorContextIndexStack.Remove(IteratorContextIndexStack.Front())

	return
}

func (s *ExecIterator) GenerateLoopTimes(log domain.ExecLog) (ret domain.ExecIterator, err error) {
	ret.ProcessorCategory = log.ProcessorCategory
	ret.ProcessorType = log.ProcessorType

	if log.Output.Times > 0 {
		for i := 0; i < log.Output.Times; i++ {
			ret.Times = append(ret.Times, i+1)
		}
	}

	return
}

func (s *ExecIterator) GenerateLoopRange(log domain.ExecLog, stepStr string, isRand bool) (ret domain.ExecIterator, err error) {
	ret.ProcessorCategory = log.ProcessorCategory
	ret.ProcessorType = log.ProcessorType

	start, end, step, precision, typ, err := execHelper.GetRange(log.Output, stepStr)
	if err == nil {
		ret.DataType = typ
		ret.Items, _ = execHelper.GenerateRangeItems(start, end, step, precision, isRand, typ)
	}

	return
}
func (s *ExecIterator) GenerateLoopList(log domain.ExecLog) (ret domain.ExecIterator, err error) {
	ret.ProcessorCategory = log.ProcessorCategory
	ret.ProcessorType = log.ProcessorType

	ret.Items, ret.DataType, err = execHelper.GenerateListItems(log.Output.List)

	return
}

func (s *ExecIterator) GenerateData(log domain.ExecLog, data model.ProcessorData) (ret domain.ExecIterator, err error) {
	ret.ProcessorCategory = log.ProcessorCategory
	ret.ProcessorType = log.ProcessorType

	ret.Data, _ = execHelper.GenerateDataItems(data)

	return
}

func (s *ExecIterator) RetrieveIteratorsVal(processor *agentDomain.Processor) (item interface{}, desc string, err error) {
	valueElem := IteratorContextValueStack.Front()
	indexElem := IteratorContextIndexStack.Front()
	if valueElem == nil || indexElem == nil {
		return
	}

	value := valueElem.Value.(domain.ExecIterator)
	index := indexElem.Value.(int)

	if value.ProcessorType == consts.ProcessorLoopTime {
		items := value.Times

		if index > len(items)-1 {
			index = 0
		}
		if len(items) == 0 {
			return
		}

		item = items[index]
		desc = fmt.Sprintf("(%v / %d)", item, len(items))

	} else if value.ProcessorType == consts.ProcessorLoopRange {
		items := value.Items
		if index > len(items)-1 {
			index = 0
		}
		if len(items) == 0 {
			return
		}

		item = items[index]

		po, _ := s.ScenarioProcessorRepo.Get(processor.ID)
		loopRangeProcessor, _ := s.ScenarioProcessorRepo.GetLoop(po)
		desc = fmt.Sprintf("变量%s = \"%d\"", loopRangeProcessor.VariableName, item)

	} else if value.ProcessorType == consts.ProcessorLoopIn {
		items := value.Items
		if index > len(items)-1 {
			index = 0
		}
		if len(items) == 0 {
			return
		}

		item = items[index]

		po, _ := s.ScenarioProcessorRepo.Get(processor.ID)
		loopRangeProcessor, _ := s.ScenarioProcessorRepo.GetLoop(po)
		desc = fmt.Sprintf("变量%s = \"%d\"", loopRangeProcessor.VariableName, item)

	} else if value.ProcessorType == consts.ProcessorLoopBreak {
		_logUtils.Info("")

	}

	index++
	indexElem.Value = index

	return
}
