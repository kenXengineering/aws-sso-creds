package cache

import (
	"encoding/json"
	"github.com/aws/aws-sdk-go-v2/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"time"
)

func GetCredentials(cmd *cobra.Command) (*config.JSON, error) {
	profile := viper.GetString("profile")
	homeDir := viper.GetString("home-directory")
	useCache, _ := cmd.Flags().GetBool("cache")
	refreshCache, _ := cmd.Flags().GetBool("refresh-cache")
	var ssoCreds *sso.GetRoleCredentialsOutput
	var accountID string
	var err error
	var creds *config.JSON
	cacheGood := false

	if useCache && !refreshCache {
		creds, err = GetCache(profile)
		if err != nil {
			return nil, err
		}
		// Cache is good if we get an AWS Account ID and the expiration time is in the future
		cacheGood = creds.AccountID != "" && creds.ExpireAt.After(time.Now())
	}

	if !useCache || !cacheGood || refreshCache {
		ssoCreds, accountID, err = credentials.GetSSOCredentials(profile, homeDir)
		if err != nil {
			return nil, err
		}
		creds = &config.JSON{
			AwsAccessKeyID:     *ssoCreds.RoleCredentials.AccessKeyId,
			AwsSecretAccessKey: *ssoCreds.RoleCredentials.SecretAccessKey,
			SessionToken:       *ssoCreds.RoleCredentials.SessionToken,
			ExpireAt:           time.UnixMilli(ssoCreds.RoleCredentials.Expiration),
			AccountID:          accountID,
		}
	}

	if (useCache && !cacheGood) || refreshCache {
		if err := SetCache(profile, creds); err != nil {
			return nil, err
		}
	}

	return creds, nil
}

func GetCachePath() string {
	homeDir, err := homedir.Dir()
	if err != nil {
		panic("Cannot find home directory, fatal error")
	}
	configPath := homeDir + "/.config"
	if configPath[:len(configPath)-1] != "/" {
		configPath += "/"
	}
	return configPath + "aws-sso-creds"
}

func GetCache(profile string) (*config.JSON, error) {
	if CacheFileExists(profile) {
		cacheFile, err := os.Open(GetCacheFilePath(profile))
		if err != nil {
			return &config.JSON{}, err
		}
		jsonParser := json.NewDecoder(cacheFile)
		data := &config.JSON{}
		if err := jsonParser.Decode(&data); err != nil {
			return &config.JSON{}, err
		}
		return data, nil
	}
	return &config.JSON{}, nil
}

func SetCache(profile string, data *config.JSON) error {
	if !CacheFileExists(profile) {
		if err := os.MkdirAll(GetCachePath(), 0700); err != nil {
			return err
		}
	}
	cacheFile, err := os.OpenFile(GetCacheFilePath(profile), os.O_CREATE|os.O_TRUNC|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	jsonEncoder := json.NewEncoder(cacheFile)
	if err := jsonEncoder.Encode(data); err != nil {
		return err
	}
	return nil
}

func CacheFileExists(profile string) bool {
	if _, err := os.Stat(GetCacheFilePath(profile)); err != nil {
		return false
	}
	return true
}

func GetCacheFilePath(profile string) string {
	return GetCachePath() + "/" + profile + "_cache.json"
}
