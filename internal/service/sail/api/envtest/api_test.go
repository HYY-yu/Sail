package envtest

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApiEnvTest_StartAndYaml(t *testing.T) {
	ap := &ApiEnvTest{}
	err := ap.Start()
	if err != nil {
		t.Error(err)
		return
	}

	v, err := ap.MetaConfigYaml()
	if err != nil {
		t.Error(err)
		return
	}
	assert.NotEmpty(t, v)
	t.Log(v)
}
