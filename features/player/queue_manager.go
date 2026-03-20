package player

import (
	"sync"

	"github.com/disgoorg/snowflake/v2"
)

type QueueManager struct {
	queues sync.Map // map[snowflake.ID]*Queue
}

func NewQueueManager() *QueueManager {
	return &QueueManager{}
}

func (qm *QueueManager) Get(guildID snowflake.ID) *Queue {
	v, ok := qm.queues.Load(guildID)
	if !ok {
		q := NewQueue()
		qm.queues.Store(guildID, q)
		return q
	}
	q, ok := v.(*Queue)
	if !ok {
		q = NewQueue()
		qm.queues.Store(guildID, q)
		return q
	}
	return q
}

func (qm *QueueManager) Delete(guildID snowflake.ID) {
	qm.queues.Delete(guildID)
}
