package agglayer

import (
	"fmt"
	"testing"
)

func TestOffsetLimitToPagination(t *testing.T) {
	arr := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	fmt.Println(OffsetLimitToPagination(arr, 5, 4))
	fmt.Println(OffsetLimitToPagination(arr, 1, 7))
	fmt.Println(OffsetLimitToPagination(arr, 8, 2))
	fmt.Println(OffsetLimitToPagination(arr, 5, 5))
	fmt.Println(OffsetLimitToPagination(arr, 10, 10))
	fmt.Println(OffsetLimitToPagination(arr, 100, 100))
	fmt.Println(OffsetLimitToPagination([]int{}, 100, 100))
}

func OffsetLimitToPagination(arr []int, offset, limit int) []int {
	pageSize := limit
	pageNumber := offset / limit
	skipItems := offset % limit

	totalPages := len(arr) / pageSize
	if len(arr)%pageSize > 0 {
		totalPages++
	}

	pages := make([][]int, totalPages)

	if pageNumber >= totalPages {
		return []int{}
	}
	items := []int{}

	for i := 0; i < totalPages; i++ {
		start := pageSize * i
		end := pageSize * (i + 1)
		if end > len(arr)-1 {
			end = len(arr)
		}
		pages[i] = arr[start:end]
	}

	start := skipItems
	end := pageSize
	page := pages[pageNumber]
	items = append(items, page[start:end]...)

	if skipItems > 0 && pageNumber+1 < totalPages {
		start := 0
		end := skipItems
		page := pages[pageNumber+1]
		pageItensCount := len(page)
		if pageItensCount < end {
			end = pageItensCount
		}

		items = append(items, page[start:end]...)
	}

	return items
}
