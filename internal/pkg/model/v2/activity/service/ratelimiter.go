package service

import (
	"context"
	"errors"

	"github.com/ulule/limiter/drivers/store/memory"

	"github.com/ulule/limiter"
)

var limiterInstance *limiter.Limiter

// RateLimiter is a rate limiter service
// Limit can be specified in the format "<limit>-<period>"
//
// Valid periods:
// * "S": second
// * "M": minute
// * "H": hour
//
// Examples:
// * 5 requests / second : "5-S"
// * 5 requests / minute : "5-M"
// * 5 requests / hour : "5-H"
type RateLimiter struct {
	Limit string `json:"limit"`
	Token string `json:"token"`

	QuotaExceeded  bool  `json:"quotaExceeded"`
	AvailableQuota int64 `json:"availableQuota"`

	limiterInstance *limiter.Limiter
}

// Execute invokes this service
func (rl *RateLimiter) Execute() (err error) {

	//check for request token
	if rl.Token == "" {
		//TODO: Need to handle 'token not found' case elegantly. Time being set to dummy value.
		rl.Token = "DUMMYTOKEN"
	}

	//consume quota
	limiterContext, err := rl.limiterInstance.Get(context.TODO(), rl.Token)
	if err != nil {
		return nil
	}

	//check the ratelimit
	rl.AvailableQuota = limiterContext.Remaining
	if limiterContext.Reached {
		rl.QuotaExceeded = true
	} else {
		rl.QuotaExceeded = false
	}

	return nil
}

// UpdateRequest updates a request on an existing RateLimiter service instance with new values.
func (rl *RateLimiter) UpdateRequest(values map[string]interface{}) (err error) {
	err = rl.setRequestValues(values)
	return err
}

// InitializeRateLimiter initializes a RateLimiter services with provided settings.
func InitializeRateLimiter(settings map[string]interface{}) (rl *RateLimiter, err error) {
	rl = &RateLimiter{}
	err = rl.setRequestValues(settings)

	//create limiter instance once
	if limiterInstance == nil {
		//create rate
		rate, err := limiter.NewRateFromFormatted(rl.Limit)
		if err != nil {
			panic(err)
		}
		//create store
		store := memory.NewStore()
		//create limiter
		limiterInstance = limiter.New(store, rate)
		rl.limiterInstance = limiterInstance
	} else {
		rl.limiterInstance = limiterInstance
	}

	return
}

func (rl *RateLimiter) setRequestValues(settings map[string]interface{}) (err error) {
	for k, v := range settings {
		switch k {
		case "limit":
			limit, ok := v.(string)
			if !ok {
				return errors.New("Invalid type for limit")
			}
			rl.Limit = limit
		case "token":
			token, ok := v.(string)
			if !ok {
				rl.Token = ""
			}
			rl.Token = token
		default:
			//ignore the seting
		}
	}
	return nil
}
