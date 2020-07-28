package clients

// Segment represents a list segment.
type Segment struct {
	// ID segment id.
	ID string `json:"SegmentID"`
	// Title segment title.
	Title string
	// ListID the ID of the list the segment belongs to.
	ListID string
}
