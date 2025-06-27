package logger

import "time"

type LogEntry struct {
	Message   string         `bson:"message" json:"message"`
	Timestamp time.Time      `bson:"timestamp" json:"timestamp"`
	Metadata  map[string]any `bson:"metadata,omitempty" json:"metadata,omitempty"`
}
