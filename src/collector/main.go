package collector

import "log"

func Main() {
	collector, err := NewCollector()
	if err != nil {
		log.Fatalf("failed to create collecotr: %v", err)
	}

	err = collector.Start()
	if err != nil {
		log.Fatalf("failed to start collector: %v", err)
	}
}
