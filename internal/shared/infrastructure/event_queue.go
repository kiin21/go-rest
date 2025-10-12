package infrastructure

import (
	"context"
	"log"
	"sync"
	"time"
)

// SyncEvent represents a sync event
type SyncEvent struct {
	Type      string
	Domain    string
	Data      interface{}
	Timestamp time.Time
	Retries   int
}

// SyncEventQueue manages async sync events with reliability
type SyncEventQueue struct {
	events      chan *SyncEvent
	deadLetter  chan *SyncEvent
	workers     int
	maxRetries  int
	retryConfig RetryConfig
	wg          sync.WaitGroup
	ctx         context.Context
	cancel      context.CancelFunc
	handler     func(context.Context, *SyncEvent) error
}

// NewSyncEventQueue creates a new event queue
func NewSyncEventQueue(bufferSize, workers, maxRetries int, handler func(context.Context, *SyncEvent) error) *SyncEventQueue {
	ctx, cancel := context.WithCancel(context.Background())

	return &SyncEventQueue{
		events:      make(chan *SyncEvent, bufferSize),
		deadLetter:  make(chan *SyncEvent, bufferSize),
		workers:     workers,
		maxRetries:  maxRetries,
		retryConfig: DefaultRetryConfig(),
		ctx:         ctx,
		cancel:      cancel,
		handler:     handler,
	}
}

// Start starts the queue workers
func (q *SyncEventQueue) Start() {
	log.Printf("[SyncEventQueue] Starting %d workers...", q.workers)

	for i := 0; i < q.workers; i++ {
		q.wg.Add(1)
		go q.worker(i)
	}

	// Start dead letter queue processor
	q.wg.Add(1)
	go q.deadLetterProcessor()
}

// Stop gracefully stops the queue
func (q *SyncEventQueue) Stop() {
	log.Printf("[SyncEventQueue] Stopping queue...")
	q.cancel()
	close(q.events)
	q.wg.Wait()
	close(q.deadLetter)
	log.Printf("[SyncEventQueue] Queue stopped")
}

// Enqueue adds an event to the queue
func (q *SyncEventQueue) Enqueue(event *SyncEvent) {
	select {
	case q.events <- event:
		// Event queued successfully
	case <-q.ctx.Done():
		log.Printf("[SyncEventQueue] Queue is shutting down, event dropped: %+v", event)
	default:
		// Queue is full, log warning
		log.Printf("[SyncEventQueue] WARNING: Queue is full, event dropped: %+v", event)
	}
}

// worker processes events from the queue
func (q *SyncEventQueue) worker(id int) {
	defer q.wg.Done()

	log.Printf("[SyncEventQueue] Worker %d started", id)

	for {
		select {
		case event, ok := <-q.events:
			if !ok {
				log.Printf("[SyncEventQueue] Worker %d: Channel closed, exiting", id)
				return
			}

			// Process event with retry
			err := RetryWithBackoff(q.ctx, q.retryConfig, func() error {
				return q.handler(q.ctx, event)
			}, event.Type)

			if err != nil {
				event.Retries++
				if event.Retries >= q.maxRetries {
					log.Printf("[SyncEventQueue] Worker %d: Event failed after max retries, moving to dead letter: %+v", id, event)
					q.deadLetter <- event
				} else {
					log.Printf("[SyncEventQueue] Worker %d: Event failed, retrying (%d/%d): %v", id, event.Retries, q.maxRetries, err)
					q.events <- event // Re-queue
				}
			} else {
				log.Printf("[SyncEventQueue] Worker %d: Event processed successfully: %s - %s", id, event.Type, event.Domain)
			}

		case <-q.ctx.Done():
			log.Printf("[SyncEventQueue] Worker %d: Context cancelled, exiting", id)
			return
		}
	}
}

// deadLetterProcessor handles failed events
func (q *SyncEventQueue) deadLetterProcessor() {
	defer q.wg.Done()

	log.Printf("[SyncEventQueue] Dead letter processor started")

	for {
		select {
		case event, ok := <-q.deadLetter:
			if !ok {
				log.Printf("[SyncEventQueue] Dead letter queue closed, exiting")
				return
			}

			log.Printf("[SyncEventQueue] DEAD LETTER: %s - %s (retries: %d, timestamp: %v)",
				event.Type, event.Domain, event.Retries, event.Timestamp)

			// TODO: Persist to database or external queue for manual recovery

		case <-q.ctx.Done():
			log.Printf("[SyncEventQueue] Dead letter processor: Context cancelled, exiting")
			return
		}
	}
}
