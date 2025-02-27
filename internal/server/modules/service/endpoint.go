package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/474420502/requests"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	curlHelper "github.com/aaronchen2k/deeptest/internal/pkg/helper/gcurl"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/openapi"
	schemaHelper "github.com/aaronchen2k/deeptest/internal/pkg/helper/schema"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
	_commUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"github.com/getkin/kin-openapi/openapi3"
	encoder "github.com/zwgblue/yaml-encoder"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"
)

type EndpointService struct {
	EndpointRepo             *repo.EndpointRepo          `inject:""`
	ServeRepo                *repo.ServeRepo             `inject:""`
	EndpointInterfaceRepo    *repo.EndpointInterfaceRepo `inject:""`
	ServeServerRepo          *repo.ServeServerRepo       `inject:""`
	UserRepo                 *repo.UserRepo              `inject:""`
	CategoryRepo             *repo.CategoryRepo          `inject:""`
	DiagnoseInterfaceService *DiagnoseInterfaceService   `inject:""`
	EndpointTagRepo          *repo.EndpointTagRepo       `inject:""`
	EndpointTagService       *EndpointTagService         `inject:""`
	ServeService             *ServeService               `inject:""`
	MessageService           *MessageService             `inject:""`
	ThirdPartySyncService    *ThirdPartySyncService      `inject:""`
	DebugInterfaceRepo       *repo.DebugInterfaceRepo    `inject:""`
	EnvironmentRepo          *repo.EnvironmentRepo       `inject:""`
}

func (s *EndpointService) Paginate(req v1.EndpointReqPaginate) (ret _domain.PageData, err error) {
	ret, err = s.EndpointRepo.Paginate(req)
	return
}

func (s *EndpointService) Save(endpoint model.Endpoint) (res uint, err error) {

	if endpoint.ServerId == 0 {
		server, _ := s.ServeServerRepo.GetDefaultByServe(endpoint.ServeId)
		endpoint.ServerId = server.ID
	}

	if endpoint.Curl != "" {
		err = s.curlToEndpoint(&endpoint)
		if err != nil {
			return
		}
	}

	ret, _ := s.EndpointRepo.Get(endpoint.ID)

	err = s.EndpointRepo.SaveAll(&endpoint)

	//go func() {
	//	_ = s.SendEndpointMessage(endpoint.ProjectId, endpoint.ID, userId)
	//}()

	s.DebugInterfaceRepo.SyncPath(ret.ID, endpoint.ServeId, endpoint.Path, ret.Path)

	return endpoint.ID, err
}

func (s *EndpointService) SendEndpointMessage(projectId, endpointId, userId uint) (err error) {
	messageContent, err := s.MessageService.GetEndpointMcsData(projectId, endpointId)
	messageContentByte, _ := json.Marshal(messageContent)
	messageReq := v1.MessageReq{
		MessageBase: v1.MessageBase{
			MessageSource: consts.MessageSourceEndpoint,
			Content:       string(messageContentByte),
			ReceiverRange: 4,
			SenderId:      userId,
			ReceiverId:    projectId,
			SendStatus:    consts.MessageCreated,
			ServiceType:   consts.ServiceTypeInfo,
			BusinessId:    endpointId,
		},
	}
	_, _ = s.MessageService.Create(messageReq)

	return
}

func (s *EndpointService) GetById(id uint, version string) (res model.Endpoint) {
	res, _ = s.EndpointRepo.GetAll(id, version)
	s.SchemasConv(&res)
	return
}

func (s *EndpointService) DeleteById(id uint) (err error) {
	//var count int64
	//count, err = s.EndpointRepo.GetUsedCountByEndpointId(id)
	//if err != nil {
	//	return err
	//}
	//
	//if count > 0 {
	//	err = fmt.Errorf("this interface has already been used by scenarios, not allowed to delete")
	//	return err
	//}

	err = s.EndpointRepo.DeleteById(id)
	err = s.EndpointInterfaceRepo.DeleteByEndpoint(id)

	return
}

func (s *EndpointService) DeleteByCategories(categoryIds []uint) (err error) {
	endpointIds, err := s.EndpointRepo.ListEndpointByCategories(categoryIds)

	err = s.EndpointRepo.DeleteByIds(endpointIds)
	err = s.EndpointInterfaceRepo.DeleteByEndpoints(endpointIds)

	return
}

func (s *EndpointService) DisableById(id uint) (err error) {
	err = s.EndpointRepo.UpdateStatus(id, serverConsts.Abandoned)
	return
}

func (s *EndpointService) Publish(id uint) (err error) {
	err = s.EndpointRepo.UpdateStatus(id, serverConsts.Published)
	return
}

func (s *EndpointService) Develop(id uint) (err error) {
	err = s.EndpointRepo.UpdateStatus(id, serverConsts.Developing)
	return
}

func (s *EndpointService) Copy(id uint, version string) (res uint, err error) {

	endpoint, _ := s.EndpointRepo.GetAll(id, version)
	s.removeIds(&endpoint)
	endpoint.Title += "_copy"
	err = s.EndpointRepo.SaveAll(&endpoint)
	return endpoint.ID, err
}

func (s *EndpointService) removeIds(endpoint *model.Endpoint) {
	endpoint.ID = 0
	endpoint.CreatedAt = nil
	endpoint.UpdatedAt = nil
	for key, _ := range endpoint.PathParams {
		endpoint.PathParams[key].ID = 0
	}
	for key, _ := range endpoint.Interfaces {
		endpoint.Interfaces[key].ID = 0
		endpoint.Interfaces[key].RequestBody.ID = 0
		endpoint.Interfaces[key].RequestBody.SchemaItem.ID = 0
		endpoint.Interfaces[key].DebugInterfaceId = 0
		for key1, _ := range endpoint.Interfaces[key].Headers {
			endpoint.Interfaces[key].Headers[key1].ID = 0
		}
		for key1, _ := range endpoint.Interfaces[key].Cookies {
			endpoint.Interfaces[key].Cookies[key1].ID = 0
		}
		for key1, _ := range endpoint.Interfaces[key].Params {
			endpoint.Interfaces[key].Params[key1].ID = 0
		}
		for key1, _ := range endpoint.Interfaces[key].ResponseBodies {
			endpoint.Interfaces[key].ResponseBodies[key1].ID = 0
			endpoint.Interfaces[key].ResponseBodies[key1].SchemaItem.ID = 0
			for key2, _ := range endpoint.Interfaces[key].ResponseBodies[key1].Headers {
				endpoint.Interfaces[key].ResponseBodies[key1].Headers[key2].ID = 0
			}
		}
	}

}

func (s *EndpointService) Yaml(endpoint model.Endpoint) (res *openapi3.T) {
	var serve model.Serve
	if endpoint.ServeId != 0 {
		var err error
		serve, err = s.ServeRepo.Get(endpoint.ServeId)
		if err != nil {
			return
		}

		serveComponent, err := s.ServeRepo.GetSchemasByServeId(serve.ID)
		if err != nil {
			return
		}
		serve.Components = serveComponent

		serveServer, err := s.ServeRepo.ListServer(serve.ID)
		if err != nil {
			return
		}
		serve.Servers = serveServer

		securities, err := s.ServeRepo.ListSecurity(serve.ID)
		if err != nil {
			return
		}
		serve.Securities = securities
	}

	serve2conv := openapi.NewServe2conv(serve, []model.Endpoint{endpoint})
	res = serve2conv.ToV3()
	return
}

func (s *EndpointService) UpdateStatus(id uint, status int64) (err error) {
	err = s.EndpointRepo.UpdateStatus(id, status)
	return
}

func (s *EndpointService) BatchDelete(ids []uint) (err error) {
	err = s.EndpointRepo.DeleteByIds(ids)
	return
}

func (s *EndpointService) GetVersionsByEndpointId(endpointId uint) (res []model.EndpointVersion, err error) {
	res, err = s.EndpointRepo.GetVersionsByEndpointId(endpointId)
	return
}

func (s *EndpointService) GetLatestVersion(endpointId uint) (version string) {
	version = "v0.1.0"
	if res, err := s.EndpointRepo.GetLatestVersion(endpointId); err != nil {
		version = res.Version
	}
	return
}

func (s *EndpointService) AddVersion(version *model.EndpointVersion) (err error) {
	err = s.EndpointRepo.FindVersion(version)
	if err != nil {
		err = s.EndpointRepo.Save(0, version)
	} else {
		err = fmt.Errorf("version already exists")
	}
	return
}

func (s *EndpointService) SaveEndpoints(endpoints []*model.Endpoint, dirs *openapi.Dirs, components map[string]*model.ComponentSchema, req v1.ImportEndpointDataReq) (err error) {

	if dirs.Id == 0 || dirs.Id == -1 {
		root, _ := s.CategoryRepo.ListByProject(serverConsts.EndpointCategory, req.ProjectId)
		dirs.Id = int64(root[0].ID)
	}
	s.createDirs(dirs, req)
	wg := new(sync.WaitGroup)
	wg.Add(2)
	go s.createComponents(wg, components, req)
	go s.createEndpoints(wg, endpoints, dirs, req)
	wg.Wait()
	return
}

func (s *EndpointService) createEndpoints(wg *sync.WaitGroup, endpoints []*model.Endpoint, dirs *openapi.Dirs, req v1.ImportEndpointDataReq) (err error) {
	defer func() {
		wg.Done()
	}()

	userName := ""
	if req.UserId != 0 {
		user, _ := s.UserRepo.FindById(req.UserId)
		userName = user.Username
	}

	if req.CategoryId == -1 {
		if len(dirs.Dirs) > 0 {
			req.CategoryId = dirs.Id
		} else {
			dirs.Id = -1
		}
	}

	for _, endpoint := range endpoints {
		endpoint.ProjectId, endpoint.ServeId, endpoint.CategoryId = req.ProjectId, req.ServeId, req.CategoryId
		endpoint.Status = 1
		endpoint.SourceType = req.SourceType
		if endpoint.CreateUser == "" {
			endpoint.CreateUser = userName
		}
		endpoint.CategoryId = s.getCategoryId(endpoint.Tags, dirs)

		res, err := s.EndpointRepo.GetByItem(endpoint.SourceType, endpoint.ProjectId, endpoint.Path, endpoint.ServeId, req.CategoryId)

		//非Notfound
		if err != nil && err != gorm.ErrRecordNotFound {
			logUtils.Logger.Error(fmt.Sprintf("swagger import error:%s", err.Error()))
			continue
		}

		res, _ = s.EndpointRepo.GetAll(res.ID, "v0.1.0")

		//对比endpoint的时候不需要对比组件，所以服务ID设置为0
		endpoint.ServeId, res.ServeId = 0, 0
		openAPIDoc := s.Yaml(*endpoint)
		endpoint.Snapshot = _commUtils.JsonEncode(openAPIDoc)

		if req.DataSyncType == consts.FullCover {
			if err == nil {
				endpoint.ID = res.ID
				endpoint.CategoryId = res.CategoryId
				endpoint.ChangedStatus = consts.NoChanged
			}

		} else if req.DataSyncType == consts.AutoAdd {

			if err == nil {

				//远端无更新，则不做任何修改
				if endpoint.Snapshot == res.Snapshot {
					continue
				}

				//本地快照和本地数据不一致,更新快照,说明有修改，更新快照
				localEndpoint := s.Yaml(res)
				localEndpointJson := _commUtils.JsonEncode(localEndpoint)
				if res.Snapshot != localEndpointJson {
					s.EndpointRepo.UpdateSnapshot(res.ID, endpoint.Snapshot)
					continue
				} else { //一致覆盖数据
					endpoint.ID = res.ID
					endpoint.CategoryId = res.CategoryId
					now := time.Now()
					endpoint.ChangedTime = &now
				}

			}
		}
		endpoint.ServeId = req.ServeId //前面销毁了ID，现在补充上
		_, err = s.Save(*endpoint)
		if err != nil {
			//遇到错误跳过
			logUtils.Logger.Error(fmt.Sprintf("swagger import error:%s", err.Error()))
			//return err
		}
	}

	return
}

func (s *EndpointService) createComponents(wg *sync.WaitGroup, components map[string]*model.ComponentSchema, req v1.ImportEndpointDataReq) {
	defer func() {
		wg.Done()
	}()
	var newComponents []*model.ComponentSchema
	for _, component := range components {
		component.ServeId = int64(req.ServeId)
		component.SourceType = req.SourceType

		res, err := s.ServeRepo.GetComponentByItem(component.SourceType, uint(component.ServeId), component.Ref)
		if err != nil && err != gorm.ErrRecordNotFound {
			continue
		}

		if req.DataSyncType == consts.FullCover {
			if err == nil {
				component.ID = res.ID
			}
		} else if req.DataSyncType == consts.AutoAdd {
			if err == nil {
				continue
			}
		} else { //相同ref组件能创建新的
			continue
		}

		newComponents = append(newComponents, component)
	}

	s.ServeRepo.SaveSchemas(newComponents)

}

func (s *EndpointService) createDirs(data *openapi.Dirs, req v1.ImportEndpointDataReq) (err error) {
	for name, dirs := range data.Dirs {

		category := model.Category{Name: name, ParentId: int(data.Id), ProjectId: req.ProjectId, UseID: req.UserId, Type: serverConsts.EndpointCategory, SourceType: req.SourceType}
		//全覆盖更新目录
		res, err := s.CategoryRepo.GetByItem(uint(category.ParentId), category.Type, category.ProjectId, category.Name)
		if err != nil && err != gorm.ErrRecordNotFound {
			logUtils.Logger.Error(fmt.Sprintf("swagger import error:%s", err.Error()))
			continue
		}

		if err == nil { //同级目录下不创建同名目录
			category.ID = res.ID
			goto here
		}

		err = s.CategoryRepo.Save(&category)
		if err != nil {
			logUtils.Logger.Error(fmt.Sprintf("swagger import error:%s", err.Error()))
			return err
		}

	here:
		dirs.Id = int64(category.ID)
		err = s.createDirs(dirs, req)
		if err != nil {
			logUtils.Logger.Error(fmt.Sprintf("swagger import error:%s", err.Error()))
			return err
		}
	}
	return
}

func (s *EndpointService) getCategoryId(tags []string, dirs *openapi.Dirs) int64 {
	rootId := dirs.Id
	for _, tag := range tags {
		dirs = dirs.Dirs[tag]
	}
	if dirs.Id == rootId && rootId == 0 {
		return -1
	}
	return dirs.Id
}

func (s *EndpointService) BatchUpdateByField(req v1.BatchUpdateReq) (err error) {
	if _commUtils.InSlice(req.FieldName, []string{"status", "categoryId", "serveId", "description"}) {
		err = s.EndpointRepo.BatchUpdate(req.EndpointIds, map[string]interface{}{_commUtils.Camel2Case(req.FieldName): req.Value})
		if req.FieldName == "serveId" { //修改debug表serveId
			if serveId, ok := req.Value.(float64); ok {
				s.DebugInterfaceRepo.SyncServeId(req.EndpointIds, uint(serveId))
			}
		}

	} else {
		err = errors.New("字段错误")
	}
	return
}

func (s *EndpointService) curlToEndpoint(endpoint *model.Endpoint) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("curl格式错误")
		}
	}()
	curlObj := curlHelper.Parse(endpoint.Curl)
	wf := curlObj.CreateTemporary(curlObj.CreateSession())

	endpoint.Path = curlObj.ParsedURL.Path

	endpoint.Interfaces = s.getInterfaces(endpoint.Title, curlObj, wf)

	return
}

func (s *EndpointService) getInterfaces(name string, cURL *curlHelper.CURL, wf *requests.Temporary) (interfaces []model.EndpointInterface) {
	interf := model.EndpointInterface{}
	interf.Name = name
	interf.Params = s.getQueryParams(wf.GetQuery())
	interf.Headers = s.getHeaders(wf.Header)
	interf.Cookies = s.getCookies(wf.Cookies)
	bodyType := ""
	contentType := strings.Split(cURL.ContentType, ";")
	if len(contentType) >= 1 {
		bodyType = contentType[0]
	}
	interf.BodyType = consts.HttpContentType(bodyType)
	interf.RequestBody = s.getRequestBody(wf.Body.String())
	interf.RequestBody.MediaType = string(interf.BodyType)
	interf.Method = s.getMethod(bodyType, cURL.Method)
	interfaces = append(interfaces, interf)

	return
}

func (s *EndpointService) getMethod(contentType, method string) (ret consts.HttpMethod) {
	if method == "" && contentType == "application/json" {
		method = "POST"
	}

	return consts.HttpMethod(method)

}

func (s *EndpointService) getQueryParams(params url.Values) (ret []model.EndpointInterfaceParam) {
	m := map[string]bool{}
	for key, arr := range params {
		for _, item := range arr {
			if _, ok := m[key]; ok {
				continue
			}
			ret = append(ret, model.EndpointInterfaceParam{
				SchemaParam: model.SchemaParam{Name: key, Type: "string", Value: item, Default: item, Example: item},
			})
			m[key] = true
		}
	}

	return
}

func (s *EndpointService) getHeaders(header http.Header) (ret []model.EndpointInterfaceHeader) {
	for key, arr := range header {
		for _, item := range arr {
			ret = append(ret, model.EndpointInterfaceHeader{
				SchemaParam: model.SchemaParam{Name: key, Type: "string", Value: item, Default: item, Example: item},
			})
		}
	}

	return
}

func (s *EndpointService) getCookies(cookies map[string]*http.Cookie) (ret []model.EndpointInterfaceCookie) {
	for _, item := range cookies {
		ret = append(ret, model.EndpointInterfaceCookie{
			SchemaParam: model.SchemaParam{Name: item.Name, Type: "string", Value: item.Value, Default: item.Value, Example: item.Value},
		})
	}

	return
}

func (s *EndpointService) getRequestBody(body string) (requestBody model.EndpointInterfaceRequestBody) {
	requestBody = model.EndpointInterfaceRequestBody{}

	if body != "" {
		var examples []map[string]string
		examples = append(examples, map[string]string{"content": body, "name": "defaultExample"})
		requestBody.Examples = _commUtils.JsonEncode(examples)
	}

	requestBody.SchemaItem = s.getRequestBodyItem(body)
	return
}

func (s *EndpointService) getRequestBodyItem(body string) (requestBodyItem model.EndpointInterfaceRequestBodyItem) {
	requestBodyItem = model.EndpointInterfaceRequestBodyItem{}
	requestBodyItem.Type = "object"
	schema2conv := schemaHelper.NewSchema2conv()
	var obj interface{}
	schema := schemaHelper.Schema{}
	_commUtils.JsonDecode(body, &obj)
	schema2conv.Example2Schema(obj, &schema)
	requestBodyItem.Content = _commUtils.JsonEncode(schema)
	return
}

func (s *EndpointService) UpdateTags(req v1.EndpointTagReq, projectId uint) (err error) {
	if err = s.EndpointTagRepo.DeleteRelByEndpointAndProject(req.Id, projectId); err != nil {
		return
	}

	if len(req.TagNames) > 0 {
		err = s.EndpointTagRepo.BatchAddRel(req.Id, projectId, req.TagNames)
	}
	return
	//oldTagIds, err := s.EndpointTagRepo.GetTagIdsByEndpointId(req.Id)
	//if err != nil && err != gorm.ErrRecordNotFound {
	//	return
	//}
	//
	//intTagIds, err := s.EndpointTagService.GetTagIdsNyName(req.TagNames, projectId)
	//if err != nil && err != gorm.ErrRecordNotFound {
	//	return
	//}
	//
	//tagsNeedDeleted := _commUtils.DifferenceUint(oldTagIds, intTagIds)
	//
	//if err = s.EndpointTagRepo.DeleteRelByEndpointId(req.Id); err != nil {
	//	return
	//}
	//
	//if len(intTagIds) > 0 {
	//	err = s.EndpointTagRepo.AddRel(req.Id, intTagIds)
	//	if err != nil {
	//		return err
	//	}
	//}
	//
	//for _, v := range tagsNeedDeleted {
	//	relations, err := s.EndpointTagRepo.ListRelByTagId(v)
	//	if err != nil && err != gorm.ErrRecordNotFound {
	//		return err
	//	}
	//
	//	if len(relations) == 0 {
	//		if err = s.EndpointTagRepo.DeleteById(v); err != nil {
	//			return err
	//		}
	//	}
	//}
	//return
}

func (s *EndpointService) SchemasConv(endpoint *model.Endpoint) {
	schema2conv := schemaHelper.NewSchema2conv()
	schema2conv.Components = s.ServeService.Components(endpoint.ServeId)
	for key, intef := range endpoint.Interfaces {
		for k, response := range intef.ResponseBodies {
			schema := new(schemaHelper.SchemaRef)
			_commUtils.JsonDecode(response.SchemaItem.Content, schema)
			if endpoint.SourceType == 1 && schema.Value != nil && len(schema.Value.AllOf) > 0 {
				schema2conv.CombineSchemas(schema)
			}
			endpoint.Interfaces[key].ResponseBodies[k].SchemaItem.Content = _commUtils.JsonEncode(schema)
		}
	}

}

func (s *EndpointService) SchemaConv(interf *model.EndpointInterface, serveId uint) {
	schema2conv := schemaHelper.NewSchema2conv()
	schema2conv.Components = s.ServeService.Components(serveId)
	for k, response := range interf.ResponseBodies {
		schema := new(schemaHelper.SchemaRef)
		_commUtils.JsonDecode(response.SchemaItem.Content, schema)
		if schema.Value != nil && len(schema.Value.AllOf) > 0 {
			schema2conv.CombineSchemas(schema)
		}
		interf.ResponseBodies[k].SchemaItem.Content = _commUtils.JsonEncode(schema)
	}
}

func (s *EndpointService) UpdateAdvancedMockDisabled(endpointId uint, disabled bool) (err error) {
	err = s.EndpointRepo.UpdateAdvancedMockDisabled(endpointId, disabled)
	return
}

func (s *EndpointService) CreateExample(req v1.CreateExampleReq) (ret interface{}, err error) {
	var endpoint model.Endpoint
	endpoint, err = s.EndpointRepo.Get(req.EndpointId)
	if err != nil {
		return
	}

	var bodyItem model.EndpointInterfaceResponseBodyItem
	bodyItem, err = s.EndpointInterfaceRepo.GetResponseDefine(req.EndpointId, req.Method, req.Code)
	if err != nil || bodyItem.Content == "" {
		return
	}

	schema2conv := schemaHelper.NewSchema2conv()
	schema2conv.Components = s.ServeService.Components(endpoint.ServeId)

	schema := schemaHelper.SchemaRef{}
	_commUtils.JsonDecode(bodyItem.Content, &schema)

	ret = schema2conv.Schema2Example(schema)

	return

}

func (s *EndpointService) SyncFromThirdParty(endpointId uint) (err error) {
	endpoint, err := s.EndpointRepo.Get(endpointId)
	endpoint.Interfaces, _ = s.EndpointInterfaceRepo.ListByEndpointId(endpoint.ID, "v0.1.0")
	if err != nil {
		return
	}

	if endpoint.SourceType != consts.ThirdPartySync || endpoint.CategoryId == -1 || len(endpoint.Interfaces) == 0 {
		return
	}

	pathArr := strings.Split(endpoint.Path, "/")

	err = s.ThirdPartySyncService.SyncFunctionBody(endpoint.ProjectId, endpoint.ServeId, endpoint.Interfaces[0].ID, pathArr[2], pathArr[3])
	if err != nil {
		return
	}

	err = s.EndpointRepo.UpdateBodyIsChanged(endpointId, consts.Changed)

	return
}

func (s *EndpointService) GetDiff(endpointId uint) (res v1.EndpointDiffRes, err error) {
	var endpoint model.Endpoint
	var resYaml []byte
	endpoint, err = s.EndpointRepo.GetAll(endpointId, "v0.1.0")
	if err != nil {
		return
	}

	var sourceName string
	if endpoint.SourceType == consts.SwaggerSync {
		sourceName = "Swagger"
	} else if endpoint.SourceType == consts.SwaggerImport {
		sourceName = "接口定义"
	} else if endpoint.SourceType == consts.ThirdPartySync {
		sourceName = "乐仓智能体厂"
	}

	res.ChangedStatus = endpoint.ChangedStatus

	res.CurrentDesc = fmt.Sprintf("%s于%s在系统中手动更新", endpoint.CreateUser, endpoint.UpdatedAt.Format("2006-01-02 15:04:05"))
	res.LatestDesc = fmt.Sprintf("%s从%s自动同步", endpoint.ChangedTime.Format("2006-01-02 15:04:05"), sourceName)

	var ret interface{}
	endpoint.ServeId = 0
	_commUtils.JsonDecode(_commUtils.JsonEncode(s.Yaml(endpoint)), &ret)
	resYaml, err = encoder.NewEncoder(ret).Encode()
	if err != nil {
		return
	}
	res.Current = string(resYaml)

	_commUtils.JsonDecode(endpoint.Snapshot, &ret)
	resYaml, err = encoder.NewEncoder(ret).Encode()
	if err != nil {
		return
	}
	res.Latest = string(resYaml)
	return
}

func (s *EndpointService) SaveDiff(endpointId uint, isChanged bool, userName string) (err error) {
	endpoint, err := s.EndpointRepo.GetAll(endpointId, "v0.1.0")
	if err != nil {
		return
	}

	if isChanged {
		var doc openapi3.T
		_commUtils.JsonDecode(endpoint.Snapshot, &doc)
		endpoints, _, _ := openapi.NewOpenapi2endpoint(&doc, endpoint.CategoryId).Convert()
		endpoints[0].ID = endpoint.ID
		endpoints[0].Title = endpoint.Title
		endpoints[0].ServeId = endpoint.ServeId
		endpoints[0].ChangedStatus = consts.NoChanged
		endpoints[0].ProjectId = endpoint.ProjectId
		endpoints[0].GlobalParams = endpoint.GlobalParams
		endpoints[0].UpdateUser = userName
		err = s.EndpointRepo.SaveAll(endpoints[0])
	} else {
		err = s.EndpointRepo.UpdateBodyIsChanged(endpointId, consts.IgnoreChanged)
	}

	return
}

func (s *EndpointService) isEqualEndpoint(old, new model.Endpoint) bool {
	var ret interface{}
	oldYaml := s.Yaml(old)
	_commUtils.JsonDecode(_commUtils.JsonEncode(oldYaml), &ret)
	oldYamlByte, _ := encoder.NewEncoder(ret).Encode()

	newYaml := s.Yaml(new)
	_commUtils.JsonDecode(_commUtils.JsonEncode(newYaml), &ret)
	newYamlByte, _ := encoder.NewEncoder(ret).Encode()
	res1, res2 := string(oldYamlByte), string(newYamlByte)

	return res1 == res2

}

func (s *EndpointService) UpdateName(id uint, name string) (err error) {
	err = s.EndpointRepo.UpdateName(id, name)
	return
}
