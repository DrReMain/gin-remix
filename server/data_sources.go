package main

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-remix/config"
	"go-remix/model"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

type dataSources struct {
	DB          *gorm.DB
	RedisClient *redis.Client
}

func initDS(ctx context.Context, cfg config.Config) (*dataSources, error) {
	log.Printf("Initializing data sources\n")

	log.Printf("Connecting to Postgresql\n")
	db, err := gorm.Open(postgres.Open(cfg.DatabaseUrl))
	if err != nil {
		return nil, fmt.Errorf("error opening db: %w", err)
	}

	if err = db.AutoMigrate(
		&model.User{},
	); err != nil {
		return nil, fmt.Errorf("error migrating models: %w", err)
	}

	//if err = db.SetupJoinTable(); err != nil {
	//	return nil, fmt.Errorf("error creating join table: %w", err)
	//}

	// Redis
	opt, err := redis.ParseURL(cfg.RedisUrl)
	if err != nil {
		return nil, fmt.Errorf("error parsing the redis url: %w", err)
	}

	log.Println("Connecting to Redis")
	rdb := redis.NewClient(opt)
	_, err = rdb.Ping(ctx).Result()
	if err != nil {
		return nil, fmt.Errorf("error connecting to redis: %w", err)
	}

	return &dataSources{
		DB:          db,
		RedisClient: rdb,
	}, nil

}

func (d *dataSources) close() error {
	if err := d.RedisClient.Close(); err != nil {
		return fmt.Errorf("error closing Redis Client: %w", err)
	}

	return nil
}
