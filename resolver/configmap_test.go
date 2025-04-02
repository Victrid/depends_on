package resolver

import (
	"reflect"
	"testing"
)

func Test_parseConfigMapDependency(t *testing.T) {
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
			name: "regular",
			args: args{
				input: `
- raw: "service:postgresql?available"
- resource: "deployment"
  name: "prometheus"
  namespace: "metric"`,
			},
			want: []*Dependency{
				{
					Resource: "service",
					Locator: Locator{
						Namespace: nil,
						Name:      "postgresql",
					},
					Status: ptr("available"),
				},
				{
					Resource: "deployment",
					Locator: Locator{
						Namespace: ptr("metric"),
						Name:      "prometheus",
					},
					Status: nil,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseConfigMapDependency(tt.args.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseConfigMapDependency() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseConfigMapDependency() got = %v, want %v", got, tt.want)
			}
		})
	}
}
