package repo

import (
	serverDomain "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"gorm.io/gorm"
	"strconv"
)

type EndpointCaseAlternativeRepo struct {
	*BaseRepo `inject:""`
	DB        *gorm.DB `inject:""`

	EndpointRepo       *EndpointRepo       `inject:""`
	DebugInterfaceRepo *DebugInterfaceRepo `inject:""`
	ProjectRepo        *ProjectRepo        `inject:""`
	CategoryRepo       *CategoryRepo       `inject:""`
}

func (r *EndpointCaseAlternativeRepo) List(caseId uint) (pos []model.EndpointCaseAlternative, err error) {
	err = r.DB.
		Where("base_id=?", caseId).
		Where("NOT deleted").Order("created_at desc").
		Find(&pos).Error

	return
}

func (r *EndpointCaseAlternativeRepo) LoadFactor(caseId uint) (pos []model.EndpointCaseAlternativeFactor, err error) {
	err = r.DB.
		Where("case_id=?", caseId).
		Where("NOT deleted").
		Find(&pos).Error

	return
}

func (r *EndpointCaseAlternativeRepo) Get(id uint) (po model.EndpointCase, err error) {
	err = r.DB.Where("id = ?", id).First(&po).Error
	return
}

func (r *EndpointCaseAlternativeRepo) GetDetail(caseId uint) (endpointCase model.EndpointCase, err error) {
	if caseId <= 0 {
		return
	}

	endpointCase, err = r.Get(caseId)

	debugInterface, _ := r.DebugInterfaceRepo.Get(endpointCase.DebugInterfaceId)

	debugData, _ := r.DebugInterfaceRepo.GetDetail(debugInterface.ID)
	endpointCase.DebugData = &debugData

	return
}

func (r *EndpointCaseAlternativeRepo) Save(po *model.EndpointCase) (err error) {
	err = r.DB.Save(po).Error

	err = r.UpdateSerialNumber(po.ID, po.ProjectId)

	return
}

func (r *EndpointCaseAlternativeRepo) UpdateName(req serverDomain.EndpointCaseSaveReq) (err error) {
	err = r.DB.Model(&model.EndpointCase{}).
		Where("id=?", req.ID).
		Update("name", req.Name).Error

	return
}

func (r *EndpointCaseAlternativeRepo) Remove(id uint) (err error) {
	err = r.DB.Model(&model.EndpointCase{}).
		Where("id = ?", id).
		Update("deleted", true).Error

	return
}

func (r *EndpointCaseAlternativeRepo) SaveDebugData(interf *model.EndpointCase) (err error) {
	r.DB.Transaction(func(tx *gorm.DB) error {
		err = r.UpdateDebugInfo(interf)
		if err != nil {
			return err
		}

		// TODO: save debug data

		return err
	})

	return
}

func (r *EndpointCaseAlternativeRepo) UpdateDebugInfo(interf *model.EndpointCase) (err error) {
	values := map[string]interface{}{
		"server_id": interf.DebugData.ServerId,
		"base_url":  interf.DebugData.BaseUrl,
		"url":       interf.DebugData.Url,
		"method":    interf.DebugData.Method,
	}

	err = r.DB.Model(&model.EndpointCase{}).
		Where("id=?", interf.ID).
		Updates(values).
		Error

	return
}

func (r *EndpointCaseAlternativeRepo) UpdateSerialNumber(id, projectId uint) (err error) {
	var project model.Project
	project, err = r.ProjectRepo.Get(projectId)
	if err != nil {
		return
	}

	err = r.DB.Model(&model.EndpointCase{}).
		Where("id=?", id).
		Update("serial_number", project.ShortName+"-TC-"+strconv.Itoa(int(id))).Error
	return
}

func (r *EndpointCaseAlternativeRepo) ListByProjectIdAndServeId(projectId, serveId uint) (endpointCases []*serverDomain.InterfaceCase, err error) {
	err = r.DB.Model(&model.EndpointCase{}).
		Joins("left join biz_debug_interface d on biz_endpoint_case.debug_interface_id=d.id").
		Select("biz_endpoint_case.*, d.method as method").
		Where("biz_endpoint_case.project_id = ? and biz_endpoint_case.serve_id = ? and processor_interface_src = '' and not biz_endpoint_case.deleted and not biz_endpoint_case.disabled", projectId, serveId).
		Find(&endpointCases).Error
	//err = r.DB.Where("project_id = ? and serve_id = ? and not deleted and not disabled", projectId, serveId).Find(&endpointCases).Error
	return
}

func (r *EndpointCaseAlternativeRepo) GetEndpointCount(projectId, serveId uint) (result []serverDomain.EndpointCount, err error) {
	err = r.DB.Raw("select count(id) count,endpoint_id from "+model.EndpointCase{}.TableName()+" where not deleted and not disabled and project_id=? and serve_id =? group by endpoint_id", projectId, serveId).Scan(&result).Error
	return
}

func (r *EndpointCaseAlternativeRepo) GetCategoryEndpointCase(projectId, serveId uint) (result []serverDomain.CategoryEndpointCase, err error) {
	err = r.DB.Raw("select concat('case_',ec.id) as case_unique_id,concat('endpoint_',e.id) as endpoint_unique_id,ec.id as case_id,ec.name as case_name,i.method,ec.`desc` as case_desc,ec.endpoint_id as case_endpoint_id,ec.debug_interface_id as case_debug_interface_id,ec.project_id,ec.serve_id,e.id as endpoint_id,e.title as endpoint_title,e.description as endpoint_description,e.category_id as category_id from biz_endpoint_case ec left join biz_endpoint e on ec.endpoint_id=e.id left join biz_debug_interface i on ec.debug_interface_id=i.id Where ec.project_id= ? and ec.serve_id=? and not e.deleted and not ec.deleted", projectId, serveId).Scan(&result).Error
	return
}

func (r *EndpointCaseAlternativeRepo) UpdateDebugInterfaceId(debugInterfaceId, id uint) (err error) {
	err = r.DB.Model(&model.EndpointCase{}).
		Where("id=?", id).
		Update("debug_interface_id", debugInterfaceId).Error

	return
}

func (r *EndpointCaseAlternativeRepo) SaveFactor(req serverDomain.EndpointCaseFactorSaveReq) (err error) {
	var po model.EndpointCaseAlternativeFactor
	err = r.DB.Where("case_id = ? AND path = ?", req.CaseId, req.Path).First(&po).Error

	if po.ID > 0 {
		err = r.DB.Model(&model.EndpointCaseAlternativeFactor{}).
			Where("case_id = ? AND path = ?", req.CaseId, req.Path).
			Update("value", req.Value).Error

	} else {
		po = model.EndpointCaseAlternativeFactor{
			CaseId: uint(req.CaseId),
			Path:   req.Path,
			Value:  req.Value,
		}

		err = r.DB.Save(&po).Error
	}

	return
}
