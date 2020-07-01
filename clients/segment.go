package clients

// Segment represents a list segment.
type Segment struct {
	// Id segment id.
	Id string `json:"SegmentID"`
	// Title segment title.
	Title string
	// ListId the Id of the list the segment belongs to.
	ListId string `json:"ListID"`
}
