package ec2

import (
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/pchico83/i2kit/cli/schemas/environment"
)

//CreateSG creates the project security group
func CreateSG(e *environment.Environment, config *aws.Config) error {
	vpc, err := GetVPC(e, config)
	if err != nil {
		return err
	}
	svc := ec2.New(session.New(), config)
	description := fmt.Sprintf("Security Group for environment %s", e.Name)
	name := fmt.Sprintf("%s-i2kit", e.Name)
	csgi := &ec2.CreateSecurityGroupInput{
		Description: &description,
		GroupName:   &name,
		VpcId:       &vpc,
	}
	sg, err := svc.CreateSecurityGroup(csgi)
	if err != nil {
		if !strings.Contains(err.Error(), "InvalidGroup.Duplicate") {
			return err
		}
		dsgi := &ec2.DescribeSecurityGroupsInput{
			GroupNames: []*string{&name},
		}
		sgs, err2 := svc.DescribeSecurityGroups(dsgi)
		if err2 != nil {
			return err2
		}
		for _, i := range sgs.SecurityGroups {
			if *i.VpcId == vpc {
				e.Provider.SecurityGroup = *i.GroupId
			}
		}
	} else {
		e.Provider.SecurityGroup = *sg.GroupId
	}

	if e.Provider.SecurityGroup == "" {
		return fmt.Errorf("Error retrieving SG '%s'", name)
	}

	asgii := &ec2.AuthorizeSecurityGroupIngressInput{
		GroupName:               &name,
		SourceSecurityGroupName: &name,
	}
	protocol := "-1"
	if _, err = svc.AuthorizeSecurityGroupIngress(asgii); err != nil {
		if strings.Contains(err.Error(), "InvalidParameterValue") {
			asgii = &ec2.AuthorizeSecurityGroupIngressInput{
				GroupId: &e.Provider.SecurityGroup,
				IpPermissions: []*ec2.IpPermission{
					&ec2.IpPermission{
						IpProtocol: &protocol,
						UserIdGroupPairs: []*ec2.UserIdGroupPair{
							&ec2.UserIdGroupPair{
								GroupId: &e.Provider.SecurityGroup,
							},
						},
					},
				},
			}
			_, err = svc.AuthorizeSecurityGroupIngress(asgii)
		}
	}
	if err != nil && !strings.Contains(err.Error(), "InvalidPermission.Duplicate") {
		return err
	}
	var ports int64
	ports = 22
	cidrIP := "0.0.0.0/0"
	asgii = &ec2.AuthorizeSecurityGroupIngressInput{
		FromPort:   &ports,
		ToPort:     &ports,
		GroupName:  &name,
		IpProtocol: &protocol,
		CidrIp:     &cidrIP,
	}
	if _, err = svc.AuthorizeSecurityGroupIngress(asgii); err != nil {
		if strings.Contains(err.Error(), "InvalidParameterValue") {
			asgii.SetGroupName("")
			asgii.SetGroupId(e.Provider.SecurityGroup)
			_, err = svc.AuthorizeSecurityGroupIngress(asgii)
		}
	}
	if err != nil && !strings.Contains(err.Error(), "InvalidPermission.Duplicate") {
		return err
	}
	return nil
}
