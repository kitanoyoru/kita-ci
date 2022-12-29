package worker

import (
	"fmt"
	"strings"

	"github.com/golang/protobuf/proto"

	grpcTypes "github.com/kitanoyoru/kita-proto/go"

	"github.com/kitanoyoru/kita-ci/pkg/structs"
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
	payload := w.createDockerBuilderPayload()
	logs, err := w.dockerClient.RunCIBuilderContainer(payload)

	success := true
	if err != nil {
		w.logger.Error("Build failed", err)
		success = false
	}

	// TODO: Write to DB
}

func (w *CIWorker) createDockerBuilderPayload(job grpcTypes.CIJob) structs.CIBuilderPayload {
	payload := structs.CIBuilderPayload{}

	// REFACTOR: process string logic move to utils
	payload.RepoName = strings.TrimSuffix(strings.TrimPrefix(job.RepoURL, "https://github.com/"), ".git")
	payload.Username = strings.Split(strings.TrimPrefix(job.RepoURL, "https://github.com/"), "/")[0]
	payload.Branch = job.Branch
	payload.Tag = fmt.Sprintf("%s:%s-%s", payload.RepoName, payload.Branch, job.HeadSHA[:7])

	repoURL := job.RepoURL
	if job.AccessToken != "" {
		repoURL = fmt.Sprintf("https://%s:%s@%s", payload.Username, job.AccessToken, strings.TrimPrefix(job.RepoURL, "https://"))
	}
	payload.RepoURL = repoURL

	return payload
}
