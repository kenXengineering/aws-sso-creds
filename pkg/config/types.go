package config

import "time"

type SSOConfig struct {
	StartURL  string
	Region    string
	AccountID string
	RoleName  string
}

type SSOCacheConfig struct {
	StartURL    string `json:"startUrl"`
	Region      string `json:"region"`
	AccessToken string `json:"accessToken"`
	ExpiresAt   string `json:"expiresAt"`
}

type JSON struct {
	AwsAccessKeyID     string    `json:"aws_access_key_id"`
	AwsSecretAccessKey string    `json:"aws_secret_access_key"`
	SessionToken       string    `json:"aws_session_token"`
	ExpireAt           time.Time `json:"expire_at"`
	AccountID          string    `json:"accountId"`
}
