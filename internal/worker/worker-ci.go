package worker

import (
	"github.com/golang/protobuf/proto"

	grpcTypes "github.com/kitanoyoru/kita-proto/go"
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
  
}
