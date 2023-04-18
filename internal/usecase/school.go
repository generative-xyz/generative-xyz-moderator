package usecase

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"rederinghub.io/internal/entity"
	"rederinghub.io/internal/usecase/structure"
)

type JobProgress struct {
	Epoch       int
	IsCompleted bool
}

type AIJobInstance struct {
	u            Usecase
	job          *entity.AISchoolJob
	currentEpoch int
	IsCompleted  bool
	progCh       chan JobProgress
}

var currentAIJobs map[string]AIJobInstance

func (u Usecase) JobAIS_WatchPending() error {
	jobList, err := u.Repo.GetAISchoolJobByStatus([]string{"running", "waiting"})
	if err != nil {
		return err
	}

	if currentAIJobs == nil {
		currentAIJobs = make(map[string]AIJobInstance)
	}
	for jobID, job := range currentAIJobs {
		if job.IsCompleted {
			delete(currentAIJobs, jobID)
		}
	}

	for _, job := range jobList {
		if job.Status == "waiting" {
			jobParams := &structure.AISchoolModelParams{}
			err := json.Unmarshal([]byte(job.Params), jobParams)
			if err != nil {
				// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while unmarshalling job params: "+err.Error(), "error")
				continue
			}

			newJob := AIJobInstance{
				u:   u,
				job: &job,
			}
			currentAIJobs[job.JobID] = newJob
			job.Status = "running"
			job.ExecutedAt = time.Now().Unix()
			err = u.Repo.UpdateAISchoolJob(&job)
			if err != nil {
				// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
				continue
			}
			go newJob.Start()
		}
		if job.Status == "running" {
			if _, exist := currentAIJobs[job.JobID]; !exist {
				job.Status = "waiting"
				job.ExecutedAt = 0
				job.Progress = 0
				err = u.Repo.UpdateAISchoolJob(&job)
				if err != nil {
					// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
					continue
				}
			}
		}
	}
	return nil
}

const basePath = "./ai-school-work/"

func createAISchoolWorkFolder(jobID string, params structure.AISchoolModelParams) error {
	if err := os.MkdirAll(basePath+jobID, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(basePath+jobID+"/dataset", os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(basePath+jobID+"/output", os.ModePerm); err != nil {
		return err
	}
	content, err := json.Marshal(params)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(basePath+jobID+"/params.json", content, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func clearAISchoolWorkFolder(jobID string) error {
	err := os.RemoveAll("./ai-school-work/" + jobID + "/")
	if err != nil {
		return err
	}
	return nil
}

func (job *AIJobInstance) Start() {
	defer func() {
		job.IsCompleted = true
	}()
	jobID := job.job.JobID
	err := clearAISchoolWorkFolder(jobID)
	if err != nil {
		job.job.Errors = err.Error()
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
	}
	params := structure.AISchoolModelParams{}
	err = json.Unmarshal([]byte(job.job.Params), &params)
	if err != nil {
		job.job.Errors = err.Error()
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
	}
	err = createAISchoolWorkFolder(jobID, params)
	if err != nil {
		job.job.Errors = err.Error()
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
	}

	progCh := make(chan JobProgress)
	job.progCh = progCh
	go func() {
		for prog := range job.progCh {
			job.job.Progress = prog.Epoch
			err = job.u.Repo.UpdateAISchoolJob(job.job)
			if err != nil {
				// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
				log.Println(err)
			}
		}
	}()
	scriptPath := os.Getenv("AI_SCHOOL_SCRIPT")
	jobPath := basePath + jobID
	err = executeAISchoolJob(scriptPath, jobPath+"/params.json", jobPath+"/dataset", jobPath+"/output", job.progCh)
	if err != nil {
		job.job.Errors = err.Error()
		err = job.u.Repo.UpdateAISchoolJob(job.job)
		if err != nil {
			// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
			return
		}
	}

	job.job.CompletedAt = time.Now().Unix()
	job.job.Status = "completed"
	err = job.u.Repo.UpdateAISchoolJob(job.job)
	if err != nil {
		// go u.Slack.SendMessageToSlackWithChannel("Error", "Error while updating job status: "+err.Error(), "error")
		return
	}
}
func executeAISchoolJob(scriptPath string, params string, dataset string, output string, progCh chan JobProgress) error {
	// 1. Get params
	// 2. Get dataset
	// 3. Run job
	// 4. Update job
	args := fmt.Sprintf("%v -c %v -d %v -o %v", scriptPath, params, dataset, output)
	cmd := exec.Command("python3", strings.Split(args, " ")...)
	// cmd := exec.Command("ls", "-a")
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Start()
	scanner := bufio.NewScanner(stderr)
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		m := scanner.Text()
		fmt.Println("err", m)
	}

	scanner2 := bufio.NewScanner(stdout)
	scanner2.Split(bufio.ScanLines)
	for scanner2.Scan() {
		m := scanner2.Text()
		if strings.Contains(strings.ToLower(m), "epoch") {
			epochStr := strings.Split(m, "Epoch ")
			epochs := strings.Split(epochStr[1], "/")
			currentEpoch := epochs[0]
			currentEpochInt, err := strconv.ParseInt(currentEpoch, 10, 64)
			if err != nil {
				return err
			}
			progCh <- JobProgress{
				Epoch: int(currentEpochInt),
			}
		}
	}

	cmd.Wait()
	time.Sleep(100 * time.Millisecond)
	close(progCh)
	return nil
}
