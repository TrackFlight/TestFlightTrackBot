package types

type Config struct {
	LimitFree    int64 `json:"limit_free"`
	LimitPremium int64 `json:"limit_premium"`

	MaxFollowingLinks int64 `json:"max_following_links"`
}
