package campaigns

type List struct {
	ID   string `json:"ListID"`
	Name string
}

type Segment struct {
	ID     string `json:"SegmentID"`
	ListID string
	Title  string
}

type ListsAndSegments struct {
	Lists    []List
	Segments []Segment
}
