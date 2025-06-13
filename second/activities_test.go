package iplocate_test

import (
	"io"
	"net/http"
	"strings"
	"testing"

	"temporal-ip-geolocation/iplocate"

	"github.com/stretchr/testify/assert"
	"go.temporal.io/sdk/testsuite"
)

type MockHTTPClient struct {
	Response *http.Response
	Err      error
}

func (m *MockHTTPClient) Get(url string) (*http.Response, error) {
	return m.Response, m.Err
}

// TestGetIPは、モックサーバーを使ってGetIPアクティビティをテストします。
func TestGetIP(t *testing.T) {
	// テスト環境をセットアップ
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()

	// フェイクIPアドレスを返すモックレスポンスを作成
	mockResponse := &http.Response{
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader("127.0.0.1\n")),
	}

	// Activitiesをロードし、モックレスポンスを注入
	ipActivities := &iplocate.IPActivities{
		HTTPClient: &MockHTTPClient{Response: mockResponse},
	}
	env.RegisterActivity(ipActivities)

	// GetIP関数を呼び出す
	val, err := env.ExecuteActivity(ipActivities.GetIP)
	if err != nil {
		t.Fatalf("エラーは発生しないはずですが、発生しました: %v", err)
	}

	// アクティビティの結果を取得
	var ip string
	val.Get(&ip)

	// 返されたIPを検証
	expectedIP := "127.0.0.1"
	assert.Equal(t, ip, expectedIP)
}

// TestGetLocationInfoは、モックサーバーを使ってGetLocationInfoアクティビティをテストします。
func TestGetLocationInfo(t *testing.T) {
	// テスト環境をセットアップ
	testSuite := &testsuite.WorkflowTestSuite{}
	env := testSuite.NewTestActivityEnvironment()

	mockResponse := &http.Response{
		StatusCode: 200,
		Body: io.NopCloser(strings.NewReader(`{
            "city": "San Francisco",
            "regionName": "California",
            "country": "United States"
        }`)),
	}

	ipActivities := &iplocate.IPActivities{
		HTTPClient: &MockHTTPClient{Response: mockResponse},
	}

	env.RegisterActivity(ipActivities)

	ip := "127.0.0.1"
	val, err := env.ExecuteActivity(ipActivities.GetLocationInfo, ip)
	if err != nil {
		t.Fatalf("エラーは発生しないはずですが、発生しました: %v", err)
	}

	var location string
	val.Get(&location)

	expectedLocation := "San Francisco, California, United States"
	assert.Equal(t, location, expectedLocation)
}
