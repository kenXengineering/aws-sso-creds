package main

import (
	"fmt"
	"github.com/spf13/viper"
	"os"

	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/get"
	"github.com/jaxxstorm/aws-sso-creds/cmd/aws-sso-creds/set"
	"github.com/jaxxstorm/aws-sso-creds/pkg/contract"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

var (
	profile string
)

func configureCLI() *cobra.Command {
	rootCommand := &cobra.Command{
		Use:  "aws-sso-creds",
		Long: "A helper utility to interact with AWS SSO",
	}

	rootCommand.AddCommand(get.Command())
	rootCommand.AddCommand(set.Command())

	homeDir, err := homedir.Dir()

	if err != nil {
		panic("Cannot find home directory, fatal error")
	}

	rootCommand.PersistentFlags().StringVarP(&profile, "profile", "p", "", "the AWS profile to use")
	rootCommand.PersistentFlags().StringVarP(&homeDir, "home-directory", "H", homeDir, "specify a path to a home directory")
	viper.BindEnv("profile", "AWS_PROFILE")
	viper.BindPFlag("profile", rootCommand.PersistentFlags().Lookup("profile"))
	viper.BindPFlag("home-directory", rootCommand.PersistentFlags().Lookup("home-directory"))

	return rootCommand
}

func main() {
	rootCommand := configureCLI()

	if err := rootCommand.Execute(); err != nil {
		contract.IgnoreIoError(fmt.Fprintf(os.Stderr, "%s", err))
		os.Exit(1)
	}
}