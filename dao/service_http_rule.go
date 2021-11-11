package dao

import (
	"go_gateway/common"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

// HTTPRule
type HTTPRule struct {
	ID             int64  `json:"id" gorm:"primary_key"`
	ServiceID      int64  `json:"service_id" gorm:"column:service_id" description:"服务id"`
	RuleType       int    `json:"rule_type" gorm:"column:rule_type" description:"匹配类型 domain=域名, url_prefix=url前缀"`
	Rule           string `json:"rule" gorm:"column:rule" description:"type=domain表示域名，type=url_prefix时表示url前缀"`
	NeedHttps      int    `json:"need_https" gorm:"column:need_https" description:"type=支持https 1=支持"`
	NeedWebsocket  int    `json:"need_websocket" gorm:"column:need_websocket" description:"启用websocket 1=启用"`
	NeedStripUri   int    `json:"need_strip_uri" gorm:"column:need_strip_uri" description:"启用strip_uri 1=启用"`
	UrlRewrite     string `json:"url_rewrite" gorm:"column:url_rewrite" description:"url重写功能，每行一个	"`
	HeaderTransfor string `json:"header_transfor" gorm:"column:header_transfor" description:"header转换支持增加(add)、删除(del)、修改(edit) 格式: add headname headvalue	"`
}

// TableName
func (t *HTTPRule) TableName() string {
	return "gateway_service_http_rule"
}

// Find
func (t *HTTPRule) Find(c *gin.Context, tx *gorm.DB, search *HTTPRule) (*HTTPRule, error) {
	model := &HTTPRule{}
	err := tx.SetCtx(common.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save
func (t *HTTPRule) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(common.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// ListByServiceID
func (t *HTTPRule) ListByServiceID(c *gin.Context, tx *gorm.DB, serviceID int64) ([]HTTPRule, int64, error) {
	var list []HTTPRule
	var count int64
	query := tx.SetCtx(common.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("`id`,`service_id`,`rule_type`,`rule`,`need_https`,`need_websocket`,`need_strip_uri`,`url_rewrite`,`header_transfor`")
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
