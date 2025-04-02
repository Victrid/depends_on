package resolver

import (
	"reflect"
	"testing"
)

func Test_parseDependsOnString(t *testing.T) {
	type args struct {
		input string
	}
	tests := []struct {
		name    string
		args    args
		want    []*Dependency
		wantErr bool
	}{
		{
			name: "valid input",
			args: args{
				input: "service:service,statefulset:ns/set?good,deployment:ns2/dep?bad",
			},
			want: []*Dependency{
				{
					Resource: "service",
					Locator: Locator{
						Namespace: nil,
						Name:      "service",
					},
					Status: nil,
				},
				{
					Resource: "statefulset",
					Locator: Locator{
						Namespace: ptr("ns"),
						Name:      "set",
					},
					Status: ptr("good"),
				},
				{
					Resource: "deployment",
					Locator: Locator{
						Namespace: ptr("ns2"),
						Name:      "dep",
					},
					Status: ptr("bad"),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseDependsOnString(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseDependsOnString() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseDependsOnString() got = %v, want %v", got, tt.want)
			}
		})
	}
}
