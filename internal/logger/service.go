package logger

import (
	"context"
	"sync"
	"time"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/database/mongodb"
	"go.mongodb.org/mongo-driver/mongo"
)

type IService interface {
	Log(message string, metadata map[string]any) error
}

var (
	service     IService
	onceService sync.Once
)

func NewService() IService {
	onceService.Do(func() {
		collection, err := GetCollection()
		if err != nil {
			panic("Failed to get logger collection name: " + err.Error())
		}
		service = &Service{
			collection: mongodb.GetCollection(collection),
		}
	})
	return service
}

type Service struct {
	collection *mongo.Collection
}

func (s *Service) Log(message string, metadata map[string]any) error {
	entry := &LogEntry{
		Message:   message,
		Metadata:  metadata,
		Timestamp: time.Now(),
	}
	writeTimeout, err := GetWriteTimeOut()
	if err != nil {
		panic("Faile to get mongo write time out: " + err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout*time.Second)
	defer cancel()

	_, err = s.collection.InsertOne(ctx, entry)
	return err
}
