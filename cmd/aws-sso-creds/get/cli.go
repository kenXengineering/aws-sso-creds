package get

import (
	"encoding/json"
	"fmt"
	"github.com/aws/aws-sdk-go/service/sso"
	"github.com/jaxxstorm/aws-sso-creds/pkg/config"
	"os"
	"time"

	"github.com/jaxxstorm/aws-sso-creds/pkg/credentials"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/logrusorgru/aurora"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:          "get",
		Short:        "Get AWS temporary credentials to use on the command line",
		Long:         "Retrieve AWS temporary credentials",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true

			profile := viper.GetString("profile")
			homeDir := viper.GetString("home-directory")
			exportJSON, _ := cmd.Flags().GetBool("json")
			exportEnv, _ := cmd.Flags().GetBool("env")
			useCache, _ := cmd.Flags().GetBool("cache")
			refreshCache, _ := cmd.Flags().GetBool("refresh-cache")

			var ssoCreds *sso.GetRoleCredentialsOutput
			var accountID string
			var err error
			var creds *config.JSON
			cacheGood := false

			if useCache && !refreshCache {
				creds, err = config.GetCache()
				if err != nil {
					return err
				}
				// Cache is good if we get an AWS Account ID and the expiration time is in the future
				cacheGood = creds.AccountID != "" && creds.ExpireAt.After(time.Now())
			}

			if !useCache || !cacheGood || refreshCache {
				ssoCreds, accountID, err = credentials.GetSSOCredentials(profile, homeDir)
				if err != nil {
					return err
				}
				creds = &config.JSON{
					AwsAccessKeyID:     *ssoCreds.RoleCredentials.AccessKeyId,
					AwsSecretAccessKey: *ssoCreds.RoleCredentials.SecretAccessKey,
					SessionToken:       *ssoCreds.RoleCredentials.SessionToken,
					ExpireAt:           time.UnixMilli(*ssoCreds.RoleCredentials.Expiration),
					AccountID:          accountID,
				}
			}

			if exportJSON {
				output, err := json.Marshal(creds)
				if err != nil {
					return err
				}
				fmt.Println(string(output))
			} else if exportEnv {
				fmt.Println("export AWS_ACCESS_KEY_ID=" + creds.AwsAccessKeyID)
				fmt.Println("export AWS_SECRET_ACCESS_KEY=" + creds.AwsSecretAccessKey)
				fmt.Println("export AWS_SESSION_TOKEN=" + creds.SessionToken)
			} else {

				fmt.Println(aurora.Sprintf("Your temporary credentials for account %s are:", aurora.White(accountID)))
				fmt.Println("")

				fmt.Fprintln(os.Stdout, "AWS_ACCESS_KEY_ID\t", creds.AwsAccessKeyID)
				fmt.Fprintln(os.Stdout, "AWS_SECRET_ACCESS_KEY\t", creds.AwsSecretAccessKey)
				fmt.Fprintln(os.Stdout, "AWS_SESSION_TOKEN\t", creds.SessionToken)

				fmt.Println("")

				fmt.Println("These credentials will expire at:", aurora.Red(creds.ExpireAt))
			}

			if (useCache && !cacheGood) || refreshCache {
				return config.SetCache(creds)
			}

			return nil
		},
	}
	command.PersistentFlags().BoolP("json", "j", false, "print output in json format")
	command.PersistentFlags().BoolP("env", "e", false, "print output as environment variables")
	command.PersistentFlags().Bool("cache", true, "cache credentials.  System will update credentials if they are expired")
	command.PersistentFlags().Bool("refresh-cache", false, "retrieves new credentials and updates the cache with them")
	return command
}
