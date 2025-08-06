package types

type Config struct {
	LimitFree    int64 `json:"limit_free"`
	LimitPremium int64 `json:"limit_premium"`

	MaxFollowingPerUser int64 `json:"max_following_per_user"`
}
