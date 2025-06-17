package schedule

import (
	"context"
	"time"

	"go.temporal.io/sdk/activity"
	"go.temporal.io/sdk/converter"
	"go.temporal.io/sdk/workflow"
)

// SampleScheduleWorkflow は指定されたスケジュールで実行されます
func SampleScheduleWorkflow(ctx workflow.Context) error {

	workflow.GetLogger(ctx).Info("スケジュールワークフローが開始されました。", "StartTime", workflow.Now(ctx))

	ao := workflow.ActivityOptions{
		StartToCloseTimeout: 10 * time.Second,
	}
	ctx1 := workflow.WithActivityOptions(ctx, ao)

	info := workflow.GetInfo(ctx1)

	// スケジュールによって開始されたワークフロー実行には、以下の追加プロパティが検索属性に追加されます
	//lint:ignore SA1019 - これはサンプルです
	scheduledByIDPayload := info.SearchAttributes.IndexedFields["TemporalScheduledById"]
	var scheduledByID string
	err := converter.GetDefaultDataConverter().FromPayload(scheduledByIDPayload, &scheduledByID)
	if err != nil {
		return err
	}
	//lint:ignore SA1019 - これはサンプルです
	startTimePayload := info.SearchAttributes.IndexedFields["TemporalScheduledStartTime"]
	var startTime time.Time
	err = converter.GetDefaultDataConverter().FromPayload(startTimePayload, &startTime)
	if err != nil {
		return err
	}

	err = workflow.ExecuteActivity(ctx1, DoSomething, scheduledByID, startTime).Get(ctx, nil)
	if err != nil {
		workflow.GetLogger(ctx).Error("スケジュールワークフローが失敗しました。", "Error", err)
		return err
	}

	return nil
}

// DoSomething はアクティビティです
func DoSomething(ctx context.Context, scheduleByID string, startTime time.Time) error {
	activity.GetLogger(ctx).Info("スケジュールジョブが実行中です。", "scheduleByID", scheduleByID, "startTime", startTime)
	// データベースへのクエリ、外部APIの呼び出し、またはその他の非決定的なアクションを実行します。
	return nil
}
