// Package sqlxyz Pager 提供了一些针对于数据库分页的计算
// 开发者： 无双(luciferlai@qq.com)
package sqlxyz

import (
	"math"
)

// Pager 分页数据结构体
type Pager struct {
	Page        int64       //当前页
	PageSize    int64       //页记录数
	RecordCount int64       //总记录数
	PageCount   int64       //总页数
	FirstPage   int64       //第一页
	BackPage    int64       // 上一页
	NextPage    int64       // 下一页
	LastPage    int64       //最后一页
	Offset      int64       //偏移多少条记录
	Begin       int64       //起始记录序号
	End         int64       //结束记录序号
	Pages       []int64     //分页页码集合
	Data        interface{} //分页数据集合
}

// paging 分页计算
func (p *Pager) paging() []int64 {
	var data []int64
	if p.PageCount <= 1 {
		return data
	}
	var pageNum int64 = 10
	var beginPage, endPage int64
	if p.PageCount <= pageNum {
		beginPage = 1
		endPage = p.PageCount
	} else {
		var pageNumOffset int64 = int64(pageNum / 2)
		if p.Page < pageNumOffset {
			beginPage = 1
			endPage = pageNum
		} else if p.Page >= p.PageCount-pageNumOffset {
			beginPage = p.PageCount - pageNum
			endPage = p.PageCount
		} else {
			beginPage = int64(math.Max(1, float64(p.Page-pageNumOffset)))
			endPage = beginPage + pageNum
		}
	}
	if beginPage != 1 {
		data = append(data, 1)
	}
	for i := beginPage; i <= endPage; i++ {
		data = append(data, i)
	}
	if endPage != p.PageCount {
		data = append(data, p.PageCount)
	}
	return data
}

// SearchPager 搜索分页计算 Page 当前页  Total 总记录数, PageSize 页记录大小
func SearchPager(Page, Total, PageSize int) (NewPage, PageCount, Offset int) {
	if Total <= PageSize {
		return 1, 1, 0
	}
	if Page <= 0 {
		NewPage = 1
	} else {
		NewPage = Page
	}
	if PageSize < 1 {
		if Total > 100 {
			PageSize = 100
		} else {
			PageSize = Total
		}
	}
	PageCount = int(math.Ceil(float64(Total) / float64(PageSize)))
	if NewPage > PageCount {
		NewPage = PageCount
	}
	Offset = (NewPage - 1) * PageSize
	return NewPage, PageCount, Offset
}

// NewPager 新的分页计算 Page 当前页 ReacordCount 总记录数, PageSize 页记录大小
func NewPager(Page, RecordCount, PageSize int64) Pager {
	var p Pager
	if RecordCount <= PageSize {
		p.Page = 1
		p.PageCount = 1
		p.RecordCount = RecordCount
		p.PageSize = RecordCount
		p.Pages = []int64{1}
		p.Offset = 0
		p.FirstPage = 1
		p.LastPage = 1
		p.Begin = 1
		return p
	}
	p.Page = Page
	p.FirstPage = 1

	if Page < 0 {
		p.Page = 1
	}

	if p.Page > 1 {
		p.BackPage = p.Page - 1
	}
	if PageSize < 1 {
		if RecordCount > 100 {
			PageSize = 100
		} else {
			PageSize = RecordCount
		}
	}
	p.RecordCount = RecordCount
	p.PageSize = PageSize

	p.PageCount = int64(math.Ceil(float64(p.RecordCount) / float64(p.PageSize)))
	if p.Page > p.PageCount {
		p.Page = p.PageCount
	}
	if p.PageCount > p.Page {
		p.LastPage = p.PageCount
	}
	if p.Page < p.PageCount {
		p.NextPage = p.Page + 1
	}
	p.Pages = p.paging()
	p.Offset = (p.Page - 1) * p.PageSize
	p.Begin = p.Offset + 1
	p.End = p.Offset + p.PageSize
	return p
}
