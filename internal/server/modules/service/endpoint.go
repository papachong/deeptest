package service

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/openapi"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	_domain "github.com/aaronchen2k/deeptest/pkg/domain"
)

type EndpointService struct {
	EndpointRepo          *repo.EndpointRepo          `inject:""`
	ServeRepo             *repo.ServeRepo             `inject:""`
	EndpointInterfaceRepo *repo.EndpointInterfaceRepo `inject:""`
	ServeServerRepo       *repo.ServeServerRepo       `inject:""`
	UserRepo              *repo.UserRepo              `inject:""`
	CategoryRepo          *repo.CategoryRepo          `inject:""`
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

	err = s.EndpointRepo.SaveAll(&endpoint)
	return endpoint.ID, err
}

func (s *EndpointService) GetById(id uint, version string) (res model.Endpoint) {
	res, _ = s.EndpointRepo.GetAll(id, version)
	return
}

func (s *EndpointService) DeleteById(id uint) (err error) {
	var count int64
	count, err = s.EndpointRepo.GetUsedCountByEndpointId(id)
	if err != nil {
		return err
	}

	if count > 0 {
		err = fmt.Errorf("this interface has already been used by scenarios, not allowed to delete")
		return err
	}

	err = s.EndpointRepo.DeleteById(id)
	err = s.EndpointInterfaceRepo.DeleteByEndpoint(id)

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

func (s *EndpointService) Yaml(endpoint model.Endpoint) (res interface{}) {
	serve, err := s.ServeRepo.Get(endpoint.ServeId)
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

	Securities, err := s.ServeRepo.ListSecurity(serve.ID)
	if err != nil {
		return
	}
	serve.Securities = Securities

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

func (s *EndpointService) SaveEndpoints(endpoints []*model.Endpoint, dirs *openapi.Dirs, components []*model.ComponentSchema, req v1.ImportEndpointDataReq) (err error) {
	user, _ := s.UserRepo.FindById(req.UserId)

	if dirs.Id == 0 || dirs.Id == -1 {
		root, _ := s.CategoryRepo.ListByProject(serverConsts.EndpointCategory, req.ProjectId, 0)
		dirs.Id = int64(root[0].ID)
	}

	s.createDirs(dirs, req)
	s.createComponents(components, req)
	for _, endpoint := range endpoints {
		endpoint.ProjectId, endpoint.ServeId, endpoint.CategoryId, endpoint.CreateUser = req.ProjectId, req.ServeId, req.CategoryId, user.Name
		endpoint.Status = 1
		endpoint.CategoryId = s.getCategoryId(endpoint.Tags, dirs)

		_, err = s.Save(*endpoint)
		if err != nil {
			return
		}

	}
	//err = s.EndpointRepo.CreateEndpoints(endpoints)
	return
}

func (s *EndpointService) createComponents(components []*model.ComponentSchema, req v1.ImportEndpointDataReq) {
	for _, component := range components {
		component.ServeId = int64(req.ServeId)
	}
	s.ServeRepo.CreateSchemas(components)
}

func (s *EndpointService) createDirs(data *openapi.Dirs, req v1.ImportEndpointDataReq) (err error) {

	for name, dirs := range data.Dirs {
		category := model.Category{Name: name, ParentId: int(data.Id), ProjectId: req.ProjectId, UseID: req.UserId, Type: serverConsts.EndpointCategory}
		err = s.CategoryRepo.Save(&category)
		if err != nil {
			return
		}

		dirs.Id = int64(category.ID)
		err = s.createDirs(dirs, req)
		if err != nil {
			return err
		}
	}
	return
}

func (s *EndpointService) getCategoryId(tags []string, dirs *openapi.Dirs) int64 {
	for _, tag := range tags {
		dirs = dirs.Dirs[tag]
	}

	return dirs.Id
}
