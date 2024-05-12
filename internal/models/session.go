package models

import (
	"encoding/json"
	"time"
)

type Session struct {
	ID           string    `json:"id"`
	UserID       string    `json:"user_id"`
	RefreshToken string    `json:"refresh_token"`
	ExpiresAt    time.Time `json:"expires_at"`
}

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func MarshalSession(s Session) ([]byte, error) {
	return json.Marshal(s)
}

func UnmarshalSession(b []byte) (Session, error) {
	var s Session
	if err := json.Unmarshal(b, &s); err != nil {
		return Session{}, err
	}

	return s, nil
}
