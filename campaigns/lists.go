package campaigns

type List struct {
	// ID represents the ID of the list
	ID string `json:"ListID"`
	// Name represents the name of the list
	Name string
}

type Segment struct {
	// ID represents the ID of the segment
	ID string `json:"SegmentID"`
	// ListID represents the ID of the list that the segment is associated with
	ListID string
	// Title is the title of the segment
	Title string
}

// ListsAndSegments represents a grouping of related lists and segments
type ListsAndSegments struct {
	Lists    []List
	Segments []Segment
}
