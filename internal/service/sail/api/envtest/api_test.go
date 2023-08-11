package envtest

import "testing"

func TestApiEnvTest_Start(t *testing.T) {
	ap := &ApiEnvTest{}
	err := ap.Start()
	if err != nil {
		t.Error(err)
	}
}
