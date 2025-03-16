package flag

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/asyauqi15/go-flag/model"
	"github.com/redis/go-redis/v9"
	"strconv"
)

func IsActive(ctx context.Context, c *Client, name string) (bool, error) {
	if c.rdb == nil {
		return false, errors.New("redis client is not initialized")
	}

	val, err := c.rdb.Get(ctx, "flag:"+name).Result()
	if errors.Is(err, redis.Nil) {
		return false, nil
	} else if err != nil {
		return false, err
	}

	var flag model.Flag
	err = json.Unmarshal([]byte(val), &flag)
	if err != nil {
		return false, err
	}

	return flag.Active, nil
}

func GetValue[T any](ctx context.Context, c *Client, name string) (value T, err error) {
	if c.rdb == nil {
		return value, errors.New("redis client is not initialized")
	}

	v, err := c.rdb.Get(ctx, "flag:"+name).Result()
	if errors.Is(err, redis.Nil) {
		return value, errors.New("feature not found")
	} else if err != nil {
		return value, err
	}

	var flag model.Flag
	err = json.Unmarshal([]byte(v), &flag)
	if err != nil {
		return value, err
	}

	switch any(value).(type) {
	case int64:
		int64Val, convErr := strconv.ParseInt(flag.Value, 10, 64)
		if convErr != nil {
			return value, errors.New("failed to convert to int64")
		}
		return any(int64Val).(T), nil

	case float64:
		floatVal, convErr := strconv.ParseFloat(flag.Value, 64)
		if convErr != nil {
			return value, errors.New("failed to convert to float64")
		}
		return any(floatVal).(T), nil

	case bool:
		boolVal, convErr := strconv.ParseBool(flag.Value)
		if convErr != nil {
			return value, errors.New("failed to convert to bool")
		}
		return any(boolVal).(T), nil

	case string:
		return any(flag.Value).(T), nil

	default:
		return value, fmt.Errorf("unsupported type: %T", value)
	}
}

func GetStructValue[T any](ctx context.Context, c *Client, name string) (value T, err error) {
	if c.rdb == nil {
		return value, errors.New("redis client is not initialized")
	}

	v, err := c.rdb.Get(ctx, "flag:"+name).Result()
	if errors.Is(err, redis.Nil) {
		return value, errors.New("feature not found")
	} else if err != nil {
		return value, err
	}

	var flag model.Flag
	err = json.Unmarshal([]byte(v), &flag)
	if err != nil {
		return value, err
	}

	err = json.Unmarshal([]byte(flag.Value), &value)
	if err != nil {
		return value, errors.New("failed to unmarshal flag's value")
	}
	return value, nil
}
