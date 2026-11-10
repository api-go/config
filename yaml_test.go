package config_test

import (
	"os"
	"testing"

	"github.com/api-go/config"
)

func TestForYaml_Support(t *testing.T) {
	os.Clearenv()
	// create a test.yaml without extension in the load call
	os.WriteFile("test3.yaml", []byte(`
name: yaml-config
sets: [4, 5, 6]
`), 0644)
	defer os.Remove("test3.yaml")

	var testConf RegressionConf
	if err := config.Load(&testConf, "test3"); err != nil {
		t.Errorf("Load failed: %v", err)
	}

	if testConf.Name != "yaml-config" {
		t.Errorf("Expected name yaml-config, got %v", testConf.Name)
	}
	if len(testConf.Sets) != 3 || testConf.Sets[0] != 4 {
		t.Errorf("Sets failed: %v", testConf.Sets)
	}
}
