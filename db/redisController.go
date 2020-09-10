package db

import "github.com/go-redis/redis"

//RedisController ...
type RedisController struct {
	redisClient *redis.ClusterClient
}

//key에 여러개의 field와 value를 저장할 수 있다. 기존에 같은 field가 있으면 덮어쓴다
func (s *RedisController) Set(key string, field string, value interface{}) {
	err := s.redisClient.HSet(key, field, value).Err()
	if err != nil {
		panic(err)
	}
}

func (s *RedisController) HWWKDDSet(key string, field string, value interface{}) {
	err := s.redisClient.HSet(key, field, value).Err()
	if err != nil {
		panic(err)
	}
}

func (s *RedisController) HMSet(key string, fields interface{}) {
	err := s.redisClient.HMSet(key, fields.(map[string]interface{})).Err()
	if err != nil {
		panic(err)
	}
}

func (s *RedisController) Get(key string, field string) string {
	data, err := s.redisClient.HGet(key, field).Result()
	// Execute cached commands in for loop
	if err == redis.Nil {
		return ""
	} else if err != nil {
		panic(err)
	}
	return data
}

func (s *RedisController) HIncrBy(key string, field string, incr int64) int64 {
	hIncrBy := s.redisClient.HIncrBy(key, field, incr)
	err := hIncrBy.Err()
	data := hIncrBy.Val()
	if err != nil {
		panic(err)
	}
	// Execute cached commands in for loop
	return data
}

//key에 속한 모든 field와 value를 조회
func (s *RedisController) GetAll(key string) map[string]string {
	data, err := s.redisClient.HGetAll(key).Result()
	// Execute cached commands in for loop
	if err == redis.Nil {
		return nil
	} else if err != nil {
		panic(err)
	}
	return data
}

func (s *RedisController) Delete(key string, fields ...string) {
	_, err := s.redisClient.HDel(key, fields...).Result()

	if err != nil {
		panic(err)
	}
}

func (s *RedisController) DeleteAll(key ...string) {
	_, err := s.redisClient.Del(key...).Result()

	if err != nil {
		panic(err)
	}
}

func NewRedisController(addrs []string) *RedisController {
	rc := new(RedisController)

	rc.redisClient = redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    addrs,
		Password: "", // no password set
		//DB:       0,  // use default DB
	})

	return rc
}
