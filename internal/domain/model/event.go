package model

import "time"

type EventType int

const (
	Undefined EventType = 0 << iota
	Meeting
	Notification
)

type Event struct {
	Id          uint64
	Type        EventType
	Date        time.Time
	Description string
}
