package middleware

import (
	"math"
	"net/http"
	"sync"
	"time"

	"github.com/TrackFlight/TestFlightTrackBot/internal/api/utils"
)

type FloodInfo struct {
	Hits          int
	LastHit       time.Time
	WaitUntil     time.Time
	LastPenalty   time.Duration
	PenaltyExpire time.Time
}

func shouldAllow(
	floodMu *sync.Mutex,
	floodMap map[int64]*FloodInfo,
	key int64,
	maxHits int,
	timeWindow,
	initialPenalty,
	maxPenalty,
	penaltyDecayWindow time.Duration,
) (bool, time.Duration) {
	now := time.Now()
	floodMu.Lock()
	defer floodMu.Unlock()

	info, exists := floodMap[key]
	if !exists {
		info = &FloodInfo{
			Hits:        1,
			LastHit:     now,
			LastPenalty: initialPenalty,
		}
		floodMap[key] = info
		return true, 0
	}

	if now.Before(info.WaitUntil) {
		return false, info.WaitUntil.Sub(now)
	}

	if now.Sub(info.LastHit) > timeWindow {
		info.Hits = 1
		info.LastHit = now
		return true, 0
	}

	info.Hits++
	info.LastHit = now

	if info.Hits > maxHits {
		if now.Before(info.PenaltyExpire) {
			info.LastPenalty *= 2
			if info.LastPenalty > maxPenalty {
				info.LastPenalty = maxPenalty
			}
		} else {
			info.LastPenalty = initialPenalty
		}

		info.WaitUntil = now.Add(info.LastPenalty)
		info.PenaltyExpire = now.Add(penaltyDecayWindow)
		info.Hits = 0
		return false, info.LastPenalty
	}

	return true, 0
}

func AntiFlood(maxHits int, timeWindow, initialPenalty, maxPenalty, penaltyDecayWindow time.Duration) func(next http.Handler) http.Handler {
	var floodMu sync.Mutex
	floodMap := make(map[int64]*FloodInfo)

	go func() {
		for {
			time.Sleep(1 * time.Minute)
			now := time.Now()

			floodMu.Lock()
			for k, info := range floodMap {
				if now.Sub(info.LastHit) > 10*time.Minute && now.After(info.WaitUntil) {
					delete(floodMap, k)
				}
			}
			floodMu.Unlock()
		}
	}()

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			allowed, wait := shouldAllow(
				&floodMu,
				floodMap,
				r.Context().Value(UserIDKey).(int64),
				maxHits,
				timeWindow,
				initialPenalty,
				maxPenalty,
				penaltyDecayWindow,
			)
			if !allowed {
				utils.JSONErrorFloodWait(w, int(math.Ceil(wait.Seconds())))
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
