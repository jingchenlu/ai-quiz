package utils

import "gorm.io/gorm"

type Page struct {
	PageNum  int `json:"page_num" form:"page_num" default:"1"`
	PageSize int `json:"page_size" form:"page_size" default:"10"`
}

func NewPage(pageNum, pageSize int) Page {
	if pageNum <= 0 {
		pageNum = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	} else if pageSize > 100 {
		pageSize = 100
	}
	return Page{PageNum: pageNum, PageSize: pageSize}
}

// Paginate 通用分页函数
func Paginate(p Page) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		if p.PageNum <= 0 {
			p.PageNum = 1
		}
		switch {
		case p.PageSize > 100:
			p.PageSize = 100
		case p.PageSize <= 0:
			p.PageSize = 10
		}
		offset := (p.PageNum - 1) * p.PageSize
		return db.Offset(offset).Limit(p.PageSize)
	}
}

// PageResult 分页查询结果结构体
type PageResult struct {
	List       interface{} `json:"list"`       // 当前页数据列表
	Total      int64       `json:"total"`      // 总记录数
	Page       int         `json:"page"`       // 当前页码
	PageSize   int         `json:"pageSize"`   // 每页大小
	TotalPages int         `json:"totalPages"` // 总页数
}

// NewPageResult 创建分页结果实例
func NewPageResult(list interface{}, total int64, page, pageSize int) *PageResult {
	totalPages := int(total) / pageSize
	if int(total)%pageSize != 0 {
		totalPages++
	}
	return &PageResult{
		List:       list,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: totalPages,
	}
}
