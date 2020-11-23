package order

// Page represents a single page of paginated records.
type Page struct {
	// OrderDirection the order in which the results were sorted.
	OrderDirection Direction
	// Number the current page number.
	Number int `json:"PageNumber"`
	// Size the original page size.
	Size int `json:"PageSize"`
	// Records number of records on this page.
	Records int `json:"RecordsOnThisPage"`
	// Total the total number of records.
	Total int `json:"TotalNumberOfRecords"`
	// NumberOfPages the total number of pages.
	NumberOfPages int
}
