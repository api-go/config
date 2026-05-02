package config

import (
	"os"
	"testing"
)

func BenchmarkLoadConfig(b *testing.B) {
	os.WriteFile("bench.json", []byte(`{"name":"bench-test","sets":[1,2,3]}`), 0644)
	defer os.Remove("bench.json")
	
	var conf struct {
		Name string
		Sets []int
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = LoadConfig("bench", &conf)
	}
}

func BenchmarkApplyEnvOverrides(b *testing.B) {
	os.Setenv("APP_DB_HOST", "localhost")
	os.Setenv("APP_DB_PORT", "5432")
	os.Setenv("APP_SERVER_TIMEOUT", "30s")
	
	var conf struct {
		App struct {
			Db struct {
				Host string
				Port int
			}
			Server struct {
				Timeout Duration
			}
		}
	}
	
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		applyEnvOverrides(&conf)
	}
}
