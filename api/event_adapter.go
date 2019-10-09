package api

import (
	"github.com/golang/protobuf/ptypes"
	"golang_lesson17/internal/domain/model"
)

func Convert(event *Event) (*model.Event, error) {

	time, err := ptypes.Timestamp(event.Date)

	if err != nil {
		return nil, err
	}

	return &model.Event{
		Id:          event.MsgId,
		Type:        model.EventType(event.Type),
		Date:        time,
		Description: event.Description,
	}, nil
}
