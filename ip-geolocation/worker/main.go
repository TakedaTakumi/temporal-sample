package main

import (
	"log"
	"net/http"

	"temporal-ip-geolocation/iplocate"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"
)

func main() {
	// Temporalクライアントを作成
	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Temporalクライアントの作成に失敗しました", err)
	}
	defer c.Close()

	// Temporalワーカーを作成
	w := worker.New(c, iplocate.TaskQueueName, worker.Options{})

	// HTTPクライアントをActivities構造体に注入
	activities := &iplocate.IPActivities{
		HTTPClient: http.DefaultClient,
	}

	// ワークフローとアクティビティを登録
	w.RegisterWorkflow(iplocate.GetAddressFromIP)
	w.RegisterActivity(activities)

	// ワーカーを開始
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("Temporalワーカーの起動に失敗しました", err)
	}
}
