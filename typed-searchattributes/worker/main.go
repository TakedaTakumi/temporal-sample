package main

import (
	"log"

	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/worker"

	typedsearchattributes "typed-searchattributes/attribute"
)

func main() {
	// クライアントとワーカーは、プロセスごとに1回作成すべき重量級オブジェクトです。
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("クライアントの作成ができません", err)
	}
	defer c.Close()

	w := worker.New(c, "typed-search-attributes", worker.Options{})

	w.RegisterWorkflow(typedsearchattributes.SearchAttributesWorkflow)

	err = w.Run(worker.InterruptCh())
	if err != nil {
		log.Fatalln("ワーカーの起動ができません", err)
	}
}
