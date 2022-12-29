package worker

import (
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
