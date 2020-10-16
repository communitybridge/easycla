// Copyright The Linux Foundation and each contributor to CommunityBridge.
// SPDX-License-Identifier: MIT

package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	log "github.com/communitybridge/easycla/infra/ssm/logging"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	region       string
	stage        string
	parameterKey string
)

func init() {
	log.Info("Running init...")

	// cobra.OnInitialize(initConfig)
	viper.AutomaticEnv()
	defaults := map[string]interface{}{
		"REGION": "us-east-1",
		"STAGE":  "dev",
	}

	for key, value := range defaults {
		viper.SetDefault(key, value)
	}

	rootCmd.PersistentFlags().StringVarP(&region, "region", "r", "us-east-1", "the AWS region")
	rootCmd.PersistentFlags().StringVarP(&stage, "stage", "s", "dev", "the environment stage")
	rootCmd.PersistentFlags().StringVarP(&parameterKey, "key", "k", "", "the SSM parameter key to get")
	err := viper.BindPFlag("region", rootCmd.PersistentFlags().Lookup("region"))
	if err != nil {
		log.Warnf("unable to set 'region' in the configuration - error: %+v", err)
	}
	err = viper.BindPFlag("stage", rootCmd.PersistentFlags().Lookup("stage"))
	if err != nil {
		log.Warnf("unable to set 'stage' in the configuration - error: %+v", err)
	}
	err = viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))
	if err != nil {
		log.Warnf("unable to set 'parameterKey' in the configuration - error: %+v", err)
	}
}

// GetAWSSession returns an AWS session based on the region and credentials
func GetAWSSession(awsRegion string) (*session.Session, error) {
	log.Debugf("Creating a new AWS session for region: %s", awsRegion)
	awsSession := session.Must(session.NewSession(
		&aws.Config{
			Region:                        aws.String(awsRegion),
			CredentialsChainVerboseErrors: aws.Bool(true),
			MaxRetries:                    aws.Int(5),
		},
	))

	log.Debugf("Successfully created a new AWS session for region: %s...", awsRegion)

	return awsSession, nil
}

// getSSMString is a generic routine to fetch the specified key value
func getSSMString(ssmClient *ssm.SSM, key string) (string, error) {
	log.Debugf("Loading SSM parameter: %s", key)
	value, err := ssmClient.GetParameter(&ssm.GetParameterInput{
		Name:           aws.String(key),
		WithDecryption: aws.Bool(false),
	})
	if err != nil {
		log.Warnf("unable to read SSM parameter %s - error: %+v", key, err)
		return "", err
	}

	return strings.TrimSpace(*value.Parameter.Value), nil
}

func invoke(cmd *cobra.Command, args []string) {
	if viper.GetString("key") == "" {
		println("missing --key command line argument")
		println(cmd.UsageString())
		return
	}

	awsSession, err := GetAWSSession(viper.GetString("region"))
	if err != nil {
		log.Warnf("Unable to load AWS session - Error: %v", err)
		return
	}
	ssmClient := ssm.New(awsSession)

	log.Debugf("Loading parameter '%s' from region %s, stage: %s", viper.GetString("key"), viper.GetString("region"), viper.GetString("stage"))
	value, err := getSSMString(ssmClient, viper.GetString("key"))
	if err != nil {
		return
	}
	log.Debugf("Parameter '%s' = %s", viper.GetString("key"), value)
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "get-ssm-parameter",
	Short: "Gets a SSM Parameter",
	Long:  "Fetches the specified SSM parameter from the AWS SSM configuration",
	Run:   invoke,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
