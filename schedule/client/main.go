package main

import (
	"context"
	"log"
	"time"

	"schedule-sample/schedule"

	"github.com/pborman/uuid"
	"go.temporal.io/api/common/v1"
	"go.temporal.io/api/enums/v1"
	"go.temporal.io/sdk/client"
	"go.temporal.io/sdk/converter"
)

func main() {
	ctx := context.Background()
	// クライアントは、プロセスごとに一度だけ作成されるべき重いオブジェクトです。
	c, err := client.Dial(client.Options{
		HostPort: client.DefaultHostPort,
	})
	if err != nil {
		log.Fatalln("クライアントの作成ができません", err)
	}
	defer c.Close()
	// このスケジュールIDはユーザーのビジネスロジック識別子としても使用できます。
	scheduleID := "schedule_" + uuid.New()
	workflowID := "schedule_workflow_" + uuid.New()
	// スケジュールを作成し、スケジュールが実行されないように仕様なしで開始します。
	scheduleHandle, err := c.ScheduleClient().Create(ctx, client.ScheduleOptions{
		ID:   scheduleID,
		Spec: client.ScheduleSpec{},
		Action: &client.ScheduleWorkflowAction{
			ID:        workflowID,
			Workflow:  schedule.SampleScheduleWorkflow,
			Args:      []interface{}{"Input Args"}, // workflowの引数
			TaskQueue: "schedule",
		},
	})
	if err != nil {
		log.Fatalln("スケジュールの作成ができません", err)
	}
	// サンプルが終了したらスケジュールを削除します
	defer func() {
		log.Println("スケジュールを削除しています", "ScheduleID", scheduleHandle.GetID())
		err = scheduleHandle.Delete(ctx)
		if err != nil {
			log.Fatalln("スケジュールの削除ができません", err)
		}
	}()
	// スケジュールを手動で一度トリガーします
	log.Println("スケジュールを手動でトリガーしています", "ScheduleID", scheduleHandle.GetID())

	err = scheduleHandle.Trigger(ctx, client.ScheduleTriggerOptions{
		Overlap: enums.SCHEDULE_OVERLAP_POLICY_ALLOW_ALL,
	})
	if err != nil {
		log.Fatalln("スケジュールのトリガーができません", err)
	}
	// スケジュールを仕様で更新して定期的に実行されるようにします
	log.Println("スケジュールを更新しています", "ScheduleID", scheduleHandle.GetID())
	err = scheduleHandle.Update(ctx, client.ScheduleUpdateOptions{
		DoUpdate: func(schedule client.ScheduleUpdateInput) (*client.ScheduleUpdate, error) {
			schedule.Description.Schedule.Spec = &client.ScheduleSpec{
				// 金曜日の午後5時にスケジュールを実行します
				Calendars: []client.ScheduleCalendarSpec{
					{
						Hour: []client.ScheduleRange{
							{
								Start: 17,
							},
						},
						DayOfWeek: []client.ScheduleRange{
							{
								Start: 5,
							},
						},
					},
				},
				// 5秒ごとにスケジュールを実行します
				Intervals: []client.ScheduleIntervalSpec{
					{
						Every: 5 * time.Second,
					},
				},
			}
			// スケジュールの一時停止を解除する方法を示すために、一時停止状態で開始します
			schedule.Description.Schedule.State.Paused = true
			schedule.Description.Schedule.State.LimitedActions = true
			schedule.Description.Schedule.State.RemainingActions = 10

			return &client.ScheduleUpdate{
				Schedule: &schedule.Description.Schedule,
			}, nil
		},
	})
	if err != nil {
		log.Fatalln("スケジュールの更新ができません", err)
	}
	// スケジュールの一時停止を解除します
	log.Println("スケジュールの一時停止を解除しています", "ScheduleID", scheduleHandle.GetID())
	err = scheduleHandle.Unpause(ctx, client.ScheduleUnpauseOptions{})
	if err != nil {
		log.Fatalln("スケジュールの一時停止解除ができません", err)
	}
	// スケジュールが10アクションを実行するのを待ちます
	log.Println("スケジュールが10アクションを完了するのを待っています", "ScheduleID", scheduleHandle.GetID())

	for {
		description, err := scheduleHandle.Describe(ctx)
		if err != nil {
			log.Fatalln("スケジュールの説明ができません", err)
		}
		if description.Schedule.State.RemainingActions != 0 {
			log.Println("スケジュールには残りのアクションがあります", "ScheduleID", scheduleHandle.GetID(), "RemainingActions", description.Schedule.State.RemainingActions)

			// workflowの引数を取得する
			if action, ok := description.Schedule.Action.(*client.ScheduleWorkflowAction); ok {
				if payload, ok := action.Args[0].(*common.Payload); ok {
					var inputString string
					if err := converter.GetDefaultDataConverter().FromPayload(payload, &inputString); err != nil {
						log.Fatalln("Failed to convert payload", err)
					}

					log.Println("Input:", inputString)
				}
			}

			time.Sleep(5 * time.Second)
		} else {
			break
		}
	}
}
