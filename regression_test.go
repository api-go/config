package config

import (
	"os"
	"testing"
	"time"
)

type RegressionConf struct {
	Name     string
	Sets     []int
	List     map[string]*TestList
	List2    []string
	Duration Duration
	Ext      map[string]interface{}
	ExtList  []interface{}
}

func TestForMap_Regression(t *testing.T) {
	os.Clearenv()
	os.WriteFile("test.json", []byte(`{"name":"test-config"}`), 0644)
	defer os.Remove("test.json")
	
	os.Setenv("TEST_LIST_CCC", "333")
	
	testConf := map[string]interface{}{}
	err := LoadConfig("test", &testConf)
	if err != nil {
		t.Fatalf("LoadConfig failed: %v", err)
	}
	if testConf["name"] != "test-config" {
		t.Errorf("Expected name test-config, got %v", testConf["name"])
	}
}

func TestForStruct_Regression(t *testing.T) {
	os.Clearenv()
	os.WriteFile("test.json", []byte(`{
		"name": "test-config",
		"sets": [1, 2, 3],
		"list": {"aaa": {"name": "222"}, "bbb": {"name": "xxx"}},
		"duration": "100s"
	}`), 0644)
	defer os.Remove("test.json")
	
	os.Setenv("LIST_CCC", "{\"name\":\"333\"}")
	os.Setenv("LIST_DDD", "{\"name\":\"444\"}")

	var testConf RegressionConf
	if err := LoadConfig("test", &testConf); err != nil {
		t.Errorf("LoadConfig failed: %v", err)
	}

	if testConf.Name != "test-config" || len(testConf.Sets) != 3 || testConf.Sets[1] != 2 {
		t.Errorf("Basic struct mapping failed: %+v", testConf)
	}
	if testConf.List["aaa"].Name != "222" || testConf.List["ccc"].Name != "333" {
		t.Errorf("List mapping failed: %+v", testConf.List)
	}
	if testConf.Duration.TimeDuration() != 100*time.Second {
		t.Errorf("Duration failed: %v", testConf.Duration.TimeDuration())
	}
}

func TestForYml_Regression(t *testing.T) {
	os.Clearenv()
	os.WriteFile("test2.yml", []byte(`
name: test-config
sets: [1, 2, 3]
list:
  aaa: {name: "222"}
  bbb: {name: "xxx"}
list2: ["a", "b"]
ext:
  bbb: "222"
  ccc: {val: "111"}
extList: ["a", "222", {val: "111"}]
`), 0644)
	defer os.Remove("test2.yml")

	var testConf RegressionConf
	if err := LoadConfig("test2", &testConf); err != nil {
		t.Errorf("LoadConfig failed: %v", err)
	}

	if testConf.Name != "test-config" || len(testConf.List2) != 2 {
		t.Errorf("YML structure failed: %+v", testConf)
	}
	if testConf.Ext["bbb"] != "222" {
		t.Errorf("Any map failed: %v", testConf.Ext["bbb"])
	}
	if len(testConf.ExtList) != 3 || testConf.ExtList[1] != "222" {
		t.Errorf("Any slice failed: %v", testConf.ExtList)
	}
}
