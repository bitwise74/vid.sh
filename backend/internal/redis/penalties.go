package redis

import (
	"context"
	"fmt"
	"strconv"
	"time"
)

const maxPenalties = 3

// CheckPenalties checks if a user has exceeded the maximum number of penalties.
func CheckPenalties(c context.Context, ip string) (*time.Time, error) {
	val := Rdb.Get(c, ip+":penalty").Val()

	valInt, err := strconv.Atoi(val)
	if err != nil {
		return nil, fmt.Errorf("penalty value conversion error: %w", err)
	}

	if valInt >= maxPenalties {
		expiration, err := Rdb.TTL(c, ip+":penalty").Result()
		if err != nil {
			return nil, fmt.Errorf("penalty TTL retrieval error: %w", err)
		}
		expirationTime := time.Now().Add(expiration)
		return &expirationTime, nil
	}

	return nil, nil
}

// AddPenalty adds a penalty to a user for a specified duration.
func AddPenalty(c context.Context, ip string, duration time.Duration) error {
	val := Rdb.Get(c, ip+":penalty").Val()

	valInt, err := strconv.Atoi(val)
	if err != nil && val != "" {
		return fmt.Errorf("penalty value conversion error: %w", err)
	}

	if valInt >= maxPenalties {
		return nil
	}

	valInt++
	err = Rdb.Set(c, ip+":penalty", valInt, duration).Err()
	if err != nil {
		return fmt.Errorf("setting penalty error: %w", err)
	}

	return nil
}
