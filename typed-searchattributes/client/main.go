package main

import (
	"context"
	"log"

	"github.com/pborman/uuid"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/temporal"

	typedsearchattributes "typed-searchattributes/attribute"
)

func main() {
	// クライアントは、プロセスごとに1回作成すべき重量級オブジェクトです。
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("クライアントの作成ができません", err)
	}
	defer c.Close()

	workflowOptions := client.StartWorkflowOptions{
		ID:                    "typed-search_attributes_" + uuid.New(),
		TaskQueue:             "typed-search-attributes",
		TypedSearchAttributes: temporal.NewSearchAttributes(typedsearchattributes.CustomIntKey.ValueSet(1)),
	}

	we, err := c.ExecuteWorkflow(context.Background(), workflowOptions, typedsearchattributes.SearchAttributesWorkflow)
	if err != nil {
		log.Fatalln("ワークフローの実行ができません", err)
	}
	log.Println("ワークフローを開始しました", "WorkflowID", we.GetID(), "RunID", we.GetRunID())
}
