package service

import (
	"reflect"
	"testing"
)

func Test_extractPorts(t *testing.T) {
	tests := []struct {
		name    string
		ports   []string
		want    []*Port
		wantErr bool
	}{
		{
			name:  "https-with-cert",
			ports: []string{"https:443:http:8000:arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8"},
			want: []*Port{&Port{
				Protocol:         "HTTPS",
				Port:             "443",
				InstanceProtocol: "HTTP",
				InstancePort:     "8000",
				Certificate:      "arn:aws:acm:us-west-2:062762192540:certificate/12de3ac5-abcd-461a-1234-9e81250b33d8",
			}},
			wantErr: false,
		},
		{
			name:  "http-to-http",
			ports: []string{"http:8000:http:80"},
			want: []*Port{&Port{
				Protocol:         "HTTP",
				Port:             "8000",
				InstanceProtocol: "HTTP",
				InstancePort:     "80",
			}},
			wantErr: false,
		},
		{
			name:    "https-to-http-no-cert",
			ports:   []string{"https:8000:http:80"},
			want:    nil,
			wantErr: true,
		},
		{
			name:    "unknown-schema",
			ports:   []string{"ftp:8000:http:80"},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pointerPorts := make([]*string, len(tt.ports))
			for i := range tt.ports {
				pointerPorts[0] = &tt.ports[i]
			}

			got, err := extractPorts(tt.name, pointerPorts)
			if (err != nil) != tt.wantErr {
				t.Errorf("extractPorts() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(&got, &tt.want) {
				t.Errorf("extractPorts() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
