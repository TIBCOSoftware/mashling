package app

import (
	"testing"

	"github.com/TIBCOSoftware/mashling/cli/cli"
)

func Test_cmdVersion_Exec(t *testing.T) {
	type fields struct {
		option        *cli.OptionInfo
		versionNumber string
	}
	type args struct {
		args []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "TestCase1",
			args: args{},
			fields: fields{
				option:        nil,
				versionNumber: "1.8.0",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &cmdVersion{
				option:        tt.fields.option,
				versionNumber: tt.fields.versionNumber,
			}
			if err := c.Exec(tt.args.args); (err != nil) != tt.wantErr {
				t.Errorf("cmdVersion.Exec() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
