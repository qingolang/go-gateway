package dao

import (
	"go_gateway/common"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

type GRPCRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id	"`
	Port           int    `json:"port" gorm:"column:port" description:"端口	"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue"`
}

// TableName
func (t *GRPCRule) TableName() string {
	return "gateway_service_grpc_rule"
}

// Find
func (t *GRPCRule) Find(c *gin.Context, tx *gorm.DB, search *GRPCRule) (*GRPCRule, error) {
	model := &GRPCRule{}
	err := tx.SetCtx(common.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save
func (t *GRPCRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(common.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// ListByServiceID
func (t *GRPCRule) ListByServiceID(c *gin.Context, tx *gorm.DB, serviceID int64) ([]GRPCRule, int64, error) {
	var list []GRPCRule
	var count int64
	query := tx.SetCtx(common.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("`id`,`service_id`,`port`,`header_transfor`")
	query = query.Where("service_id=?", serviceID)
	err := query.Order("id desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}
