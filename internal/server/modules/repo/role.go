package repo

import (
	"errors"
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/server/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/casbin"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	"github.com/aaronchen2k/deeptest/internal/server/modules/model"
	myZap "github.com/aaronchen2k/deeptest/pkg/core/zap"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"strconv"
)

type RoleRepo struct {
	DB *gorm.DB `inject:""`
}

// Paginate
func (r *RoleRepo) Paginate(req v1.RoleReqPaginate) (data _domain.PageData, err error) {
	var count int64

	db := r.DB.Model(&model.SysRole{})
	if req.Name != "" {
		db = db.Where("name LIKE ?", fmt.Sprintf("%s%%", req.Name))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("获取角色总数错误", zap.String("错误:", err.Error()))
		return
	}

	var roles []*v1.RoleResp
	err = db.Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).Find(&roles).Error
	if err != nil {
		logUtils.Errorf("获取角色分页数据错误", zap.String("错误:", err.Error()))
		return
	}

	data.Result = roles
	data.Populate(roles, count, req.Page, req.PageSize)

	return
}

// FindByName
func (r *RoleRepo) FindByName(name string, ids ...uint) (v1.RoleResp, error) {
	role := v1.RoleResp{}
	db := r.DB.Model(&model.SysRole{}).Where("name = ?", name)
	if len(ids) == 1 {
		db.Where("id != ?", ids[0])
	}
	err := db.First(&role).Error
	if err != nil {
		return role, err
	}
	return role, nil
}

// FindByName
func (r *RoleRepo) FindFirstAdminUser() (v1.RoleResp, error) {
	role := v1.RoleResp{}
	err := r.DB.Model(&model.SysRole{}).Where("true").First(&role).Error

	if err != nil {
		logUtils.Errorf("管理员角色不存在，错误%s。", err.Error())
		return role, err
	}
	return role, nil
}

func (r *RoleRepo) Create(req v1.RoleReq) (roleId uint, err error) {
	role := model.SysRole{RoleBase: req.RoleBase}
	roleRes, err := r.FindByName(req.Name)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		err = r.DB.Create(&role).Error
		if err != nil {
			logUtils.Errorf("create data err ", zap.String("错误:", err.Error()))
			return
		}
		roleId = role.ID
	} else {
		roleId = roleRes.Id
	}

	err = r.AddPermForRole(roleId, req.Perms)
	if err != nil {
		logUtils.Errorf("添加权限到角色错误", zap.String("错误:", err.Error()))
		return
	}

	return
}

func (r *RoleRepo) Update(id uint, req v1.RoleReq) error {
	if b, err := r.IsAdminRole(id); err != nil {
		return err
	} else if b {
		return errors.New("不能编辑超级管理员")
	}
	_, err := r.FindByName(req.Name, id)
	if !errors.Is(err, gorm.ErrRecordNotFound) {
		logUtils.Errorf("角色名称已经被使用")
		return err
	}
	role := model.SysRole{RoleBase: req.RoleBase}
	err = r.DB.Model(&model.SysRole{}).Where("id = ?", id).Updates(&role).Error
	if err != nil {
		logUtils.Errorf("更新角色错误", zap.String("错误:", err.Error()))
		return err
	}
	err = r.AddPermForRole(role.ID, req.Perms)
	if err != nil {
		logUtils.Errorf("添加权限到角色错误", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

func (r *RoleRepo) IsAdminRole(id uint) (bool, error) {
	role, err := r.FindById(id)
	if err != nil {
		return false, err
	}
	return role.Name == serverConsts.AdminRoleName, nil
}

func (r *RoleRepo) FindById(id uint) (v1.RoleResp, error) {
	role := v1.RoleResp{}
	err := r.DB.Model(&model.SysRole{}).Where("id = ?", id).First(&role).Error
	if err != nil {
		logUtils.Errorf("根据id查询角色错误", zap.String("错误:", err.Error()))
		return role, err
	}
	return role, nil
}

func (r *RoleRepo) DeleteById(id uint) error {
	if b, err := r.IsAdminRole(id); err != nil {
		return err
	} else if b {
		return errors.New("不能删除超级管理员")
	}
	err := r.DB.Unscoped().Delete(&model.SysRole{}, id).Error
	if err != nil {
		logUtils.Errorf("删除角色错误", zap.String("错误:", err.Error()))
		return err
	}
	return nil
}

func (r *RoleRepo) FindInId(ids []string) (roles []v1.RoleResp, error error) {
	err := r.DB.Model(&model.SysRole{}).Where("id in ?", ids).Find(&roles).Error
	if err != nil {
		logUtils.Errorf("通过ids查询角色错误", zap.String("错误:", err.Error()))
		return
	}
	return
}

// AddPermForRole
func (r *RoleRepo) AddPermForRole(id uint, perms [][]string) error {
	roleId := strconv.FormatUint(uint64(id), 10)
	oldPerms := casbin.GetPermissionsForUser(roleId)
	_, err := casbin.Instance().RemovePolicies(oldPerms)
	if err != nil {
		logUtils.Errorf("add policy err: %+v", zap.String("错误:", err.Error()))
		return err
	}

	if len(perms) == 0 {
		logUtils.Debug("没有权限")
		return nil
	}
	var newPerms [][]string
	for _, perm := range perms {
		newPerms = append(newPerms, append([]string{roleId}, perm...))
	}
	logUtils.Debugf("添加权限到角色", myZap.Strings("新权限", newPerms))
	_, err = casbin.Instance().AddPolicies(newPerms)
	if err != nil {
		logUtils.Errorf("add policy err: %+v", zap.String("错误:", err.Error()))
		return err
	}

	return nil
}

func (r *RoleRepo) GetRoleIds() ([]uint, error) {
	var roleIds []uint
	err := r.DB.Model(&model.SysRole{}).Select("id").Find(&roleIds).Error
	if err != nil {
		return roleIds, fmt.Errorf("获取角色ids错误 %w", err)
	}
	return roleIds, nil
}

func (r *RoleRepo) GetAllRoles() (res []v1.RoleResp, err error) {
	err = r.DB.Model(&model.SysRole{}).Find(&res).Error
	return
}
