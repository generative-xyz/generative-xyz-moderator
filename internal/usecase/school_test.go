package usecase

import (
	"fmt"
	"path/filepath"
	"testing"
)

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
			filep, _ := filepath.Abs("./ai-school-work/")
			fmt.Println(filep)
			// progCh := make(chan JobProgress)

			// go func() {
			// 	for prog := range progCh {
			// 		fmt.Println(prog.Epoch)
			// 	}
			// }()

			// if err := executeAISchoolJob("./training_user.py", tt.args.params, tt.args.dataset, tt.args.output, progCh); (err != nil) != tt.wantErr {
			// 	t.Errorf("executeAISchoolJob() error = %v, wantErr %v", err, tt.wantErr)
			// }
		})
	}
}
