package service

import (
	"context"
	"github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/helper/openapi"
	serverConsts "github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	repo "github.com/aaronchen2k/deeptest/internal/server/modules/repo"
	"github.com/getkin/kin-openapi/openapi3"
)

type ImportService struct {
	InterfaceRepo    *repo.InterfaceRepo   `inject:""`
	EnvironmentRepo  *repo.EnvironmentRepo `inject:""`
	InterfaceService *InterfaceService     `inject:""`
}

func (s *ImportService) Import(req domain.InterfaceImportReq, targetId int) (err error) {
	interf, _ := s.InterfaceRepo.Get(uint(targetId))

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}

	doc3, err := loader.LoadFromFile(req.File)
	if err != nil {
		return
	}

	err = s.GenerateInterface(doc3, uint(targetId), interf.ProjectId)
	if err != nil {
		return
	}

	err = s.GenerateEnvironment(doc3, interf.ProjectId)
	if err != nil {
		return
	}

	return
}

func (s *ImportService) GenerateInterface(doc *openapi3.T, targetId, projectId uint) (err error) {
	interfaces, err := openapi.ConvertPathsToInterfaces(doc)
	if err != nil {
		return
	}

	for _, interf := range interfaces {
		interf.ProjectId = projectId
		interf.ParentId = targetId

		s.Create(&interf)
	}

	return
}

func (s *ImportService) GenerateEnvironment(doc *openapi3.T, projectId uint) (err error) {
	env, err := s.EnvironmentRepo.GetByProject(projectId)
	if err != nil {
		return
	}

	envVars, err := openapi.ConvertServersToEnvironments(doc.Servers)
	if err != nil {
		return
	}

	for _, vari := range envVars {
		po, _ := s.EnvironmentRepo.GetSameVar(vari, env.ID)

		if po.ID == 0 {
			vari.EnvironmentId = env.ID
			s.EnvironmentRepo.SaveVar(&vari)
		}
	}

	return
}

func (s *ImportService) Create(interf *model.Interface) (err error) {
	interf.ParentId, interf.Ordr = s.InterfaceRepo.UpdateOrder(serverConsts.Inner, interf.ParentId)
	err = s.InterfaceRepo.Save(interf)
	if err != nil {
		return err
	}

	err = s.InterfaceRepo.UpdateParams(interf.ID, interf.Params)
	if err != nil {
		return err
	}

	err = s.InterfaceRepo.UpdateHeaders(interf.ID, interf.Headers)
	if err != nil {
		return err
	}

	err = s.InterfaceRepo.UpdateBasicAuth(interf.ID, interf.BasicAuth)
	if err != nil {
		return err
	}

	err = s.InterfaceRepo.UpdateBearerToken(interf.ID, interf.BearerToken)
	if err != nil {
		return err
	}

	err = s.InterfaceRepo.UpdateOAuth20(interf.ID, interf.OAuth20)
	if err != nil {
		return err
	}

	err = s.InterfaceRepo.UpdateApiKey(interf.ID, interf.ApiKey)
	if err != nil {
		return err
	}

	err = s.InterfaceRepo.UpdateBodyFormData(interf.ID, interf.BodyFormData)
	if err != nil {
		return err
	}

	return
}

//func (s *ImportService) OpenApi2To3(src []byte) (ret []byte, err error) {
//	var doc2 openapi2.T
//	err = json.Unmarshal(src, &doc2)
//
//	doc3, err := openapi.ToV3(&doc2)
//	err = doc3.Validate(context.Background())
//
//	ret, err = json.Marshal(doc3)
//
//	return
//}
//
//func (s *ImportService) PostmanToOpenApi3(pth string) (ret []byte, err error) {
//	cmd := fmt.Sprintf(`p2o %s`, pth)
//
//	out, _ := _shellUtils.ExeShell(cmd)
//
//	ret = []byte(out)
//
//	return
//}
