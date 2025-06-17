package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"schedule-sample/schedule"
)

func main() {
	// クライアントとワーカーは、プロセスごとに一度だけ作成されるべき重いオブジェクトです。
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("クライアントの作成ができません", err)
	}
	defer c.Close()

	w := worker.New(c, "schedule", worker.Options{})

	w.RegisterWorkflow(schedule.SampleScheduleWorkflow)
	w.RegisterActivity(schedule.DoSomething)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("ワーカーの開始ができません", err)
	}
}
