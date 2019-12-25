package sqlxyz

import (
	"fmt"
	"testing"
)

func TestNewPager(t *testing.T) {
	fmt.Println(fmt.Sprintf("%#v", NewPager(0, 2, 0)))
}

func TestSearchPagerr(t *testing.T) {
	newPage, pageCount, offset := SearchPager(0, 2, 0)
	fmt.Println(fmt.Sprintf("newPage:%d, pageCount:%d, offset:%d", newPage, pageCount, offset))
}
