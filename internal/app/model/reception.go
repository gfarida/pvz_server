package model

import "time"

type ReceptionStatus string

const (
	InProgress ReceptionStatus = "in_progress"
	Closed     ReceptionStatus = "close"
)

var AllowedReceptionStatuses = map[ReceptionStatus]bool{
	InProgress: true,
	Closed:     true,
}

type Reception struct {
	ID       string          `json:"id"`
	DateTime time.Time       `json:"dateTime"`
	PvzID    string          `json:"pvzId"`
	Status   ReceptionStatus `json:"status"`
}
