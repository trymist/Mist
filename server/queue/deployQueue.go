package queue

import (
	"context"
	"fmt"
	"sync"

	"github.com/corecollectives/mist/models"
	"github.com/rs/zerolog/log"
	"gorm.io/gorm"
)

type Queue struct {
	jobs   chan int64
	ctx    context.Context
	cancel context.CancelFunc
	wg     sync.WaitGroup
}

var queue *Queue

func NewQueue(buffer int, db *gorm.DB) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		jobs: make(chan int64, buffer),

		ctx:    ctx,
		cancel: cancel,
	}
	q.StartWorker(db)
	queue = q
	return q

}

func GetQueue() *Queue {
	return queue
}

func (q *Queue) StartWorker(db *gorm.DB) {
	q.wg.Add(1)
	go func() {
		defer q.wg.Done()
		for id := range q.jobs {
			status, err := models.GetDeploymentStatus(id)
			if err != nil {
				log.Error().Err(err).Msg("Failed to get deployment status")
				continue
			}
			if status == "stopped" {
				log.Info().Msgf("Deployment %d has been stopped before processing, skipping", id)
				continue
			}
			q.HandleWork(id, db)
		}

	}()
}

func (q *Queue) AddJob(Id int64) error {
	select {
	case q.jobs <- Id:
		return nil
	case <-q.ctx.Done():
		return fmt.Errorf("queue is closed")
	default:
		return fmt.Errorf("queue is full")
	}
}

func (q *Queue) Close() {
	q.cancel()
	close(q.jobs)
	q.wg.Wait()
	log.Info().Msg("Deployment queue closed")
}
