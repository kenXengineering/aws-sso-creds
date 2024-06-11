package get

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kenxengineering/aws-sso-creds/pkg/cache"
	"github.com/spf13/cobra"

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

			exportJSON, _ := cmd.Flags().GetBool("json")

			creds, err := cache.GetCredentials(cmd)
			if err != nil {
				return err
			}

			if exportJSON {
				output, err := json.Marshal(creds)
				if err != nil {
					return err
				}
				fmt.Println(string(output))
			} else {

				fmt.Println(aurora.Sprintf("Your temporary credentials for account %s are:", aurora.White(creds.AccountID)))
				fmt.Println("")

				fmt.Fprintln(os.Stdout, "AWS_ACCESS_KEY_ID\t", creds.AwsAccessKeyID)
				fmt.Fprintln(os.Stdout, "AWS_SECRET_ACCESS_KEY\t", creds.AwsSecretAccessKey)
				fmt.Fprintln(os.Stdout, "AWS_SESSION_TOKEN\t", creds.SessionToken)

				fmt.Println("")

				fmt.Println("These credentials will expire at:", aurora.Red(creds.ExpireAt))
			}

			return nil
		},
	}
	command.PersistentFlags().BoolP("json", "j", false, "print output in json format")
	return command
}
