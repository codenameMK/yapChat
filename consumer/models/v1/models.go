package v1

import "time"

type MessageStruct struct {
	TimeStamp  time.Time
	Message    string
	UserId     int32
	ReceiverId int32
}
