package templates

import (
	"fmt"

	gocf "github.com/crewjam/go-cloudformation"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

func logGroup(t *gocf.Template, s *service.Service, e *environment.Environment) {
	logGroupIAM(t, s, e)
	// randomID := uuid.New()
	logGroupResource := &gocf.Resource{
		Properties: &gocf.LogsLogGroup{
			LogGroupName:    gocf.String(fmt.Sprintf("i2kit-%s", s.GetFullName(e, "-"))),
			RetentionInDays: gocf.Integer(30),
		},
		// DeletionPolicy: "Retain",
	}
	t.Resources["LogGroup"] = logGroupResource
}

func logGroupIAM(t *gocf.Template, s *service.Service, e *environment.Environment) {
	policy := gocf.IAMPolicies{
		PolicyName: gocf.String(s.GetFullName(e, "-")),
		PolicyDocument: &map[string]interface{}{
			"Version": "2012-10-17",
			"Statement": map[string]interface{}{
				"Effect":   "Allow",
				"Action":   []string{"logs:CreateLogStream", "logs:PutLogEvents"},
				"Resource": fmt.Sprintf("arn:aws:logs:%s:*:log-group:i2kit-%s:log-stream:*", e.Provider.Region, s.GetFullName(e, "-")),
			},
		},
	}
	role := &gocf.IAMRole{
		AssumeRolePolicyDocument: &map[string]interface{}{
			"Statement": map[string]interface{}{
				"Effect":    "Allow",
				"Principal": map[string]interface{}{"Service": []string{"ec2.amazonaws.com"}},
				"Action":    []string{"sts:AssumeRole"},
			},
		},
		Path:     gocf.String("/"),
		Policies: &gocf.IAMPoliciesList{policy},
	}
	t.AddResource("Role", role)
	instanceProfile := &gocf.IAMInstanceProfile{
		Path:  gocf.String("/"),
		Roles: gocf.StringList(gocf.Ref("Role")),
	}
	t.AddResource("InstanceProfile", instanceProfile)
}
