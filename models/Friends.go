package models

import "github.com/google/uuid"

type Friends struct {
	FriendsID    uuid.UUID `json:"friendsID" db:"id"`
	UserSender   uuid.UUID `json:"userSender" db:"userSender"`
	UserReceiver uuid.UUID `json:"userReceiver" db:"userReceiver"`
	StatusID     int       `json:"statusID" db:"status_id"`
}
