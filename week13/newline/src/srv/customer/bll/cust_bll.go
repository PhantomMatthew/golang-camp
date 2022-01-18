package bll

import (
	"newline.com/newline/src/srv/customer/dao"
)

// Dao dao.
type CustBll struct {
	custBasicInfoDao *dao.CustBasicInfoDao

	custOrderRecordDao *dao.CustOrderRecordDao
}

// New new a bll.
func NewCustBll() (d *CustBll) {
	d = &CustBll{

		custBasicInfoDao: dao.NewCustBasicInfoDao(),

		custOrderRecordDao: dao.NewCustOrderRecordDao(),
	}
	return
}
