package config

import "os"

type Environment struct {
	LineToken string
	ReplyUri  string
}

func SetEnvironment() *Environment {
	return &Environment{
		LineToken: os.Getenv("LINE_TOKEN"),
		ReplyUri:  os.Getenv("REPLY_URI"),
	}
}
