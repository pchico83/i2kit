package templates

import (
	"fmt"

	gocf "github.com/crewjam/go-cloudformation"
	"github.com/pchico83/i2kit/cli/schemas/environment"
	"github.com/pchico83/i2kit/cli/schemas/service"
)

func route53(t *gocf.Template, s *service.Service, e *environment.Environment) {
	var resourceRecords *gocf.StringListExpr
	var dependsOn []string
	if s.Stateful {
		resourceRecords = gocf.StringList(gocf.GetAtt("EC2Instance", "PublicDnsName"))
		dependsOn = append(dependsOn, "EIP")
	} else {
		resourceRecords = gocf.StringList(gocf.GetAtt("ELB", "DNSName"))
	}
	recordName := fmt.Sprintf("%s.%s", s.GetFullName(e, "."), e.Provider.HostedZone)
	recordSetProperties := &gocf.Route53RecordSet{
		HostedZoneName:  gocf.String(e.Provider.HostedZone),
		Name:            gocf.String(recordName),
		Type:            gocf.String("CNAME"),
		TTL:             gocf.String("60"),
		ResourceRecords: resourceRecords,
	}
	resourceSetResource := &gocf.Resource{
		Properties: recordSetProperties,
		DependsOn:  dependsOn,
	}
	t.Resources["DNSRecord"] = resourceSetResource
}
