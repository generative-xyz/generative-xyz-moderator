package main

import "rederinghub.io/tools"

func main() {
	uc := tools.StartFactory()
	if uc != nil {
		uc.ReportPerceptronOwners()
	}
}
