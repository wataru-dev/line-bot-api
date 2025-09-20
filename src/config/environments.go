package config

import "os"

type Environment struct {
	LineToken       string
	ReplyUri        string
	GeminiKey       string
	GoogleProjectID string
}

func SetEnvironment() *Environment {
	return &Environment{
		LineToken:       os.Getenv("LINE_TOKEN"),
		ReplyUri:        os.Getenv("REPLY_URI"),
		GeminiKey:       os.Getenv("GEMINI_KEY"),
		GoogleProjectID: os.Getenv("GOOGLE_PROJECT_ID"),
	}
}
