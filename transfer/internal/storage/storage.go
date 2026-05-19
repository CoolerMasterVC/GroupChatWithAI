package storage

import (
	"sync"
	"time"
)

type SegmentInfo struct {
	Received int
	Total    int
	Last     time.Time
	Username string
	Segments []string
}

type Storage struct {
	mu   sync.RWMutex
	data map[time.Time]SegmentInfo
}

func NewStorage() *Storage {
	return &Storage{
		data: make(map[time.Time]SegmentInfo),
	}
}

func (s *Storage) AddOrUpdate(sendTime time.Time, segmentNumber, totalSegments int, username, payload string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, exists := s.data[sendTime]
	if !exists {
		info = SegmentInfo{
			Received: 0,
			Total:    totalSegments,
			Last:     time.Now().UTC(),
			Username: username,
			Segments: make([]string, totalSegments),
		}
	}
	info.Segments[segmentNumber-1] = payload
	info.Received++
	info.Last = time.Now().UTC()
	s.data[sendTime] = info
}

func (s *Storage) GetAndDelete(sendTime time.Time) (SegmentInfo, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	info, ok := s.data[sendTime]
	if ok {
		delete(s.data, sendTime)
	}
	return info, ok
}

func (s *Storage) GetAllKeys() []time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	keys := make([]time.Time, 0, len(s.data))
	for k := range s.data {
		keys = append(keys, k)
	}
	return keys
}

func (s *Storage) Peek(sendTime time.Time) (SegmentInfo, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	info, ok := s.data[sendTime]
	return info, ok
}
