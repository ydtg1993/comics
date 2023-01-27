package rd

import (
	"time"
)

func RPush(key string, value interface{}) error {
	return RdClient.RPush(RedisCtx, key, value).Err()
}

func LPop(key string) (string, error) {
	return RdClient.LPop(RedisCtx, key).Result()
}

func LLen(key string) (int64, error) {
	return RdClient.LLen(RedisCtx, key).Result()
}

func Set(key string, value string, expiration time.Duration) error {
	return RdClient.Set(RedisCtx, key, value, expiration).Err()
}

func Get(key string) string {
	return RdClient.Get(RedisCtx, key).Val()
}

func HSet(key string, values ...interface{}) error {
	return RdClient.HSet(RedisCtx, key, values).Err()
}

func HLen(key string) int64 {
	return RdClient.HLen(RedisCtx, key).Val()
}

func HGet(key string, field string) (string, error) {
	return RdClient.HGet(RedisCtx, key, field).Result()
}

func GetCashWithHash(key string, field string, fun func() (string, error)) (string, error) {
	result, err := HGet(key, field)
	if err != nil && err.Error() != "redis: nil" {
		return "", err
	}
	if result == "" {
		result, err := fun()
		if err != nil {
			return "", err
		}

		if result == "" {
			return "", nil
		}

		if err = HSet(key, field, result); err != nil {
			return "", err
		}
		return result, nil
	}
	return result, nil
}

func GetCash(key string, fun func() string, expiration time.Duration) (result string, err error) {
	result = Get(key)
	if result == "" {
		result = fun()
		if err = Set(key, result, expiration); err != nil {
			return
		}
	}
	return
}

func Delete(key string) error {
	return RdClient.Del(RedisCtx, key).Err()
}

const Lock = "system_lock"

/**
锁定系统
*/
func LockSystem() {
	err := Set(Lock, "1", time.Duration(30)*time.Second)
	if err != nil {
		panic(err.Error())
	}
}

/**
检查锁定
*/
func CheckLock() bool {
	lock := Get(Lock)
	if lock == "1" {
		return true
	} else {
		return false
	}
}

func UnLockSystem() {
	err := Delete(Lock)
	if err != nil {
		panic(err)
	}
	time.Sleep(time.Second * 2)
}
