package templates

import (
	"encoding/json"
	"testing"

	gocf "github.com/crewjam/go-cloudformation"
)

func TestSecurityGroupSerialization(t *testing.T) {
	template := gocf.NewTemplate()
	instanceIngressRules := gocf.EC2SecurityGroupRuleList{
		gocf.EC2SecurityGroupRule{
			SourceSecurityGroupIdXXSecurityGroupIngressXOnlyX: gocf.String("sg-11111"),
			IpProtocol: gocf.String("tcp"),
			FromPort:   gocf.Integer(80),
			ToPort:     gocf.Integer(80),
		},
	}
	securityGroup := &gocf.EC2SecurityGroup{
		GroupDescription:     gocf.String("Description"),
		SecurityGroupIngress: &instanceIngressRules,
		VpcId:                gocf.String("vpc-12345"),
	}
	template.AddResource("SecurityGroup", securityGroup)
	marshalledTemplate, err := json.Marshal(template)
	if err != nil {
		t.Fatalf("Error serialization: %s", err.Error())
	}
	if string(marshalledTemplate) != `{"AWSTemplateFormatVersion":"2010-09-09","Resources":{"SecurityGroup":{"Type":"AWS::EC2::SecurityGroup","Properties":{"GroupDescription":"Description","SecurityGroupIngress":[{"FromPort":80,"IpProtocol":"tcp","SourceSecurityGroupId":"sg-11111","ToPort":80}],"VpcId":"vpc-12345"}}}}` {
		t.Fatalf("Wrong serialization, modify 'github.com/crewjam/go-cloudformation/schemas.go': %s", string(marshalledTemplate))
	}
}
