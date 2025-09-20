package model

import "time"

type Session struct {
	Role      string    `firestore:"role"`
	Content   string    `firestore:"content"`
	Timestamp int64     `firestore:"timestamp"`
	ExpireAt  time.Time `firestore:"expireAt"`
}