// Package models model Request
package models

import "github.com/google/uuid"

// Request used to define sender and receiver
type Request struct {
	SenderID   uuid.UUID `json:"senderID"`
	ReceiverID uuid.UUID `json:"receiverID"`
}
