package iplocate

import (
	"fmt"
	"time"

	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

// GetAddressFromIPは、IPアドレスと位置情報を取得するTemporalワークフローです。
func GetAddressFromIP(ctx workflow.Context, name string) (string, error) {
	// アクティビティオプション（リトライポリシーを含む）を定義
	ao := workflow.ActivityOptions{
		StartToCloseTimeout: time.Minute,
		RetryPolicy: &temporal.RetryPolicy{
			InitialInterval:    time.Second, // 最初のリトライまでの待機時間
			MaximumInterval:    time.Minute, // リトライ間隔の最大値
			BackoffCoefficient: 2,           // リトライ間隔の増加係数
			// MaximumAttempts: 5, // 試行回数を制限したい場合はコメントを外してください
		},
	}
	ctx = workflow.WithActivityOptions(ctx, ao)

	var ipActivities *IPActivities

	var ip string
	err := workflow.ExecuteActivity(ctx, ipActivities.GetIP).Get(ctx, &ip)
	if err != nil {
		return "", fmt.Errorf("IPの取得に失敗しました: %s", err)
	}

	var location string
	err = workflow.ExecuteActivity(ctx, ipActivities.GetLocationInfo, ip).Get(ctx, &location)
	if err != nil {
		return "", fmt.Errorf("位置情報の取得に失敗しました: %s", err)
	}
	return fmt.Sprintf("こんにちは、%s。あなたのIPは%s、位置情報は%sです", name, ip, location), nil
}
