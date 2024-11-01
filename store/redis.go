package store

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

var _ Storage = &Redis{}

type Redis struct {
	rdb *redis.Client
}

func NewRedis(url, password string, db int) (*Redis, error) {
	return &Redis{
		rdb: redis.NewClient(&redis.Options{
			Addr:     url,
			Password: password,
			DB:       db,
		}),
	}, nil
}

func (rd *Redis) Close() error {
	return rd.rdb.Close()
}

func (rd *Redis) Create(objectID ID, data Blob) error {
	_, err := rd.rdb.Get(context.Background(), string(objectID)).Result()
	if err != nil && err != redis.Nil {
		return err
	}
	if err != redis.Nil { // The item is already there
		return fmt.Errorf("item with id %v already exists", objectID)
	}

	err = rd.rdb.Set(context.Background(), string(objectID), data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rd *Redis) LoadAll() ([]Item, error) {
	ctx := context.Background()
	iter := rd.rdb.Scan(ctx, 0, "", 0).Iterator()
	res := []Item{}
	for iter.Next(ctx) {
		key := iter.Val()
		val, _ := rd.rdb.Get(ctx, key).Result()
		res = append(res, Item{ID: ID(key), Blob: Blob(val)})
	}
	return res, nil
}

func (rd *Redis) Load(objectID ID) (Blob, error) {
	data, err := rd.rdb.Get(context.Background(), string(objectID)).Result()
	if err == redis.Nil {
		return nil, fmt.Errorf("id %s does not exist", objectID)
	}
	if err != nil {
		return nil, err
	}
	return Blob(data), nil
}

func (rd *Redis) Save(objectID ID, blob Blob) error {
	// Non thread safe!
	data, err := rd.rdb.Get(context.Background(), string(objectID)).Result()
	if err == redis.Nil {
		return fmt.Errorf("id %s does not exist", objectID)
	}

	err = rd.rdb.Set(context.Background(), string(objectID), data, 0).Err()
	if err != nil {
		return err
	}
	return nil
}

func (rd *Redis) Delete(objectID ID) error {
	// Non thread safe!
	data, err := rd.rdb.Del(context.Background(), string(objectID)).Result()
	if err == redis.Nil {
		return fmt.Errorf("id %s does not exist", objectID)
	}

	err = rd.rdb.Set(context.Background(), string(objectID), data, 0).Err()
	if err != nil {
		return err
	}

	return nil
}
