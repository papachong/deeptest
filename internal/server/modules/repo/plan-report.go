package repo

import (
	"fmt"
	v1 "github.com/aaronchen2k/deeptest/cmd/server/v1/domain"
	"github.com/aaronchen2k/deeptest/internal/pkg/consts"
	"github.com/aaronchen2k/deeptest/internal/server/core/dao"
	model "github.com/aaronchen2k/deeptest/internal/server/modules/model"
	"github.com/aaronchen2k/deeptest/pkg/domain"
	logUtils "github.com/aaronchen2k/deeptest/pkg/lib/log"
	"gorm.io/gorm"
	"strconv"
	"strings"
	"time"
)

type PlanReportRepo struct {
	DB                 *gorm.DB            `inject:""`
	LogRepo            *LogRepo            `inject:""`
	UserRepo           *UserRepo           `inject:""`
	ScenarioReportRepo *ScenarioReportRepo `inject:""`
	ProjectRepo        *ProjectRepo        `inject:""`
	ScenarioRepo       *ScenarioRepo       `inject:""`
	PlanRepo           *PlanRepo           `inject:""`
}

func (r *PlanReportRepo) Paginate(req v1.PlanReportReqPaginate, projectId int) (data _domain.PageData, err error) {
	req.Order = "desc"
	req.Field = "biz_plan_report.created_at"
	timeLayout := "2006-01-02 15:04:05"
	var count int64

	db := r.DB.Model(&model.PlanReport{}).
		Joins("LEFT JOIN sys_user u ON biz_plan_report.create_user_id=u.id").
		Select("biz_plan_report.*, u.name create_user_name").
		Where("biz_plan_report.project_id = ? AND NOT biz_plan_report.deleted", projectId)

	if req.Keywords != "" {
		db = db.Where("biz_plan_report.name LIKE ?", fmt.Sprintf("%%%s%%", req.Keywords))
	}
	if req.PlanId != 0 {
		db = db.Where("biz_plan_report.plan_id = ?", req.PlanId)
	}
	if req.CreateUserId != "" {
		db = db.Where("biz_plan_report.create_user_id IN (?)", strings.Split(req.CreateUserId, ","))
	}

	if req.ExecuteStartTime != 0 {
		db = db.Where("biz_plan_report.start_time > ?", time.Unix(req.ExecuteStartTime/1000, 0).Format(timeLayout))
	}
	if req.ExecuteEndTime != 0 {
		db = db.Where("biz_plan_report.end_time < ?", time.Unix(req.ExecuteEndTime/1000, 0).Format(timeLayout))
	}

	err = db.Count(&count).Error
	if err != nil {
		logUtils.Errorf("count report error %s", err.Error())
		return
	}

	results := make([]*model.PlanReportDetail, 0)

	err = db.
		Scopes(dao.PaginateScope(req.Page, req.PageSize, req.Order, req.Field)).
		Find(&results).Error
	if err != nil {
		logUtils.Errorf("query report error %s", err.Error())
		return
	}

	data.Populate(results, count, req.Page, req.PageSize)

	return
}

func (r *PlanReportRepo) Get(id uint) (report model.PlanReportDetail, err error) {
	err = r.DB.Model(model.PlanReport{}).
		Select("biz_plan_report.*, e.name exec_env, u.name exec_user_name").
		Joins("LEFT JOIN biz_environment e ON biz_plan_report.exec_env_id=e.id").
		Joins("LEFT JOIN sys_user u ON biz_plan_report.create_user_id=u.id").
		Where("biz_plan_report.id = ?", id).First(&report).Error
	if err != nil {
		logUtils.Errorf("find report by id error %s", err.Error())
		return
	}

	scenarioReports, err := r.ScenarioReportRepo.GetReportsByPlanReportId(report.ID)
	report.ScenarioReports = scenarioReports

	createUserName, _ := r.GetCreateUserName(report)
	report.CreateUserName = createUserName
	//root, err := r.getLogTree(report)
	//report.Logs = root.Logs

	return
}

func (r *PlanReportRepo) GetCreateUserName(report model.PlanReportDetail) (name string, err error) {
	if report.PlanId == 0 {
		err = r.DB.Model(model.ScenarioReport{}).
			Select("u.name").
			Joins("LEFT JOIN biz_scenario s ON biz_scenario_report.scenario_id=s.id").
			Joins("LEFT JOIN sys_user u ON s.create_user_id=u.id").
			Where("biz_scenario_report.plan_report_id=?", report.ID).
			Find(&name).Error
	} else {
		err = r.DB.Model(model.Plan{}).
			Select("u.name").
			Joins("LEFT JOIN sys_user u ON biz_plan.create_user_id=u.id").
			Where("biz_plan.id=?", report.PlanId).
			Find(&name).Error
	}
	return
}

func (r *PlanReportRepo) Create(result *model.PlanReport) (err error) {

	err = r.DB.Model(&model.PlanReport{}).Create(result).Error
	if err != nil {
		logUtils.Errorf("create plan report error %s", err.Error())
		return
	}
	if err = r.UpdateSerialNumber(result.ID, result.ProjectId); err != nil {
		logUtils.Errorf("update plan report serial number error %s", err.Error())
		return
	}

	return
}

func (r *PlanReportRepo) DeleteById(id uint) (err error) {
	err = r.DB.Model(&model.PlanReport{}).Where("id = ?", id).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete report by id error %s", err.Error())
		return
	}

	err = r.ScenarioReportRepo.BatchDelete(id)
	if err != nil {
		logUtils.Errorf("delete report's logs by id error %s", err.Error())
		return
	}

	return
}

func (r *PlanReportRepo) UpdateStatus(progressStatus consts.ProgressStatus, resultStatus consts.ResultStatus, scenarioId uint) (
	err error) {

	values := map[string]interface{}{
		"progress_status": progressStatus,
		"result_status":   resultStatus,
	}
	err = r.DB.Model(&model.ScenarioReport{}).
		Where("report_id = ? AND progress_status = ?", scenarioId, consts.InProgress).
		Updates(values).Error

	return
}

func (r *PlanReportRepo) UpdateResult(report model.ScenarioReport) (err error) {
	values := map[string]interface{}{
		"pass_num":        report.PassRequestNum,
		"fail_num":        report.FailRequestNum,
		"start_time":      report.StartTime,
		"end_time":        report.EndTime,
		"duration":        report.Duration,
		"progress_status": consts.End,
		"result_status":   report.ResultStatus,
	}
	err = r.DB.Model(&model.ScenarioReport{}).
		Where("id = ?", report.ID).
		Updates(values).Error

	return
}

func (r *PlanReportRepo) ResetResult(result model.ScenarioReport) (err error) {
	values := map[string]interface{}{
		"name":       result.Name,
		"start_time": result.StartTime,
	}
	err = r.DB.Model(&result).Where("id = ?", result.ID).Updates(values).Error
	if err != nil {
		logUtils.Errorf("update report error %s", err.Error())
		return
	}

	return
}

func (r *PlanReportRepo) ClearLogs(resultId uint) (err error) {
	err = r.DB.Model(&model.ExecLogProcessor{}).Where("result_id = ?", resultId).
		Updates(map[string]interface{}{"deleted": true}).Error
	if err != nil {
		logUtils.Errorf("delete logs by result id error %s", err.Error())
		return
	}

	return
}

func (r *PlanReportRepo) FindInProgressResult(scenarioId uint) (result model.ScenarioReport, err error) {
	err = r.DB.Model(&result).
		Where("progress_status =? AND scenario_id = ? AND  not deleted", consts.InProgress, scenarioId).
		First(&result).Error

	return
}

func (r *PlanReportRepo) getLogTree(report model.ScenarioReport) (root model.ExecLogProcessor, err error) {
	logs, err := r.LogRepo.ListByReport(report.ID)
	if err != nil {
		return
	}

	for _, log := range logs {
		if log.ProcessorType == consts.ProcessorInterfaceDefault {
			log.InterfaceExtractorsResult, _ = r.listLogExtractors(log.ID)
			log.InterfaceCheckpointsResult, _ = r.listLogCheckpoints(log.ID)
		}
	}

	root = model.ExecLogProcessor{
		Name: report.Name,
	}
	r.makeTree(logs, &root)

	return
}

func (r *PlanReportRepo) makeTree(Data []*model.ExecLogProcessor, parent *model.ExecLogProcessor) { //参数为父节点，添加父节点的子节点指针切片
	children, _ := r.haveChild(Data, parent) //判断节点是否有子节点并返回

	if children != nil {
		parent.Logs = append(parent.Logs, children[0:]...) //添加子节点

		for _, child := range children { //查询子节点的子节点，并添加到子节点
			_, has := r.haveChild(Data, child)
			if has {
				r.makeTree(Data, child) //递归添加节点
			}
		}
	}
}

func (r *PlanReportRepo) haveChild(Data []*model.ExecLogProcessor, node *model.ExecLogProcessor) (children []*model.ExecLogProcessor, yes bool) {
	for _, v := range Data {
		if v.ParentId == node.ID {
			children = append(children, v)
		}
	}

	if children != nil {
		yes = true
	}

	return
}

func (r *PlanReportRepo) listLogExtractors(logId uint) (extractors []model.ExecLogExtractor, err error) {
	err = r.DB.
		Where("log_id =? AND not deleted", logId).
		Find(&extractors).Error

	return
}

func (r *PlanReportRepo) listLogCheckpoints(logId uint) (checkpoints []model.ExecLogCheckpoint, err error) {
	err = r.DB.
		Where("log_id =? AND not deleted", logId).
		Find(&checkpoints).Error

	return
}

func (r *PlanReportRepo) GetLastByPlanId(planId uint) (report model.PlanReport, err error) {
	err = r.DB.Model(&model.PlanReport{}).Where("plan_id = ?", planId).Last(&report).Error
	if err != nil {
		logUtils.Errorf("find plan report by plan_id error %s", err.Error())
		return
	}
	return
}

func (r *PlanReportRepo) GetPlanExecNumber(planId uint) (num int64, err error) {
	err = r.DB.Model(&model.PlanReport{}).Where("plan_id = ?", planId).Count(&num).Error
	if err != nil {
		logUtils.Errorf("find plan report by plan_id error %s", err.Error())
		return
	}
	return
}

func (r *PlanReportRepo) UpdateSerialNumber(id, projectId uint) (err error) {
	var project model.Project
	project, err = r.ProjectRepo.Get(projectId)
	if err != nil {
		return
	}

	err = r.DB.Model(&model.PlanReport{}).Where("id=?", id).Update("serial_number", project.ShortName+"-TR-"+strconv.Itoa(int(id))).Error
	return
}
