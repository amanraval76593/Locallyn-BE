package user

import (
	"encoding/json"
	"locallyn-be/internal/cache"
	"locallyn-be/internal/common/utils"
	"log"
)

type UserCache struct {
	redis *cache.Redis
}

type verifyUserData struct {
	UserId string `json:"userId"`
	Email  string `json:"email"`
}

func NewUserCache(r *cache.Redis) *UserCache {
	return &UserCache{
		redis: r,
	}
}

func (c *UserCache) setVerifyUserCode(userId string, email string) error {

	data, _ := json.Marshal(verifyUserData{UserId: userId, Email: email})

	code, err := utils.GenerateNanoId()

	if err != nil {
		return err
	}

	log.Printf("code:%s", code)

	return c.redis.Set(verifyUserKey(code), string(data), verifyUserTTL)
}

func (c *UserCache) getVerifyUserData(code string) (*verifyUserData, error) {

	data, err := c.redis.Get(verifyUserKey(code))

	if err != nil {
		return nil, err
	}

	var verifyData verifyUserData

	if err := json.Unmarshal([]byte(data), &verifyData); err != nil {
		return nil, err
	}

	return &verifyData, nil
}
