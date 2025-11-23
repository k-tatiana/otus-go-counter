package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"go-server-counters/models"

	"github.com/redis/go-redis/v9"
)

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string) error
}

type MessageCounter struct {
	db RedisClient
}

func NewMessageCounter(db RedisClient) *MessageCounter {
	return &MessageCounter{
		db: db,
	}
}

func (mc *MessageCounter) IncrementCounter(ctx context.Context, key string) error {
	val, err := mc.db.Get(ctx, key)
	if err != nil {
		if errors.Is(redis.Nil, err) {
			val = ""
		} else {
			return fmt.Errorf("get data from db: %w", err)
		}
	}

	var cnt models.MessageCounter
	if val != "" {
		err := json.Unmarshal([]byte(val), &cnt)
		if err != nil {
			return fmt.Errorf("unmarshal data: %w", err)
		}
	}

	cnt.CountTotal++
	cnt.CountUnread++
	data, err := json.Marshal(cnt)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}
	err = mc.db.Set(ctx, key, string(data))
	if err != nil {
		return fmt.Errorf("put data to db: %w", err)
	}

	return nil
}

func (mc *MessageCounter) DecrementCounter(ctx context.Context, key string) error {
	val, err := mc.db.Get(ctx, key)
	if err != nil {
		if errors.Is(redis.Nil, err) {
			val = ""
		} else {
			return fmt.Errorf("get data from db: %w", err)
		}
	}

	var cnt models.MessageCounter
	if val != "" {
		err := json.Unmarshal([]byte(val), &cnt)
		if err != nil {
			return fmt.Errorf("unmarshal data: %w", err)
		}
	}

	if cnt.CountTotal > 0 {
		cnt.CountTotal--
	}
	if cnt.CountUnread > 0 {
		cnt.CountUnread--
	}

	data, err := json.Marshal(cnt)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}
	err = mc.db.Set(ctx, key, string(data))
	if err != nil {
		return fmt.Errorf("put data to db: %w", err)
	}

	return nil
}

func (mc *MessageCounter) ReadMessages(ctx context.Context, key string) (int, error) {
	var (
		readCount int
		cnt       models.MessageCounter
	)
	val, err := mc.db.Get(ctx, key)
	if err != nil {
		return 0, fmt.Errorf("get data from db: %w", err)
	}

	if val != "" {
		err := json.Unmarshal([]byte(val), &cnt)
		if err != nil {
			return 0, fmt.Errorf("unmarshal data: %w", err)
		}
	}

	readCount = cnt.CountUnread
	cnt.CountUnread = 0

	data, err := json.Marshal(cnt)
	if err != nil {
		return 0, fmt.Errorf("marshal data: %w", err)
	}
	err = mc.db.Set(ctx, key, string(data))
	if err != nil {
		return 0, fmt.Errorf("put data to db: %w", err)
	}

	return readCount, nil
}

func (mc *MessageCounter) UnreadMessages(ctx context.Context, key string, count int) error {
	var cnt models.MessageCounter

	val, err := mc.db.Get(ctx, key)
	if err != nil {
		return fmt.Errorf("get data from db: %w", err)
	}

	if val != "" {
		err := json.Unmarshal([]byte(val), &cnt)
		if err != nil {
			return fmt.Errorf("unmarshal data: %w", err)
		}
	}
	cnt.CountUnread = count

	data, err := json.Marshal(cnt)
	if err != nil {
		return fmt.Errorf("marshal data: %w", err)
	}
	err = mc.db.Set(ctx, key, string(data))
	if err != nil {
		return fmt.Errorf("put data to db: %w", err)
	}

	return nil
}

func (mc *MessageCounter) GetMessageCounter(ctx context.Context, key string) (*models.MessageCounter, error) {
	val, err := mc.db.Get(ctx, key)
	if err != nil {
		if errors.Is(redis.Nil, err) {
			return &models.MessageCounter{}, nil
		}

		return nil, fmt.Errorf("get data from db: %w", err)
	}

	var cnt models.MessageCounter
	if val != "" {
		err := json.Unmarshal([]byte(val), &cnt)
		if err != nil {
			return nil, fmt.Errorf("unmarshal data: %w", err)
		}
	}

	return &cnt, nil
}
