package exportps

import (
	"fmt"
	"github.com/jaxxstorm/aws-sso-creds/pkg/cache"

	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	command := &cobra.Command{
		Use:          "export-ps",
		Short:        "Generates a set of powershell environment assignments to define the AWS temporary creds to your environment",
		Long:         "Generates a set of powershell environment assignments to define the AWS temporary creds to your environment",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {

			cmd.SilenceUsage = true

			creds, err := cache.GetCredentials(cmd)
			if err != nil {
				return err
			}

			fmt.Printf("$env:AWS_ACCESS_KEY_ID='%s'\n", creds.AwsAccessKeyID)
			fmt.Printf("$env:AWS_SECRET_ACCESS_KEY='%s'\n", creds.AwsSecretAccessKey)
			fmt.Printf("$env:AWS_SESSION_TOKEN='%s'\n", creds.SessionToken)

			return nil
		},
	}

	return command
}
