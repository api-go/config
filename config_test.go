package config

import (
	"os"
	"testing"
	"time"
)

type TestList struct {
	Name string
}

func (l *TestList) ConfigureBy(setting string) {
	l.Name = setting
}

type testConfType struct {
	Name     string
	Sets     []int
	List     map[string]*TestList
	Duration Duration
}

func TestConfig_ComplexStruct(t *testing.T) {
	os.Clearenv()
	// 测试复杂嵌套、Map、指针以及 Duration
	os.WriteFile("complex.json", []byte(`{
		"name": "test-config",
		"sets": [1, 2, 3],
		"list": {
			"aaa": {"name": "222"}
		},
		"duration": "100s"
	}`), 0644)
	defer os.Remove("complex.json")

	// 模拟 env 覆盖
	os.Setenv("LIST_BBB", "{\"name\":\"xxx\"}")

	var conf testConfType
	if err := Load("complex", &conf); err != nil {
		t.Errorf("Load failed: %v", err)
	}

	if conf.Name != "test-config" {
		t.Errorf("Name failed, got %s", conf.Name)
	}
	if len(conf.Sets) != 3 || conf.Sets[1] != 2 {
		t.Errorf("Sets failed, got %v", conf.Sets)
	}
	if conf.List == nil || conf.List["aaa"] == nil || conf.List["aaa"].Name != "222" {
		t.Errorf("Map aaa failed, got %v", conf.List["aaa"])
	}
	if conf.List == nil || conf.List["bbb"] == nil || conf.List["bbb"].Name != "xxx" {
		t.Errorf("Map bbb (env) failed, got %v", conf.List["bbb"])
	}
	if conf.Duration.TimeDuration() != 100*time.Second {
		t.Errorf("Duration failed, got %v", conf.Duration.TimeDuration())
	}
}
