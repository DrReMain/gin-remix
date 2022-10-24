package repository

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	gonanoid "github.com/matoous/go-nanoid/v2"
	"go-remix/model"
	"go-remix/model/apperrors"
	"log"
	"time"
)

type redisRepository struct {
	rds *redis.Client
}

func NewRedisRepository(rds *redis.Client) model.RedisRepository {
	return &redisRepository{
		rds: rds,
	}
}

const (
	ForgotPasswordPrefix = "forgot-password"
)

func (r *redisRepository) SetResetToken(ctx context.Context, id string) (string, error) {
	uid, err := gonanoid.New()
	if err != nil {
		log.Printf("Failed to generate id: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	if err = r.rds.Set(ctx, fmt.Sprintf("%s:%s", ForgotPasswordPrefix, uid), id, 24*time.Hour).Err(); err != nil {
		log.Printf("Failed to set link in redis: %v\n", err.Error())
		return "", apperrors.NewInternal()
	}

	return uid, nil
}

func (r *redisRepository) GetIdFromToken(ctx context.Context, token string) (string, error) {
	key := fmt.Sprintf("%s:%s", ForgotPasswordPrefix, token)
	val, err := r.rds.Get(ctx, key).Result()

	if err == redis.Nil {
		return "", apperrors.NewBadRequest(apperrors.InvalidResetToken)
	}
	if err != nil {
		log.Printf("Failed to get value from redis: %v\n", err)
		return "", apperrors.NewInternal()
	}

	r.rds.Del(ctx, key)

	return val, nil
}
