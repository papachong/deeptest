package repo

import (
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"gorm.io/gorm"
)

type DebugRepo struct {
	DB *gorm.DB `inject:""`
}

func (r *DebugRepo) List(debugInterfaceId, endpointInterfaceId uint) (pos []model.DebugInvoke, err error) {
	db := r.DB.Select("id", "name")

	if debugInterfaceId > 0 { // debugInterfaceId first
		db.Where("debug_interface_id=?", debugInterfaceId)

	} else if endpointInterfaceId > 0 {
		db.Where("endpoint_interface_id=? AND debug_interface_id=?", endpointInterfaceId, 0)

	}

	err = db.Where("NOT deleted").
		Order("created_at DESC").
		Find(&pos).Error
	return
}

func (r *DebugRepo) GetLast(debugInterfaceId, endpointInterfaceId uint) (debug model.DebugInvoke, err error) {
	db := r.DB

	if debugInterfaceId > 0 { // debugInterfaceId first
		db = db.Where("debug_interface_id=?", debugInterfaceId)
	} else if endpointInterfaceId > 0 {
		db = db.Where("endpoint_interface_id=? AND debug_interface_id=?", endpointInterfaceId, 0)
	}

	err = db.Where("NOT deleted").
		Order("created_at DESC").
		First(&debug).Error

	return
}

func (r *DebugRepo) Get(id uint) (invocation model.DebugInvoke, err error) {
	err = r.DB.
		Where("id=?", id).
		Where("NOT deleted").
		First(&invocation).Error
	return
}

func (r *DebugRepo) Delete(id uint) (err error) {
	err = r.DB.Model(&model.DebugInvoke{}).
		Where("id=?", id).
		Update("deleted", true).
		Error

	return
}
