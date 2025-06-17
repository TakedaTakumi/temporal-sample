package app

import (
    "fmt"
    "time"
    "go.temporal.io/sdk/workflow"
    "go.temporal.io/sdk/temporal"
)

func MoneyTransfer(ctx workflow.Context, input PaymentDetails) (string, error) {

    // RetryPolicyは、アクティビティが失敗した場合に自動的に再試行を処理する方法を指定します。
    retryPolicy := &temporal.RetryPolicy{
        InitialInterval:        time.Second, // 最初の再試行までの待機時間
        BackoffCoefficient:     2.0,            // 待機時間の増加係数
        MaximumInterval:        100 * time.Second, // 最大待機時間
        MaximumAttempts:        500, // 0は無制限の再試行を意味します
        NonRetryableErrorTypes: []string{"InvalidAccountError", "InsufficientFundsError"}, // 再試行しないエラーのタイプ
    }

    options := workflow.ActivityOptions{
        // タイムアウトオプションは、アクティビティ関数を自動的にタイムアウトするタイミングを指定します。
        StartToCloseTimeout: time.Minute,
        // オプションでカスタマイズされたRetryPolicyを提供します。
        // Temporalはデフォルトで失敗したアクティビティを再試行します。
        RetryPolicy: retryPolicy,
    }

    // オプションを適用します。
    ctx = workflow.WithActivityOptions(ctx, options)

    // お金を引き出します。
    var withdrawOutput string

    withdrawErr := workflow.ExecuteActivity(ctx, Withdraw, input).Get(ctx, &withdrawOutput)

    if withdrawErr != nil {
        return "", withdrawErr
    }

    // お金を入金します。
    var depositOutput string

    depositErr := workflow.ExecuteActivity(ctx, Deposit, input).Get(ctx, &depositOutput)

    if depositErr != nil {
        // 入金に失敗しました。元のアカウントにお金を戻します。

        var result string

        refundErr := workflow.ExecuteActivity(ctx, Refund, input).Get(ctx, &result)

        if refundErr != nil {
            return "",
                fmt.Errorf("Deposit: failed to deposit money into %v: %v. Money could not be returned to %v: %w",
                    input.TargetAccount, depositErr, input.SourceAccount, refundErr)
        }

        return "", fmt.Errorf("Deposit: failed to deposit money into %v: Money returned to %v: %w",
            input.TargetAccount, input.SourceAccount, depositErr)
    }

    result := fmt.Sprintf("Transfer complete (transaction IDs: %s, %s)", withdrawOutput, depositOutput)
    return result, nil
}
