package main

import (
	"context"
	"encoding/json"
	"strings"

	redis "github.com/redis/go-redis/v9"
)

type Redis[T any] struct {
	client *redis.Client
}

func NewRedis[T any](rdb *redis.Client) *Redis[T] {
	return &Redis[T]{client: rdb}
}

// Start returns the executed result and true if the idempotency key is already executed. Otherwise, returns empty T and false
func (r *Redis[T]) Start(ctx context.Context, idempotencyKey string) (T, bool, error) {
	var t T
	tr := r.client.HSetNX(ctx, "idempotency:"+idempotencyKey, "status", "started")
	if tr.Err() != nil {
		return t, false, tr.Err()
	}
	if tr.Val() {
		return t, false, nil
	}
	b, err := r.client.HGet(ctx, "idempotency:"+idempotencyKey, "value").Bytes()
	if err != nil {
		return t, false, err
	}
	if err := json.Unmarshal(b, &t); err != nil {
		return t, false, err
	}
	return t, true, nil
}

func (r *Redis[T]) Store(ctx context.Context, idempotencyKey string, value T) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	err = r.client.HSet(ctx, "idempotency:"+idempotencyKey, "value", b).Err()
	if err != nil && strings.Contains(err.Error(), "wrong number") {
		tr := r.client.HMSet(ctx, "idempotency:"+idempotencyKey, "value", b)
		return tr.Err()
	}
	return err
}
