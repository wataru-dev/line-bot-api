package model

type Session struct {
	Role      string `firestore:"role"`
	Content   string `firestore:"content"`
	Timestamp int64  `firestore:"timestamp"`
}