package templates

import (
	"fmt"

	gocf "github.com/crewjam/go-cloudformation"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

func logGroup(t *gocf.Template, s *service.Service, e *environment.Environment) {
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
