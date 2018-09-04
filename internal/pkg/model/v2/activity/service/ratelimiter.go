package service

import (
	"context"
	"errors"
	"sync"

	"github.com/ulule/limiter/drivers/store/memory"

	"github.com/ulule/limiter"
)

// NewLimiter creates new limiter with specified limit
func NewLimiter(limit string) *limiter.Limiter {
	//create rate
	rate, err := limiter.NewRateFromFormatted(limit)
	if err != nil {
		panic(err)
	}
	//create store
	store := memory.NewStore()
	//create limiter
	limiter := limiter.New(store, rate)
	return limiter
}

// Limiters is set of rate limiters
type Limiters struct {
	limiters map[string]*limiter.Limiter
	sync.RWMutex
}

// Lookup looks up limiter with specified name
func (l *Limiters) Lookup(name string, limit string) *limiter.Limiter {
	l.RLock()
	limiter := l.limiters[name]
	l.RUnlock()

	if limiter != nil {
		return limiter
	}

	limiter = NewLimiter(limit)
	l.Lock()
	l.limiters[name] = limiter
	l.Unlock()

	return limiter
}

var limiters = Limiters{
	limiters: make(map[string]*limiter.Limiter),
}

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
	serviceName string

	Limit string `json:"limit"`
	Token string `json:"token"`

	QuotaExceeded  bool  `json:"quotaExceeded"`
	AvailableQuota int64 `json:"availableQuota"`
}

// Execute invokes this service
func (rl *RateLimiter) Execute() (err error) {

	//check for request token
	if rl.Token == "" {
		//TODO: Need to handle 'token not found' case elegantly. Time being set to dummy value.
		rl.Token = "DUMMYTOKEN"
	}

	limiter := limiters.Lookup("NAME", rl.Limit)

	//consume quota
	limiterContext, err := limiter.Get(context.TODO(), rl.Token)
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
func InitializeRateLimiter(name string, settings map[string]interface{}) (rl *RateLimiter, err error) {
	rl = &RateLimiter{
		serviceName: name,
	}
	err = rl.setRequestValues(settings)
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
