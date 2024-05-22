package models

import "time"

type Segment struct {
	ID           int       `json:"id"`
	Sender       string    `json:"sender"`
	Payload      string    `json:"payload"`
	Time         time.Time `json:"time"`
	TotalLength  int       `json:"total_length"`
	SegmentIndex int       `json:"segment_index"`
}

type Response struct {
	Segment Segment
	Error   bool `json:"error"`
}

type ChannelLevel struct {
	ProbabilityError int
	ProbabilityLoss  int
}
