package main

import (
	"log"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	"temporal-sample/app"
)

func main() {

	c, err := client.Dial(client.Options{})
	if err != nil {
		log.Fatalln("Temporalクライアントの作成に失敗しました。", err)
	}
	defer c.Close()

	w := worker.New(c, app.MoneyTransferTaskQueueName, worker.Options{})

	// このワーカーはワークフロー関数とアクティビティ関数の両方をホストします。
	w.RegisterWorkflow(app.MoneyTransfer)
	w.RegisterActivity(app.Withdraw)
	w.RegisterActivity(app.Deposit)
	w.RegisterActivity(app.Refund)

	// タスクキューのリッスンを開始します。
	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("ワーカーの起動に失敗しました", err)
	}
}
