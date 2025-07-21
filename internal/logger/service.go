package logger

import (
	"context"
	"sync"
	"time"

	"github.com/ArshiAbolghasemi/dom-cobb/internal/database/mongodb"
	"github.com/ArshiAbolghasemi/dom-cobb/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

type IService interface {
	Log(entry *LogEntry) error
	LogBatch(entries []*LogEntry) error
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

func (s *Service) Log(entry *LogEntry) error {
	writeTimeout, err := GetWriteTimeOut()
	if err != nil {
		panic("Faile to get mongo write time out: " + err.Error())
	}
	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout*time.Second)
	defer cancel()

	_, err = s.collection.InsertOne(ctx, entry)
	return err
}

func (s *Service) LogBatch(entries []*LogEntry) error {
	if len(entries) == 0 {
		return nil
	}	

	writeTimeout, err := GetWriteTimeOut()
	if err != nil {
		panic("Failed to get mongo write time out: " + err.Error())
	}

	ctx, cancel := context.WithTimeout(context.Background(), writeTimeout*time.Second)
	defer cancel()

	_, err = s.collection.InsertMany(ctx, utils.ToAnySlice(entries))
	return err
}
