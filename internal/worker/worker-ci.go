package worker

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/golang/protobuf/proto"

	grpcTypes "github.com/kitanoyoru/kita-proto/go"

	"github.com/kitanoyoru/kita-ci/pkg/structs"
	"github.com/kitanoyoru/kita-ci/pkg/utils"
)

func (w *CIWorker) StartConsuming() {
	execFlag := make(chan bool)

	msgs, err := w.jobsQueue.MakeCIMsgChan()
	if err != nil {
		w.logger.Fatal("Failed to create CI jobs channel", err)
	}

	go func() {
		for m := range msgs {
			w.logger.Info("Received message: " + string(m))
			jobMsg := grpcTypes.CIJob{}
			if err := proto.Unmarshal(m, jobMsg); err != nil {
				w.logger.Error("Failed to unmarshal message", err)
				continue
			}
			go w.executeCIJob(jobMsg)
		}
	}()

	w.logger.Info("CI Worker started")
	<-execFlag
}

func (w *CIWorker) executeCIJob(job grpcTypes.CIJob) {
	payload := w.createDockerBuilderPayload(job)
	logs, err := w.dockerClient.RunCIBuilderContainer(payload)

	success := true
	if err != nil {
		w.logger.Error("Build failed", err)
		success = false
	}

	// TODO: Send artifact id to cd worker
	_, err := w.writeBuildToDb(job.RepoID, success, job.Branch, payload.Tag, logs)
	if err != nil {
		w.logger.Error("Build failed", err)
		return
	}

	err = w.sendInfoToGithub(payload.Username, job.AccessToken, payload.RepoName, job.HeadSHA, success)
	if err != nil {
		w.logger.Error("Write Guthub status faield", err)
		return
	}
}

func (w *CIWorker) createDockerBuilderPayload(job grpcTypes.CIJob) structs.CIBuilderPayload {
	payload := structs.CIBuilderPayload{
		RepoURL: job.RepoURL,
		Branch:  job.Branch,
	}

	payload.RepoName = utils.GetRepoNameFromGithubURL(job.RepoURL)
	payload.Username = utils.GetUsernameFromGithubURL(job.RepoURL)
	payload.Tag = utils.MakeContainerTag(job.RepoName, job.Branch, job.HeadSHA)

	if job.AccessToken != "" {
		payload.RepoURL = utils.MakeRepoUrlWithAccessToken(job.Username, job.AccessToken, job.RepoURL)
	}

	return payload
}

func (w *CIWorker) writeBuildToDb(repoId int64, success bool, branch, tag, stdout string) (int64, err) {
	build := structs.Build{
		GithubRepoID:  repoId,
		Branch:        branch,
		IsSuccessfull: success,
		CreatedAt:     time.Now(),
		Stdout:        stdout,
	}

	err := w.dbClient.InsertBuild(&build)
	if err != nil {
		return 0, nil
	}

	artifact := structs.Artifact{
		BuildID: build.ID,
		Name:    tag,
	}

	err = w.dbClient.InsertArtifact(&artifact)
	if err != nil {
		return 0, err
	}

	return artifact.ID, nil
}

func (w *CIWorker) sendInfoToGithub(username, accessToken, repo, sha string, success bool) error {
	client := http.Client{}
	url := fmt.Sprintf("https://api.github.com/repos/%s/statuses/%s", repo, sha)
	body := &structs.GithubStatusMessage{
		Description: "Status set by Kita CI worker",
		Context:     "ci-build",
	}
	if success == true {
		body.State = "success"
	} else {
		body.State = "failure"
	}

	rawBody, err := json.Marshal(body)
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(rawBody)

	// TODO: Create Requests pkg
	req, err := http.NewRequest("POST", url, buf)
	if err != nil {
		return err
	}
	req.SetBasicAuth(username, accessToken)
	resp, err := client.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()

	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	w.logger.Info("Github Response: " + string(respBody))

	return nil
}
