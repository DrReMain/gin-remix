package main

import (
	"context"
	"fmt"
	"go-remix/model"
	"log"

	"github.com/go-redis/redis/v8"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-remix/config"
)

type dataSources struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func initDS(ctx context.Context, cfg config.Config) (*dataSources, error) {
	log.Println("初始化数据源")

	log.Println("连接 Postgresql ...")
	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl))
	if err != nil {
		return nil, fmt.Errorf("连接数据库错误: %w", err)
	}

	if err = db.AutoMigrate(
		&model.User{},
	); err != nil {
		return nil, fmt.Errorf("同步数据库模型错误: %w", err)
	}

	opt, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, fmt.Errorf("解析 Redis 地址错误: %w", err)
	}

	log.Println("连接 Redis ...")
	rdb := redis.NewClient(opt)
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("连接 redis 错误: %w", err)
	}

	return &dataSources{
		DB:          db,
		RedisClient: rdb,
	}, nil

}

func (d *dataSources) close() error {
	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("关闭Redis连接错误: %w", err)
	}

	return nil
}
