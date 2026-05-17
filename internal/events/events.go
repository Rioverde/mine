package events

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

const (
	streamKey = "events"
	groupName = "merchant"
	prefix    = "merchant-"
)

type Event struct {
	ID     string
	Values map[string]any
}

type Bus interface {
	Publish(ctx context.Context, values map[string]any) (string, error)
	Subscribe(ctx context.Context, buffer int) <-chan Event
	Close() error
}

type redisBus struct {
	rdb    *redis.Client
	maxLen int64
}

func Connect(ctx context.Context, maxLen int) (Bus, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", env("CACHE_HOST", "localhost"), env("CACHE_PORT", "6379")),
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, err
	}

	err := rdb.XGroupCreateMkStream(ctx, streamKey, groupName, "$").Err()
	if err != nil && !strings.Contains(err.Error(), "BUSYGROUP") {
		_ = rdb.Close()
		return nil, err
	}
	return &redisBus{rdb: rdb, maxLen: int64(maxLen)}, nil
}

func (b *redisBus) Publish(ctx context.Context, values map[string]any) (string, error) {
	return b.rdb.XAdd(ctx, &redis.XAddArgs{
		Stream: streamKey,
		MaxLen: b.maxLen,
		Approx: true,
		Values: values,
	}).Result()
}

func (b *redisBus) Subscribe(ctx context.Context, buffer int) <-chan Event {
	out := make(chan Event, buffer)
	hostname, _ := os.Hostname()
	consumer := prefix + hostname
	go func() {
		defer close(out)
		for {

			if ctx.Err() != nil {
				return
			}
			streams, err := b.rdb.XReadGroup(ctx, &redis.XReadGroupArgs{
				Group:    groupName,
				Consumer: consumer,
				Streams:  []string{streamKey, ">"},
				Block:    5 * time.Second,
				Count:    10,
			}).Result()

			if err != nil {
				if ctx.Err() != nil {
					return
				}
				if errors.Is(err, redis.Nil) {
					continue
				}
				slog.Error("events: XRead failed", "err", err)
				time.Sleep(time.Second)
				continue
			}

			for _, stream := range streams {
				for _, msg := range stream.Messages {
					event := Event{ID: msg.ID, Values: msg.Values}
					select {
					case <-ctx.Done():
						return
					case out <- event:
						//TODO Do this inside of the other service
						// err := b.rdb.XAck(ctx, streamKey, groupName, msg.ID).Err()
						// if err != nil {
						// 	slog.Error("events: XAck failed", "err", err, "id", msg.ID)
						// }
					}
				}
			}
		}
	}()
	return out
}

func (b *redisBus) Close() error {
	return b.rdb.Close()
}

func env(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
