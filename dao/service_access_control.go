package dao

import (
	"go_gateway/common"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// AccessControl
type AccessControl struct {
	ID                int64  `json:"id" gorm:"primary_key"`
	ServiceID         int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	OpenAuth          int    `json:"open_auth" gorm:"column:open_auth" description:"是否开启权限 1=开启"`
	BlackList         string `json:"black_list" gorm:"column:black_list" description:"黑名单ip	"`
	WhiteList         string `json:"white_list" gorm:"column:white_list" description:"白名单ip	"`
	OpenApiWhiteList  int    `json:"open_api_white_list" gorm:"column:open_api_white_list" description:"是否开启api白名单 它依赖于open_auth是否开启JTW校验	"`
	OpenWhiteList     int    `json:"open_white_list" gorm:"column:open_white_list" description:" 是否开启IP白名单"`
	OpenBlackList     int    `json:"open_black_list" gorm:"column:open_black_list" description:" 是否开启IP黑名单"`
	ApiWhiteList      string `json:"api_white_list" gorm:"column:api_white_list" description:"api白名单"`
	ClientIPFlowLimit int    `json:"clientip_flow_limit" gorm:"column:clientip_flow_limit" description:"客户端ip限流	"`
	ServiceFlowLimit  int    `json:"service_flow_limit" gorm:"column:service_flow_limit" description:"服务端限流	"`
}

// TableName
func (t *AccessControl) TableName() string {
	return "gateway_service_access_control"
}

// Find
func (t *AccessControl) Find(c *gin.Context, tx *gorm.DB, search *AccessControl) (*AccessControl, error) {
	model := &AccessControl{}
	err := tx.SetCtx(common.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save
func (t *AccessControl) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(common.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// Del
func (t *AccessControl) Del(c *gin.Context, tx *gorm.DB) error {
	return tx.SetCtx(common.GetGinTraceContext(c)).Exec("DELETE FROM "+t.TableName()+" WHERE `service_id` = ? ", t.ServiceID).Error
}

// ListBYServiceID
func (t *AccessControl) ListBYServiceID(c *gin.Context, tx *gorm.DB, serviceID int64) ([]AccessControl, int64, error) {
	var list []AccessControl
	var count int64
	query := tx.SetCtx(common.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("`id`,`service_id`,`open_api_white_list`,`open_white_list`,`open_black_list`,`api_white_list`,`open_auth`,`black_list`,`white_list`,`white_host_name`,`clientip_flow_limit`,`service_flow_limit`")
	query = query.Where("`service_id`=?", serviceID)
	err := query.Order("`id` desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}
