package models

import "time"

type Segment struct {
	Sender       string    `json:"sender"`
	Payload      string    `json:"payload"`
	Time         time.Time `json:"time"`
	TotalLength  int       `json:"total_length"`
	SegmentIndex int       `json:"segment_index"`
	Error        bool      `json:"error"`
}

type Response struct {
	Segment Segment `json:"segment"`
	Error   bool    `json:"error"`
}

type ChannelLevel struct {
	ProbabilityError int
	ProbabilityLoss  int
}
