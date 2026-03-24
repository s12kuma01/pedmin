// Copyright (c) 2025-2026 s12kuma01
// SPDX-License-Identifier: MPL-2.0

package model

import (
	"crypto/rand"
	"math/big"
	"sync"

	"github.com/disgoorg/disgolink/v3/lavalink"
	"github.com/disgoorg/snowflake/v2"
)

// PlayerSettings holds per-guild player configuration.
type PlayerSettings struct {
	DefaultVolume *int `json:"default_volume"` // nil = use global default
}

// TrackedMessage tracks a player message for UI updates.
type TrackedMessage struct {
	ChannelID snowflake.ID
	MessageID snowflake.ID
}

// LoopMode represents the repeat mode of the player.
type LoopMode int

const (
	LoopOff LoopMode = iota
	LoopTrack
	LoopQueue
)

func (l LoopMode) String() string {
	switch l {
	case LoopTrack:
		return "リピート: トラック"
	case LoopQueue:
		return "リピート: キュー"
	default:
		return "リピートオフ"
	}
}

func (l LoopMode) Next() LoopMode {
	return (l + 1) % 3
}

// Queue is a thread-safe playlist that manages an ordered list of tracks
// with a current position and loop mode. All methods are safe for concurrent use.
type Queue struct {
	tracks  []lavalink.Track
	current int
	loop    LoopMode
	mu      sync.Mutex
}

func NewQueue() *Queue {
	return &Queue{}
}

func (q *Queue) Add(tracks ...lavalink.Track) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tracks = append(q.tracks, tracks...)
}

func (q *Queue) Next() (lavalink.Track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tracks) == 0 {
		return lavalink.Track{}, false
	}

	switch q.loop {
	case LoopTrack:
		return q.tracks[q.current], true
	case LoopQueue:
		q.current = (q.current + 1) % len(q.tracks)
		return q.tracks[q.current], true
	default:
		q.current++
		if q.current >= len(q.tracks) {
			return lavalink.Track{}, false
		}
		return q.tracks[q.current], true
	}
}

func (q *Queue) Previous() (lavalink.Track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tracks) == 0 || q.current <= 0 {
		return lavalink.Track{}, false
	}

	q.current--
	return q.tracks[q.current], true
}

func (q *Queue) Current() (lavalink.Track, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.tracks) == 0 || q.current < 0 || q.current >= len(q.tracks) {
		return lavalink.Track{}, false
	}
	return q.tracks[q.current], true
}

func (q *Queue) SetCurrent(index int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.current = index
}

func (q *Queue) Tracks() []lavalink.Track {
	q.mu.Lock()
	defer q.mu.Unlock()
	result := make([]lavalink.Track, len(q.tracks))
	copy(result, q.tracks)
	return result
}

func (q *Queue) CurrentIndex() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.current
}

func (q *Queue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.tracks)
}

func (q *Queue) Clear() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.tracks = nil
	q.current = 0
}

func (q *Queue) LoopMode() LoopMode {
	q.mu.Lock()
	defer q.mu.Unlock()
	return q.loop
}

func (q *Queue) SetLoopMode(mode LoopMode) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.loop = mode
}

func (q *Queue) CycleLoop() LoopMode {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.loop = q.loop.Next()
	return q.loop
}

// Shuffle randomizes the order of tracks in the queue while keeping the
// currently playing track at position 0.
func (q *Queue) Shuffle() {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.tracks) <= 1 {
		return
	}

	currentTrack := q.tracks[q.current]
	rest := make([]lavalink.Track, 0, len(q.tracks)-1)
	for i, t := range q.tracks {
		if i != q.current {
			rest = append(rest, t)
		}
	}

	for i := len(rest) - 1; i > 0; i-- {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(i+1)))
		j := int(n.Int64())
		rest[i], rest[j] = rest[j], rest[i]
	}

	q.tracks = append([]lavalink.Track{currentTrack}, rest...)
	q.current = 0
}

// QueueManager manages per-guild queues.
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
