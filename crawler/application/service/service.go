package service

import "context"

type MeetingsService interface {
	Start(meeting *model.Meeting, ctx context.Context) (*model.Meeting, error)
}
