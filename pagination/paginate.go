package pagination

import (
	"fmt"
)

const PSIZE int = 5

type Page struct {
	PageNumber int
	Active     string
	Link       string
}

func Pager(page_no int, max_len int) []Page {
	var pages []Page

	if max_len == 0 || max_len < PSIZE {
		pages = append(pages, Page{Active: "page-item active", Link: fmt.Sprint("/?starting_page=", 1)})
		return pages
	}

	total_pages := int(max_len / PSIZE)

	if max_len%PSIZE != 0 {
		total_pages += 1
	}

	fmt.Println("Max length, Page size:", max_len, PSIZE)
	fmt.Println("Total pages: ", total_pages)
	fmt.Println("Values: ", max_len/PSIZE, max_len%PSIZE)

	for i := 1; i <= total_pages; i++ {
		if i == page_no {
			pages = append(pages, Page{PageNumber: i, Active: "page-item active", Link: fmt.Sprint("/?starting_page=", i)})
			fmt.Printf("%+v\n", Page{Active: "page-item active", Link: fmt.Sprint("/?starting_page=", i)})
		} else {
			pages = append(pages, Page{PageNumber: i, Active: "page-item", Link: fmt.Sprint("/?starting_page=", i)})
		}
	}
	// fmt.Print(pages)
	return pages
}
