package dao

import (
	"go_gateway/common"
	"go_gateway/dto"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// ServiceInfo
type ServiceInfo struct {
	ID          int64     `json:"id" gorm:"primary_key"`
	LoadType    int       `json:"load_type" gorm:"column:load_type" description:"负载类型 0=http 1=tcp 2=grpc"`
	ServiceName string    `json:"service_name" gorm:"column:service_name" description:"服务名称"`
	ServiceDesc string    `json:"service_desc" gorm:"column:service_desc" description:"服务描述"`
	UpdatedAt   time.Time `json:"create_at" gorm:"column:create_at" description:"更新时间"`
	CreatedAt   time.Time `json:"update_at" gorm:"column:update_at" description:"添加时间"`
	IsDelete    int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
	Status      int8      `json:"status" gorm:"column:status" description:"状态；0：禁用；1：启用"`
}

// TableName
func (t *ServiceInfo) TableName() string {
	return "gateway_service_info"
}

// ServiceDetail
func (t *ServiceInfo) ServiceDetail(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceDetail, error) {
	if search.ServiceName == "" {
		info, err := t.Find(c, tx, search)
		if err != nil {
			return nil, err
		}
		search = info
	}
	HTTPRule := &HTTPRule{ServiceID: search.ID}
	HTTPRule, err := HTTPRule.Find(c, tx, HTTPRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	tcpRule := &TcpRule{ServiceID: search.ID}
	tcpRule, err = tcpRule.Find(c, tx, tcpRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	GRPCRule := &GRPCRule{ServiceID: search.ID}
	GRPCRule, err = GRPCRule.Find(c, tx, GRPCRule)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	accessControl := &AccessControl{ServiceID: search.ID}
	accessControl, err = accessControl.Find(c, tx, accessControl)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}
	loadBalance := &LoadBalance{ServiceID: search.ID}
	loadBalance, err = loadBalance.Find(c, tx, loadBalance)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	detail := &ServiceDetail{
		Info:          search,
		HTTPRule:      HTTPRule,
		TCPRule:       tcpRule,
		GRPCRule:      GRPCRule,
		LoadBalance:   loadBalance,
		AccessControl: accessControl,
	}
	return detail, nil
}

// GroupByLoadType
func (t *ServiceInfo) GroupByLoadType(c *gin.Context, tx *gorm.DB) ([]dto.DashServiceStatItemOutput, error) {
	list := []dto.DashServiceStatItemOutput{}
	query := tx.SetCtx(common.GetGinTraceContext(c))
	if err := query.Table(t.TableName()).Where("`is_delete`=0").Select("`load_type`, count(*) as value").Group("`load_type`").Scan(&list).Error; err != nil {
		return nil, err
	}
	return list, nil
}

// PageList
func (t *ServiceInfo) PageList(c *gin.Context, tx *gorm.DB, param *dto.ServiceListInput) ([]ServiceInfo, int64, error) {
	total := int64(0)
	list := []ServiceInfo{}
	offset := (param.PageNo - 1) * param.PageSize

	query := tx.SetCtx(common.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Where("`is_delete`=0")
	if param.Info != "" {
		query = query.Where("(service_name like ? or service_desc like ?)", "%"+param.Info+"%", "%"+param.Info+"%")
	}

	if param.Status != -1 {
		query = query.Where("`status` = ? ", param.Status)
	}

	if err := query.Limit(param.PageSize).Offset(offset).Order("`id` desc").Find(&list).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	query.Limit(param.PageSize).Offset(offset).Count(&total)
	return list, total, nil
}

// Find
func (t *ServiceInfo) Find(c *gin.Context, tx *gorm.DB, search *ServiceInfo) (*ServiceInfo, error) {
	out := &ServiceInfo{}
	err := tx.SetCtx(common.GetGinTraceContext(c)).Where(search).Find(out).Error
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Save
func (t *ServiceInfo) Save(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(common.GetGinTraceContext(c)).Save(t).Error
}

// Del
func (t *ServiceInfo) Del(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(common.GetGinTraceContext(c)).Exec("DELETE FROM "+t.TableName()+" WHERE `id` = ? ", t.ID).Error
}
