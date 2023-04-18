package usecase

import (
	"encoding/json"

	"rederinghub.io/internal/usecase/structure"
)

func (u Usecase) JobAIS_WatchPending() error {

	jobList, err := u.Repo.GetAISchoolJobByStatus([]string{"running", "waiting"})
	if err != nil {
		return err
	}

	for _, job := range jobList {
		if job.Status == "waiting" {
			jobParams := &structure.AISchoolModelParams{}
			err := json.Unmarshal([]byte(job.Params), jobParams)
			if err != nil {
				// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while unmarshalling job params: "+err.Error(), "error")
				continue
			}
			execParams := jobParams.TransformToExec()
			err = executeAISchoolJob(execParams, nil)
			if err != nil {
				// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while executing job: "+err.Error(), "error")
				continue
			}
			job.Status = "running"
			err = u.Repo.UpdateAISchoolJob(&job)
			if err != nil {
				// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
				continue
			}
		}
		if job.Status == "running" {

		}
	}
	return nil
}

func executeAISchoolJob(params *structure.AISchoolModelParamsExec, dataset interface{}) error {
	// 1. Get params
	// 2. Get dataset
	// 3. Run job
	// 4. Update job
	return nil
}
