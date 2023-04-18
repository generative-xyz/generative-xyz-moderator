package usecase

import "testing"

func Test_executeAISchoolJob(t *testing.T) {
	type args struct {
		params  string
		dataset string
		output  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test executeAISchoolJob",
			args: args{
				params:  "./datatest/params.json",
				dataset: "./datatest/data_test",
				output:  "./datatest/output"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := executeAISchoolJob(tt.args.params, tt.args.dataset, tt.args.output); (err != nil) != tt.wantErr {
				t.Errorf("executeAISchoolJob() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
