package dao

import (
	"go_gateway/common"
	"go_gateway/common/lib"
	"go_gateway/dto"
	"net/http/httptest"
	"sync"
	"time"

	"github.com/e421083458/gorm"
	"github.com/gin-gonic/gin"
)

var APPManagerHandler *APPManager

func init() {
	APPManagerHandler = NewAPPManager()
}

// NewAPPManager
func NewAPPManager() *APPManager {
	return &APPManager{
		APPMap:   map[string]*APP{},
		APPSlice: []*APP{},
		Locker:   sync.RWMutex{},
		init:     sync.Once{},
	}
}

// APP
type APP struct {
	ID        int64     `json:"id" gorm:"primary_key"`
	APPID     string    `json:"app_id" gorm:"column:app_id" description:"租户id	"`
	Name      string    `json:"name" gorm:"column:name" description:"租户名称	"`
	Secret    string    `json:"secret" gorm:"column:secret" description:"密钥"`
	WhiteIPS  string    `json:"white_ips" gorm:"column:white_ips" description:"ip白名单，支持前缀匹配"`
	QPD       int64     `json:"qpd" gorm:"column:qpd" description:"日请求量限制"`
	QPS       int64     `json:"qps" gorm:"column:qps" description:"每秒请求量限制"`
	CreatedAt time.Time `json:"create_at" gorm:"column:create_at" description:"添加时间	"`
	UpdatedAt time.Time `json:"update_at" gorm:"column:update_at" description:"更新时间"`
	IsDelete  int8      `json:"is_delete" gorm:"column:is_delete" description:"是否已删除；0：否；1：是"`
}

// TableName
func (t *APP) TableName() string {
	return "gateway_app"
}

// Find
func (t *APP) Find(c *gin.Context, tx *gorm.DB, search *APP) (*APP, error) {
	model := &APP{}
	err := tx.SetCtx(common.GetGinTraceContext(c)).Where(search).Find(model).Error
	return model, err
}

// Save
func (t *APP) Save(c *gin.Context, tx *gorm.DB) error {
	if err := tx.SetCtx(common.GetGinTraceContext(c)).Save(t).Error; err != nil {
		return err
	}
	return nil
}

// APPList
func (t *APP) APPList(c *gin.Context, tx *gorm.DB, params *dto.APPListInput) ([]APP, int64, error) {
	var list []APP
	var count int64
	pageNo := params.PageNo
	pageSize := params.PageSize

	//limit offset,pagesize
	offset := (pageNo - 1) * pageSize
	query := tx.SetCtx(common.GetGinTraceContext(c))
	query = query.Table(t.TableName()).Select("`id`,`app_id`,`name`,`secret`,`white_ips`,`qpd`,`qps`,`create_at`,`update_at`,`is_delete`")
	query = query.Where("`is_delete` = ?", 0)
	if params.Info != "" {
		query = query.Where(" ( `name` like ? or `app_id` like ?)", "%"+params.Info+"%", "%"+params.Info+"%")
	}
	err := query.Limit(pageSize).Offset(offset).Order("`id` desc").Find(&list).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, 0, err
	}
	errCount := query.Count(&count).Error
	if errCount != nil {
		return nil, 0, err
	}
	return list, count, nil
}

// APPManager
type APPManager struct {
	APPMap   map[string]*APP
	APPSlice []*APP
	Locker   sync.RWMutex
	init     sync.Once
	err      error
}

// GetAppList
func (s *APPManager) GetAppList() []*APP {
	return s.APPSlice
}

// LoadOnce
func (s *APPManager) LoadOnce() error {
	s.init.Do(func() {
		appInfo := &APP{}
		c, _ := gin.CreateTestContext(httptest.NewRecorder())
		tx, err := lib.GetGormPool("default")
		if err != nil {
			s.err = err
			return
		}
		params := &dto.APPListInput{PageNo: 1, PageSize: 99999}
		list, _, err := appInfo.APPList(c, tx, params)
		if err != nil {
			s.err = err
			return
		}
		s.Locker.Lock()
		for _, listItem := range list {
			tmpItem := listItem
			s.APPMap[listItem.APPID] = &tmpItem
			s.APPSlice = append(s.APPSlice, &tmpItem)
		}
		s.Locker.Unlock()
	})
	return s.err
}
