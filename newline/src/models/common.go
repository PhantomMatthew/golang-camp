package models

import "time"

// Response Response
type Response struct {
	Result  int         `json:"result"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

// Model Model
type Model struct {
	ID        uint32     `jpath:"id" gorm:"primary_key"   comment:"自增主键"`
	CreatedAt time.Time  `jpath:"created_at" gorm:"type:timestamp" comment:"创建于"`
	UpdatedAt time.Time  `jpath:"updated_at" gorm:"type:timestamp" comment:"更新于"`
	DeletedAt *time.Time `jpath:"deleted_at" gorm:"type:timestamp" commen:"删除于"`
	Remark    string     `jpath:"Remark" gorm:"type:varchar(255)"   comment:"删除于"`
}

// PaginationParam 分页查询条件
type PaginationParam struct {
	PageIndex uint32 // 页索引
	PageSize  uint32 // 页大小
}

// PaginationResult 分页查询结果
type PaginationResult struct {
	Total uint32 // 总数据条数
}
