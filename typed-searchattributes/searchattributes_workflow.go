package typedsearchattributes

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"go.temporal.io/sdk/log"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/workflow"
)

/*
 * このサンプルは、型付き検索属性APIの使用方法を示しています。
 */

var (
	// CustomIntKeyは、カスタムint検索属性のキーです。
	CustomIntKey = temporal.NewSearchAttributeKeyInt64("CustomIntField")
	// CustomKeywordは、カスタムキーワード検索属性のキーです。
	CustomKeyword = temporal.NewSearchAttributeKeyString("CustomKeywordField")
	// CustomBoolは、カスタムブール検索属性のキーです。
	CustomBool = temporal.NewSearchAttributeKeyBool("CustomBoolField")
	// CustomDoubleは、カスタムdouble検索属性のキーです。
	CustomDouble = temporal.NewSearchAttributeKeyFloat64("CustomDoubleField")
	// CustomStringFieldは、カスタム文字列検索属性のキーです。
	CustomStringField = temporal.NewSearchAttributeKeyString("CustomStringField")
	// CustomDatetimeFieldは、カスタム日時検索属性のキーです。
	CustomDatetimeField = temporal.NewSearchAttributeKeyTime("CustomDatetimeField")
)

// SearchAttributesWorkflow ワークフロー定義
func SearchAttributesWorkflow(ctx workflow.Context) error {
	logger := workflow.GetLogger(ctx)
	logger.Info("SearchAttributes ワークフローが開始されました")

	// ワークフロー開始時に提供された検索属性を取得します。
	searchAttributes := workflow.GetTypedSearchAttributes(ctx)
	currentIntValue, ok := searchAttributes.GetInt64(CustomIntKey)
	if !ok {
		return errors.New("CustomIntFieldが設定されていません")
	}
	logger.Info("現在の検索属性値。", "CustomIntField", currentIntValue)

	// 検索属性を更新します。

	// これはコマンドがサーバーに送信されないため、サーバー上の検索属性は永続化されませんが、
	// ローカルキャッシュは更新されます。
	err := workflow.UpsertTypedSearchAttributes(ctx,
		CustomIntKey.ValueSet(2),
		CustomKeyword.ValueSet("Update1"),
		CustomBool.ValueSet(true),
		CustomDouble.ValueSet(3.14),
		CustomDatetimeField.ValueSet(workflow.Now(ctx).UTC()),
		CustomStringField.ValueSet("文字列フィールドはテキスト用です。クエリ時には、部分一致のためにトークン化されます。"),
	)
	if err != nil {
		return err
	}

	// 上記の変更を含む現在の検索属性を表示します。
	searchAttributes = workflow.GetTypedSearchAttributes(ctx)
	err = printSearchAttributes(searchAttributes, logger)
	if err != nil {
		return err
	}

	// 検索属性を再度更新します。
	err = workflow.UpsertTypedSearchAttributes(ctx,
		CustomKeyword.ValueSet("Update2"),
		CustomIntKey.ValueUnset(),
	)
	if err != nil {
		return err
	}

	// 更新が検索で見えるようになるまで待ちます。
	err = workflow.Sleep(ctx, 1*time.Second)
	if err != nil {
		return err
	}

	// 現在の検索属性を表示します。
	searchAttributes = workflow.GetTypedSearchAttributes(ctx)
	err = printSearchAttributes(searchAttributes, logger)
	if err != nil {
		return err
	}

	logger.Info("ワークフローが完了しました。")
	return nil
}

func printSearchAttributes(searchAttributes temporal.SearchAttributes, logger log.Logger) error {
	//workflowcheck:ignore
	if searchAttributes.Size() == 0 {
		logger.Info("現在の検索属性は空です。")
		return nil
	}

	var builder strings.Builder
	//workflowcheck:ignore ログ記録の理由のみで反復します
	for k, v := range searchAttributes.GetUntypedValues() {
		builder.WriteString(fmt.Sprintf("%s=%v\n", k.GetName(), v))
	}
	logger.Info(fmt.Sprintf("現在の検索属性値:\n%s", builder.String()))
	return nil
}
