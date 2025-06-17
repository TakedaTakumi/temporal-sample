package typedsearchattributes

import (
	"testing"

	"github.com/stretchr/testify/require"
	"go.temporal.io/sdk/temporal"
	"go.temporal.io/sdk/testsuite"
)

func Test_Workflow(t *testing.T) {
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestWorkflowEnvironment()

	// 開始時の検索属性をモック
	_ = env.SetTypedSearchAttributesOnStart(temporal.NewSearchAttributes(CustomIntKey.ValueSet(1)))

	// upsert操作をモック
	env.OnUpsertTypedSearchAttributes(
		temporal.NewSearchAttributes(
			CustomIntKey.ValueSet(2),
			CustomKeyword.ValueSet("Update1"),
			CustomBool.ValueSet(true),
			CustomDouble.ValueSet(3.14),
			CustomDatetimeField.ValueSet(env.Now().UTC()),
			CustomStringField.ValueSet("文字列フィールドはテキスト用です。クエリ時には、部分一致のためにトークン化されます。"),
		)).Return(nil).Once()

	env.OnUpsertTypedSearchAttributes(
		temporal.NewSearchAttributes(
			CustomKeyword.ValueSet("Update2"),
			CustomIntKey.ValueUnset(),
		)).Return(nil).Once()

	env.ExecuteWorkflow(SearchAttributesWorkflow)
	require.True(t, env.IsWorkflowCompleted())
	require.NoError(t, env.GetWorkflowError())
}
