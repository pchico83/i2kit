package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pchico83/i2kit/cli/providers/aws/cf"
	"github.com/pchico83/i2kit/cli/schemas/environment"
)

//CreateKeypair creates the "okteto" keypair if it is not yet created
func CreateKeypair(e *environment.Environment, config *aws.Config) error {
	svc := ec2.New(session.New(), config)
	cki := &ec2.CreateKeyPairInput{
		KeyName: aws.String(fmt.Sprintf("i2kit-%s", e.Name)),
	}
	_, err := svc.CreateKeyPair(cki)
	if err != nil && !cf.CheckError(err, "InvalidKeyPair.Duplicate", "") {
		return err
	}
	return nil
}
