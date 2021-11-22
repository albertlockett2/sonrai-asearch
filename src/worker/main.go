package worker

import "log"

func Main() {
	worker, err := NewWorker()
	if err != nil {
		log.Fatalf("Failed create worker: %v", err)
	}

	err = worker.Start()
	if err != nil {
		log.Fatalf("Failed start worker: %v", err)
	}
}
