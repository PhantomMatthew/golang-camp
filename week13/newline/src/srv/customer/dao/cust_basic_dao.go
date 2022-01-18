package dao

import (
	"context"
	"fmt"
	"newline.com/newline/src/models"

	"newline.com/newline/src/srv/customer/components"
	"newline.com/newline/src/srv/customer/model"
)

// Dao dao.
type CustBasicInfoDao struct {
}

// New new a dao.
func NewCustBasicInfoDao() (d *CustBasicInfoDao) {
	d = &CustBasicInfoDao{
		// account memcache
	}
	return
}

//GetByCustId
func (a *CustBasicInfoDao) GetByCustId(ctx context.Context, id uint32) (*model.CustBasicInfo, error) {
	var customerBasicRecord model.CustBasicInfo
	result := components.MainDB.Where("id=?", id).First(&customerBasicRecord)
	return &customerBasicRecord, result.Error
}

//GetByCustId
func (a *CustBasicInfoDao) GetByUnionId(ctx context.Context, unionId string) (*model.CustBasicInfo, error) {
	var customerBasic model.CustBasicInfo
	fmt.Printf("%v, %T\n", &customerBasic)
	result := components.MainDB.Where("union_id=?", unionId).First(&customerBasic)
	return &customerBasic, result.Error
}

//GetList 获取列表
//func (a *CustBasicInfoDao) GetList(ctx context.Context, params model.CustomerBasicParam) (*[]model.CustBasicInfo, error) {
//	var items []model.CustBasicInfo
//	db := components.MainDB
//	if v := params.UnionId; v != "" {
//		db = db.Where("union_id=?", v)
//	}
//	result := db.Find(&items)
//	if err := result.Error; err != nil {
//		return nil, errors.WithStack(err)
//	}
//	return &items, nil
//}
func (a *CustBasicInfoDao) GetList(ctx context.Context, params model.CustomerBasicParam, pageParams *models.PaginationParam) (*[]model.CustBasicInfo, *models.PaginationResult, error) {
	var items []model.CustBasicInfo
	var rowCount uint32
	var pageResults models.PaginationResult
	db := components.MainDB
	if v := params.ID; v > 0 {
		db = db.Where("id=?", v)
	}
	if v := params.UnionId; v != "" {
		db = db.Where("union_id=?", v)
	}
	if v := params.OwnerAccountName; v != "" {
		db = db.Where("owner_account_name=?", v)
	}
	if v := params.OwnerAccountId; v != 0 {
		db = db.Where("owner_account_id=?", v)
	}
	if v := params.Phone; v != "" {
		db = db.Where("phone=?", v)
	}
	if v := params.NickName; v != "" {
		db = db.Where("nick_name=?", v)
	}
	db = db.Order("id DESC")
	if pageParams == nil {
		db = db.Find(&items)
	} else {
		db = db.Offset((pageParams.PageIndex - 1) * pageParams.PageSize).Limit(pageParams.PageSize).Find(&items)
		db.Count(&rowCount)
		pageResults = models.PaginationResult{
			Total: rowCount,
		}
	}
	return &items, &pageResults, db.Error
}

//func (a *CustBasicInfoDao) GetListPagination(ctx context.Context, params model.CustomerBasicParam, pageParams models.PaginationParam) (*[]model.CustBasicInfo, *models.PaginationResult, error) {
//	var items []model.CustOrderRecord
//	var rowCount uint32
//	db := components.MainDB
//	if v := params.ID; v > 0 {
//		db = db.Where("id=?", v)
//	}
//
//	db = db.Order("id DESC")
//
//	db = db.Offset((pageParams.PageIndex - 1) * pageParams.PageSize).Limit(pageParams.PageSize).Find(&items)
//	db.Count(&rowCount)
//	pageResults := models.PaginationResult{
//		Total: rowCount,
//	}
//	return &items, &pageResults, db.Error
//}
// Update 创建数据
func (a *CustBasicInfoDao) Create(ctx context.Context, item *model.CustBasicInfo) (*model.CustBasicInfo, error) {
	result := components.MainDB.Create(&item)
	return item, result.Error
}

// Update 更新数据
func (a *CustBasicInfoDao) Update(ctx context.Context, id uint32, item *model.CustBasicInfo) (*model.CustBasicInfo, error) {
	result := components.MainDB.Model(&item).Where("id=?", id).Omit("id").Updates(&item)
	return item, result.Error
}

// Save 更新数据
func (a *CustBasicInfoDao) Save(ctx context.Context, item *model.CustBasicInfo) (*model.CustBasicInfo, error) {
	result := components.MainDB.Save(item)
	return item, result.Error
}

// Delete 删除数据
func (a *CustBasicInfoDao) Delete(ctx context.Context, id uint32) error {
	result := components.MainDB.Where("id=?", id).Delete(model.CustBasicInfo{})
	return result.Error
}

func (a *CustBasicInfoDao) GetCustChannelInfo(id uint32, unionId string) (*model.CustChannelInfo, error) {
	custChannelInfo := model.CustChannelInfo{}
	sqlWhere := `where c.deleted_at IS NULL `
	params := []interface{}{}
	if v := id; v != 0 {
		sqlWhere = sqlWhere + " AND c.id=?"
		params = append(params, v)
	}
	if v := unionId; v != "" {
		sqlWhere = sqlWhere + " AND c.union_id=?"
		params = append(params, v)
	}

	db := components.MainDB.Raw(`
	select * from cust_basic_info c
	LEFT JOIN cust_pa_info pa  on c.id=pa.cust_id
	LEFT JOIN cust_dz_info dz  on c.id=dz.cust_id
	LEFT JOIN cust_mb_info mb  on c.id=mb.cust_id  `+sqlWhere, params...)

	db.Scan(&custChannelInfo)
	return &custChannelInfo, db.Error
}
