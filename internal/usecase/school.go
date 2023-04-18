package usecase

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"

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

			// err = executeAISchoolJob(jobParams, nil)
			// if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while executing job: "+err.Error(), "error")
			// 	continue
			// }
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

func createAISchoolWorkFolder(jobID string, params structure.AISchoolModelParams, dataset string) (string, error) {
	return "", nil
}

func executeAISchoolJob(params string, dataset string, output string) error {
	// 1. Get params
	// 2. Get dataset
	// 3. Run job
	// 4. Update job
	args := fmt.Sprintf("training_user.py -c %v -d %v -o %v", params, dataset, output)
	cmd := exec.Command("python3", strings.Split(args, " ")...)
	stderr, _ := cmd.StderrPipe()
	cmd.Start()

	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanWords)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println(m)
	}
	cmd.Wait()
	return nil
}
