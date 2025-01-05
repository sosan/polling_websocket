package redisclient

import (
	"fmt"
	"log"
	"polling_websocket/pkg/common"
	"polling_websocket/pkg/domain/models"
	"time"
)

const (
	PollingGlobalAll = "polling:all"
	EmptyValue       = "_"
)

type PollingService interface {
}

type PollingRepository struct {
	redisClient *RedisClient
}

func NewPollingRepository(redisClient *RedisClient) *PollingRepository {
	return &PollingRepository{redisClient: redisClient}
}

func (a *PollingRepository) GetPollingGlobalAll() string {
	return PollingGlobalAll
}

func (a *PollingRepository) ValidateActionGlobalUUID(uuid *string) (bool, error) {
	exists, err := a.redisClient.Hexists(PollingGlobalAll, *uuid)
	if err != nil {
		log.Printf("ERROR | Redis HExists error: %v", err)
		return true, err
	}
	return exists, err
}

func (a *PollingRepository) AcquireLock(key, value string, expiration time.Duration) (locked bool, err error) {
	for i := 1; i < models.MaxAttempts; i++ {
		locked, err = a.redisClient.AcquireLock(key, value, expiration)
		if err == nil {
			return locked, err
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to redis for key %s, attempt %d: %v. Retrying in %v", key, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false, fmt.Errorf("ERROR | Cannot create lock for key %s. More than 10 intents", key)
}

func (a *PollingRepository) RemoveLock(key string) bool {
	for i := 1; i < models.MaxAttempts; i++ {
		countRemoved, err := a.redisClient.RemoveLock(key)
		if countRemoved == 0 {
			log.Printf("WARNING | Key already removed, previuous process take more than 20 seconds")
		}
		if err == nil && countRemoved <= 1 {
			return true
		}

		waitTime := common.RandomDuration(models.MaxRangeSleepDuration, models.MinRangeSleepDuration, i)
		log.Printf("ERROR | Cannot connect to redis for key %s, attempt %d: %v. Retrying in %v", key, i, err, waitTime)
		time.Sleep(waitTime)
	}
	return false
}

func (a *PollingRepository) SetNX(hashKey, actionID string, expiration time.Duration) (bool, error) {
	inserted, err := a.redisClient.SetEx(hashKey, actionID, expiration)
	return inserted, err
}
