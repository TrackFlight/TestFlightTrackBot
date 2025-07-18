package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/url"
	"sort"
	"strings"
)

type user struct {
	ID int64 `json:"id"`
}

func ValidateInitData(initData, botToken string) (int64, bool) {
	values, err := url.ParseQuery(initData)
	if err != nil {
		return 0, false
	}

	hash := values.Get("hash")
	values.Del("hash")

	var keys []string
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var dataStrings []string
	var userID int64
	for _, k := range keys {
		value := values.Get(k)
		if k == "user" {
			var usr user
			if err = json.Unmarshal([]byte(value), &usr); err != nil {
				return 0, false
			}
			userID = usr.ID
		}
		dataStrings = append(dataStrings, k+"="+value)
	}
	dataCheckString := strings.Join(dataStrings, "\n")

	h := hmac.New(sha256.New, []byte("WebAppData"))
	h.Write([]byte(botToken))
	secret := h.Sum(nil)

	mac := hmac.New(sha256.New, secret)
	mac.Write([]byte(dataCheckString))
	expectedHash := hex.EncodeToString(mac.Sum(nil))

	return userID, expectedHash == hash && userID > 0
}
