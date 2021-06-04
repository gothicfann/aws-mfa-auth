package mfa

import (
	"bufio"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/iam"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/spf13/cobra"
)

const mfaSuffix = "mfa"

type Account struct {
	Name         string `yaml:"name"`
	Region       string `yaml:"region"`
	AccessKey    string `yaml:"accessKey"`
	SecretKey    string `yaml:"secretKey"`
	Expiration   time.Time
	SessionToken string
	MFASerial    string
}

func (a *Account) CreateSession() *session.Session {
	s, err := session.NewSession(&aws.Config{
		Region: aws.String(a.Region),
		Credentials: credentials.NewStaticCredentials(
			a.AccessKey,
			a.SecretKey,
			a.SessionToken,
		),
	})
	cobra.CheckErr(err)
	return s
}

func (a *Account) GetCurrentUserName(s *session.Session) *string {
	svc := iam.New(s)
	var userInput *iam.GetUserInput
	userOutput, err := svc.GetUser(userInput)
	cobra.CheckErr(err)
	return userOutput.User.UserName
}

func (a *Account) GetMFADeviceSerial(s *session.Session) {
	svc := iam.New(s)
	mfaInput := iam.ListMFADevicesInput{UserName: a.GetCurrentUserName(s)}
	mfaOutput, err := svc.ListMFADevices(&mfaInput)
	cobra.CheckErr(err)
	if len(mfaOutput.MFADevices) > 0 {
		a.MFASerial = *mfaOutput.MFADevices[0].SerialNumber
	}
}

func (a *Account) GetTempCredentials(s *session.Session, mfaSerial *string, durationSeconds int64) {
	svc := sts.New(s)
	fmt.Printf("Enter one-time passcode for \"%s\" account: ", a.Name)
	in := bufio.NewScanner(os.Stdin)
	in.Scan()
	tokenCode := in.Text()
	sessionTokenInput := sts.GetSessionTokenInput{
		DurationSeconds: &durationSeconds,
		SerialNumber:    mfaSerial,
		TokenCode:       &tokenCode,
	}
	sessionTokenOutput, err := svc.GetSessionToken(&sessionTokenInput)
	cobra.CheckErr(err)
	a.AccessKey = *sessionTokenOutput.Credentials.AccessKeyId
	a.SecretKey = *sessionTokenOutput.Credentials.SecretAccessKey
	a.SessionToken = *sessionTokenOutput.Credentials.SessionToken
	a.Expiration = *sessionTokenOutput.Credentials.Expiration
}

// PrintDebug print Account in debug mode (struct)
func (a *Account) PrintDebug() {
	fmt.Println(a)
	fmt.Println()
}

// PrintEnv prints Account in non-MFAd env format
func (a *Account) PrintEnv() {
	fmt.Printf("# %s\n", a.Name)
	fmt.Printf("AWS_REGION=%s\n", a.Region)
	fmt.Printf("AWS_ACCESS_KEY_ID=%s\n", a.AccessKey)
	fmt.Printf("AWS_SECRET_ACCESS_KEY=%s\n", a.SecretKey)
	fmt.Println()
}

func (a *Account) SprintEnv() string {
	s := fmt.Sprintf("# %s\n", a.Name)
	s += fmt.Sprintf("AWS_REGION=%s\n", a.Region)
	s += fmt.Sprintf("AWS_ACCESS_KEY_ID=%s\n", a.AccessKey)
	s += fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s\n", a.SecretKey)
	s += fmt.Sprintln()
	return s
}

// PrintMFAdEnv prints Account in MFAd env format
func (a *Account) PrintMFAdEnv() {
	fmt.Printf("# %s \n", a.Name)
	fmt.Printf("AWS_REGION=%s\n", a.Region)
	fmt.Printf("AWS_ACCESS_KEY_ID=%s\n", a.AccessKey)
	fmt.Printf("AWS_SECRET_ACCESS_KEY=%s\n", a.SecretKey)
	fmt.Printf("AWS_SESSION_TOKEN=%s\n", a.SessionToken)
	fmt.Println()
}

func (a *Account) SprintMFAdEnv() string {
	s := fmt.Sprintf("# %s \n", a.Name)
	s += fmt.Sprintf("AWS_REGION=%s\n", a.Region)
	s += fmt.Sprintf("AWS_ACCESS_KEY_ID=%s\n", a.AccessKey)
	s += fmt.Sprintf("AWS_SECRET_ACCESS_KEY=%s\n", a.SecretKey)
	s += fmt.Sprintf("AWS_SESSION_TOKEN=%s\n", a.SessionToken)
	s += fmt.Sprintln()
	return s
}

// PrintAwsRegion prints Account in non-MFAd aws format for aws config file
func (a *Account) PrintAwsRegion() {
	fmt.Printf("[%s]\n", a.Name)
	fmt.Printf("region = %s\n", a.Region)
	fmt.Println()
}

func (a *Account) SprintAwsRegion() string {
	s := fmt.Sprintf("[%s]\n", a.Name)
	s += fmt.Sprintf("region = %s\n", a.Region)
	s += fmt.Sprintln()
	return s
}

// PrintAws prints Account in non-MFAd aws format for aws config file
func (a *Account) PrintAws() {
	fmt.Printf("[%s]\n", a.Name)
	fmt.Printf("aws_access_key_id = %s\n", a.AccessKey)
	fmt.Printf("aws_secret_access_key = %s\n", a.SecretKey)
	fmt.Println()
}

func (a *Account) SprintAws() string {
	s := fmt.Sprintf("[%s]\n", a.Name)
	s += fmt.Sprintf("aws_access_key_id = %s\n", a.AccessKey)
	s += fmt.Sprintf("aws_secret_access_key = %s\n", a.SecretKey)
	s += fmt.Sprintln()
	return s
}

// PrintMFAdAwsRegion prints Account in MFAd aws format for aws config file
func (a *Account) PrintMFAdAwsRegion() {
	fmt.Printf("[%s-%s]\n", a.Name, mfaSuffix)
	fmt.Printf("region = %s\n", a.Region)
	fmt.Println()
}

func (a *Account) SprintMFAdAwsRegion() string {
	s := fmt.Sprintf("[%s-%s]\n", a.Name, mfaSuffix)
	s += fmt.Sprintf("region = %s\n", a.Region)
	s += fmt.Sprintln()
	return s
}

// PrintMFAdAws prints Account in MFAd aws format for aws credentials file
func (a *Account) PrintMFAdAws() {
	fmt.Printf("[%s-%s]\n", a.Name, mfaSuffix)
	fmt.Printf("aws_access_key_id = %s\n", a.AccessKey)
	fmt.Printf("aws_secret_access_key = %s\n", a.SecretKey)
	fmt.Printf("aws_session_token = %s\n", a.SessionToken)
	fmt.Println()
}

func (a *Account) SprintMFAdAws() string {
	s := fmt.Sprintf("[%s-%s]\n", a.Name, mfaSuffix)
	s += fmt.Sprintf("aws_access_key_id = %s\n", a.AccessKey)
	s += fmt.Sprintf("aws_secret_access_key = %s\n", a.SecretKey)
	s += fmt.Sprintf("aws_session_token = %s\n", a.SessionToken)
	s += fmt.Sprintln()
	return s
}

// PrintMFAdAccount outputs Account information in MFAd formats
func (a *Account) PrintMFAdAccount(format string) {
	switch format {
	case "env":
		a.PrintMFAdEnv()
	case "aws":
		a.PrintMFAdAwsRegion()
		a.PrintMFAdAws()
	}
}

// PrintAccount outputs Account information in non-MFAd formats
func (a *Account) PrintAccount(format string) {
	switch format {
	case "env":
		a.PrintEnv()
	case "aws":
		a.PrintAwsRegion()
		a.PrintAws()
	}
}

// Print outputs Account information in pretty formats
// It automatically detects wether Account is already MFAd or not
// Use this until more specific printing is required
func (a *Account) Print(format string) {
	if a.SessionToken == "" {
		a.PrintAccount(format)
	} else {
		a.PrintMFAdAccount(format)
	}
}
