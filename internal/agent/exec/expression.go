package agentExec

import (
	"fmt"
	"github.com/Knetic/govaluate"
	valueUtils "github.com/aaronchen2k/deeptest/internal/agent/exec/utils/value"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/utils"
	"github.com/aaronchen2k/deeptest/pkg/lib/string"
	"regexp"
	"strings"
)

var (
	GovaluateFunctions = map[string]govaluate.ExpressionFunction{
		"length": func(args ...interface{}) (interface{}, error) {
			length := len(args[0].(string))
			return (float64)(length), nil
		},
		"match": func(args ...interface{}) (interface{}, error) {
			str := args[0].(string)
			regx := regexp.MustCompile("args[1].(string)")

			ret := regx.MatchString(str)
			return ret, nil
		},
		"uuid": func(args ...interface{}) (interface{}, error) {
			ret := _stringUtils.Uuid()
			return ret, nil
		},
	}
)

// called by checkpoint
func EvaluateGovaluateExpressionWithDebugVariables(expression string) (ret interface{}, err error) {
	expr := commUtils.RemoveLeftVariableSymbol(expression)

	govaluateExpression, err := govaluate.NewEvaluableExpressionWithFunctions(expr, GovaluateFunctions)
	if err != nil {
		return
	}

	// 1
	paramValMap, err := generateGovaluateParamsWithVariables(expression)
	if err != nil {
		return
	}

	ret, err = govaluateExpression.Evaluate(paramValMap)

	return
}

// called by agent processor interface
func EvaluateGovaluateExpressionByProcessorScope(expression string, scopeId uint) (ret interface{}, err error) {
	expr := commUtils.RemoveLeftVariableSymbol(expression)

	valueExpression, err := govaluate.NewEvaluableExpressionWithFunctions(expr, GovaluateFunctions)
	if err != nil {
		ret = expression
		return
	}

	// 1
	parameters, err := generateGovaluateParamsByScope(expression, scopeId)
	if err != nil {
		return
	}

	ret, err = valueExpression.Evaluate(parameters)

	return
}

// a.1
func generateGovaluateParamsByScope(expression string, scopeId uint) (ret domain.VarKeyValuePair, err error) {
	ret = make(map[string]interface{}, 8)

	variables := commUtils.GetVariablesInExpressionPlaceholder(expression)

	for _, varName := range variables {
		varNameWithoutPlus := strings.TrimLeft(varName, "+")

		var vari domain.ExecVariable
		vari, err = GetVariable(scopeId, varNameWithoutPlus)
		variValueStr := valueUtils.InterfaceToStr(vari.Value)

		if err == nil {
			var val interface{}
			if strings.Index(varName, "+") == 0 { // is a number like ${+id}
				val = _stringUtils.ParseInt(variValueStr)
			} else {
				val = variValueStr
			}

			ret[varNameWithoutPlus] = val
		}
	}

	return
}

// a.2
func generateGovaluateParamsWithVariables(expression string) (ret map[string]interface{}, err error) {
	ret = make(map[string]interface{}, 0)

	variables := commUtils.GetVariablesInExpressionPlaceholder(expression)

	for _, varName := range variables {
		varNameWithoutPlus := strings.TrimLeft(varName, "+")

		vari, _ := GetVariable(CurrScenarioProcessorId, varNameWithoutPlus)
		variValueStr := valueUtils.InterfaceToStr(vari.Value)

		var val interface{}
		if strings.Index(varName, "+") == 0 { // is a number like ${+id}
			val = _stringUtils.ParseInt(variValueStr)
		} else {
			val = variValueStr
		}

		ret[varNameWithoutPlus] = val
	}

	return
}

func ReplaceDatapoolVariInGovaluateExpress(expression string) (ret string) {
	ret = expression
	variablePlaceholders := commUtils.GetVariablesInExpressionPlaceholder(expression)

	for _, placeholder := range variablePlaceholders {
		if strings.Index(placeholder, "_dp") != 0 && strings.Index(placeholder, "_dp") != 1 {
			continue
		}

		oldVal := fmt.Sprintf("${%s}", placeholder)

		placeholderWithoutPlus := strings.TrimLeft(placeholder, "+")
		newVal := getPlaceholderVariableValue(placeholderWithoutPlus)
		if strings.Index(placeholder, "+") != 0 {
			newVal = "'" + newVal + "'"
		}

		ret = strings.ReplaceAll(ret, oldVal, newVal)
	}

	// add space to replace a==-1 to a== -1
	ret = strings.ReplaceAll(ret, "==-", "== -")

	return
}
