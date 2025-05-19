package utils

import (
	"errors"
	"sync"
	"time"
)

const (
	workerIDBits   = 10
	sequenceBits   = 12
	workerIDShift  = sequenceBits
	timestampShift = sequenceBits + workerIDBits
	sequenceMask   = -1 ^ (-1 << sequenceBits)
	epoch          = 1609459200000
)

// Snowflake 雪花算法结构体
type Snowflake struct {
	mu        sync.Mutex
	lastStamp int64
	workerID  int64
	sequence  int64
}

// NewSnowflake 创建新的雪花算法实例
func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID >= (1<<workerIDBits) {
		return nil, errors.New("worker ID must be between 0 and 1023")
	}
	return &Snowflake{
		lastStamp: 0,
		workerID:  workerID,
		sequence:  0,
	}, nil
}

// NextID 生成下一个ID
func (s *Snowflake) NextID() (int64, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now().UnixNano() / 1e6

	if now < s.lastStamp {
		return 0, errors.New("clock moved backwards")
	}

	if now == s.lastStamp {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.lastStamp {
				now = time.Now().UnixNano() / 1e6
			}
		}
	} else {
		s.sequence = 0
	}

	s.lastStamp = now

	id := ((now - epoch) << timestampShift) |
		(s.workerID << workerIDShift) |
		s.sequence

	return id, nil
}
