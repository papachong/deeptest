package repo

import (
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/pkg/domain"
	responseDefineHelper "github.com/aaronchen2k/deeptest/internal/pkg/helper/schema"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	_commUtils "github.com/aaronchen2k/deeptest/pkg/lib/comm"
	"github.com/jinzhu/copier"
	"gorm.io/gorm"
	"strings"
)

type ResponseDefineRepo struct {
	DB                    *gorm.DB               `inject:""`
	ServeRepo             *ServeRepo             `inject:""`
	EndpointRepo          *EndpointRepo          `inject:""`
	EndpointInterfaceRepo *EndpointInterfaceRepo `inject:""`
}

func (r *ResponseDefineRepo) Get(id uint) (responseDefine model.DebugConditionResponseDefine, err error) {
	err = r.DB.
		Where("id=?", id).
		Where("NOT deleted").
		First(&responseDefine).Error
	return
}

func (r *ResponseDefineRepo) Save(responseDefine *model.DebugConditionResponseDefine) (err error) {

	err = r.DB.Save(responseDefine).Error
	return
}

func (r *ResponseDefineRepo) UpdateResult(responseDefine domain.ResponseDefineBase) (err error) {
	values := map[string]interface{}{
		"result_msg":    responseDefine.ResultMsg,
		"result_status": responseDefine.ResultStatus,
	}

	err = r.DB.Model(&model.DebugConditionResponseDefine{}).
		Where("id=?", responseDefine.ConditionEntityId).
		Updates(values).
		Error

	return
}

func (r *ResponseDefineRepo) Update(id uint, data map[string]interface{}) (err error) {
	err = r.DB.Model(&model.DebugConditionResponseDefine{}).
		Where("id=?", id).
		Updates(data).
		Error
	return
}

func (r *ResponseDefineRepo) CreateLog(responseDefine domain.ResponseDefineBase) (
	log model.ExecLogResponseDefine, err error) {

	copier.CopyWithOption(&log, responseDefine, copier.Option{DeepCopy: true})

	log.ID = 0
	log.ConditionId = responseDefine.ConditionId
	log.ConditionEntityId = responseDefine.ConditionEntityId

	log.InvokeId = responseDefine.InvokeId
	log.CreatedAt = nil
	log.UpdatedAt = nil

	err = r.DB.Save(&log).Error

	return
}

func (r *ResponseDefineRepo) GetLog(conditionId, invokeId uint) (ret model.ExecLogResponseDefine, err error) {
	err = r.DB.
		Where("condition_id=? AND invoke_id=?", conditionId, invokeId).
		Where("NOT deleted").
		First(&ret).Error

	ret.ConditionEntityType = consts.ConditionTypeResponseDefine

	return
}

func (r *ResponseDefineRepo) Components(endpointInterfaceId uint) responseDefineHelper.Components {
	endpointInterface, _ := r.EndpointInterfaceRepo.Get(endpointInterfaceId)
	endpoint, _ := r.EndpointRepo.Get(endpointInterface.EndpointId)

	components := responseDefineHelper.Components{}
	result, err := r.ServeRepo.GetSchemasByServeId(endpoint.ServeId)
	if err != nil {
		return nil
	}
	responseBodies, err := r.EndpointInterfaceRepo.ListResponseBodies(endpointInterfaceId)
	if err != nil {
		return nil
	}

	for _, item := range result {
		var schema responseDefineHelper.SchemaRef
		_commUtils.JsonDecode(item.Content, &schema)
		components[item.Ref] = &schema
	}

	return r.requiredComponents(responseBodies, components)

}

func (r *ResponseDefineRepo) requiredComponents(responseBodies []model.EndpointInterfaceResponseBody, components responseDefineHelper.Components) (ret responseDefineHelper.Components) {
	ret = responseDefineHelper.Components{}
	for _, responseBody := range responseBodies {
		r.dependComponents(responseBody, components, ret)
	}

	return
}

func (r *ResponseDefineRepo) dependComponents(responseBody model.EndpointInterfaceResponseBody, components, dependComponents responseDefineHelper.Components) {
	schema := new(responseDefineHelper.SchemaRef)
	responseBody.SchemaItem.Content = strings.ReplaceAll(responseBody.SchemaItem.Content, "\\u0026", "&")
	responseBody.SchemaItem.Content = strings.ReplaceAll(responseBody.SchemaItem.Content, "\n", "")
	responseBody.SchemaItem.Content = strings.ReplaceAll(responseBody.SchemaItem.Content, "\"ref\":", "\"$ref\":")
	_commUtils.JsonDecode(responseBody.SchemaItem.Content, schema)
	schema2conv := responseDefineHelper.NewSchema2conv()
	schema2conv.Components = components
	schema2conv.SchemaComponents(*schema, dependComponents)
}
