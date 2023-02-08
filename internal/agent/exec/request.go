package agentExec

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	httpHelper "github.com/aaronchen2k/deeptest/internal/pkg/helper/http"
	_httpUtils "github.com/aaronchen2k/deeptest/pkg/lib/http"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Invoke(req v1.BaseRequest) (resp v1.InvocationResponse, err error) {
	req.Url, err = _httpUtils.AddDefaultUrlSchema(req.Url)
	if err != nil {
		return
	}

	if req.Method == consts.GET {
		resp, err = httpHelper.Get(req)
	} else if req.Method == consts.POST {
		resp, err = httpHelper.Post(req)
	} else if req.Method == consts.PUT {
		resp, err = httpHelper.Put(req)
	} else if req.Method == consts.DELETE {
		resp, err = httpHelper.Delete(req)
	} else if req.Method == consts.PATCH {
		resp, err = httpHelper.Patch(req)
	} else if req.Method == consts.HEAD {
		resp, err = httpHelper.Head(req)
	} else if req.Method == consts.CONNECT {
		resp, err = httpHelper.Connect(req)
	} else if req.Method == consts.OPTIONS {
		resp, err = httpHelper.Options(req)
	} else if req.Method == consts.TRACE {
		resp, err = httpHelper.Trace(req)
	}

	return
}

func GetRequestProps(req *v1.BaseRequest) {
	req.BodyLang = consts.LangTEXT

	arr := strings.Split(string(req.BodyType), "/")
	if len(arr) == 1 {
		return
	}

	typeName := arr[1]
	if typeName == "text" || typeName == "plain" {
		typeName = consts.LangTEXT.String()
	}

	req.BodyLang = consts.HttpRespLangType(typeName)
}

func GetContentProps(resp *v1.InvocationResponse) {
	resp.ContentLang = consts.LangTEXT

	if resp.ContentLang == "" {
		return
	}

	arr := strings.Split(string(resp.ContentType), ";")

	arr1 := strings.Split(arr[0], "/")
	if len(arr1) == 1 {
		return
	}

	typeName := arr1[1]
	if typeName == "text" || typeName == "plain" {
		typeName = "plaintext"
	}
	resp.ContentLang = consts.HttpRespLangType(typeName)

	if len(arr) > 1 {
		arr2 := strings.Split(arr[1], "=")
		if len(arr2) > 1 {
			resp.ContentCharset = consts.HttpRespCharset(arr2[1])
		}
	}

	//ret.NodeContent = mockHelper.FormatXml(ret.NodeContent)

	return
}

func ReplaceAll(req *v1.BaseRequest, variableMap domain.Variables, datapools domain.Datapools) {
	replaceUrl(req, variableMap, datapools)
	replaceParams(req, variableMap, datapools)
	replaceHeaders(req, variableMap, datapools)
	replaceBody(req, variableMap, datapools)
	replaceAuthor(req, variableMap, datapools)
}

func replaceUrl(req *v1.BaseRequest, variableMap domain.Variables, datapools domain.Datapools) {
	req.Url = ReplaceVariableValue(req.Url, variableMap, datapools)
}
func replaceParams(req *v1.BaseRequest, variableMap domain.Variables, datapools domain.Datapools) {
	for idx, param := range req.Params {
		req.Params[idx].Value = ReplaceVariableValue(param.Value, variableMap, datapools)
	}
}
func replaceHeaders(req *v1.BaseRequest, variableMap domain.Variables, datapools domain.Datapools) {
	for idx, header := range req.Headers {
		req.Headers[idx].Value = ReplaceVariableValue(header.Value, variableMap, datapools)
	}
}
func replaceBody(req *v1.BaseRequest, variableMap domain.Variables, datapools domain.Datapools) {
	req.Body = ReplaceVariableValue(req.Body, variableMap, datapools)
}
func replaceAuthor(req *v1.BaseRequest, variableMap domain.Variables, datapools domain.Datapools) {
	if req.AuthorizationType == consts.BasicAuth {
		req.BasicAuth.Username = ReplaceVariableValue(req.BasicAuth.Username, variableMap, datapools)
		req.BasicAuth.Password = ReplaceVariableValue(req.BasicAuth.Password, variableMap, datapools)

	} else if req.AuthorizationType == consts.BearerToken {
		req.BearerToken.Token = ReplaceVariableValue(req.BearerToken.Token, variableMap, datapools)

	} else if req.AuthorizationType == consts.OAuth2 {
		req.OAuth20.Name = ReplaceVariableValue(req.OAuth20.Name, variableMap, datapools)
		req.OAuth20.CallbackUrl = ReplaceVariableValue(req.OAuth20.CallbackUrl, variableMap, datapools)
		req.OAuth20.AuthURL = ReplaceVariableValue(req.OAuth20.AuthURL, variableMap, datapools)
		req.OAuth20.AccessTokenURL = ReplaceVariableValue(req.OAuth20.AccessTokenURL, variableMap, datapools)
		req.OAuth20.ClientID = ReplaceVariableValue(req.OAuth20.ClientID, variableMap, datapools)
		req.OAuth20.Scope = ReplaceVariableValue(req.OAuth20.Scope, variableMap, datapools)

	} else if req.AuthorizationType == consts.ApiKey {
		req.ApiKey.Key = ReplaceVariableValue(req.ApiKey.Key, variableMap, datapools)
		req.ApiKey.Value = ReplaceVariableValue(req.ApiKey.Value, variableMap, datapools)
		req.ApiKey.TransferMode = ReplaceVariableValue(req.ApiKey.TransferMode, variableMap, datapools)
	}
}

func ReplaceVariableValue(value string, variableMap domain.Variables, datapools domain.Datapools) (ret string) {

	variablePlaceholders := GetVariablesInVariablePlaceholder(value)
	ret = value

	for _, placeholder := range variablePlaceholders {
		variablePlaceholder := fmt.Sprintf("${%s}", placeholder)

		oldVal := variablePlaceholder
		newVal := getPlaceholderValue(placeholder, variableMap, datapools)

		ret = strings.ReplaceAll(ret, oldVal, newVal)
	}

	return
}

func getPlaceholderValue(placeholder string, variableMap domain.Variables, datapools domain.Datapools) (ret string) {

	typ := getPlaceholderType(placeholder)

	if typ == consts.PlaceholderTypeVariable {
		ret = getVariableValue(placeholder, variableMap)

	} else if typ == consts.PlaceholderTypeDatapool {
		ret = getDatapoolValue(placeholder, datapools)

	} else if typ == consts.PlaceholderTypeFunction {
	}

	return
}

func getVariableValue(placeholder string, variableMap domain.Variables) (ret string) {
	ret = fmt.Sprintf("%v", variableMap[placeholder])
	return
}

func getDatapoolValue(placeholder string, datapools domain.Datapools) (ret string) {
	// _dp(name, col, <1 | seq | rand>)

	regex := regexp.MustCompile(fmt.Sprintf("(?Ui)%s\\((.+),(.+),(.+)\\)", consts.PlaceholderPrefixDatapool))
	arrs := regex.FindAllStringSubmatch(placeholder, -1)

	if !(len(arrs) == 1 && len(arrs[0]) == 4) {
		return
	}

	dpName := strings.TrimSpace(arrs[0][1])
	dpCol := strings.TrimSpace(arrs[0][2])
	dpSeq := strings.TrimSpace(arrs[0][3])

	dp := datapools[dpName]
	if dp == nil {
		ret = fmt.Sprintf("${%s}", placeholder)
		return
	}

	rowIndex := getDatapoolRow(dpName, dpSeq, datapools)

	val := datapools[dpName][rowIndex][dpCol]
	if val == nil {
		val = "NOT_FOUND"
	}

	ret = fmt.Sprintf("%v", val)

	return
}

func getDatapoolRow(dpName, seq string, datapools domain.Datapools) (ret int) {
	dp := datapools[dpName]
	if dp == nil {
		return
	}

	total := len(dp)

	if seq == "seq" {
		ret = DatapoolCursor[dpName] % total
		DatapoolCursor[dpName]++

	} else if seq == "rand" {
		rand.Seed(time.Now().Unix())
		ret = rand.Intn(total)

	} else {
		seqInt, _ := strconv.Atoi(seq)
		ret = seqInt % total
	}

	return
}

func getPlaceholderType(placeholder string) (ret consts.PlaceholderType) {
	if strings.HasPrefix(placeholder, consts.PlaceholderPrefixDatapool.String()) {
		return consts.PlaceholderTypeDatapool
	} else if strings.HasPrefix(placeholder, consts.PlaceholderPrefixFunction.String()) {
		return consts.PlaceholderTypeFunction
	}

	return consts.PlaceholderTypeVariable
}
