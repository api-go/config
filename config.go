package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"apigo.cc/go/cast"
	"apigo.cc/go/file"
)

type Duration time.Duration

func (d *Duration) UnmarshalJSON(value []byte) error {
	result, err := time.ParseDuration(string(bytes.Trim(value, "\"")))
	if err == nil {
		*d = Duration(result)
	}
	return err
}

func (d *Duration) UnmarshalYAML(value interface{}) error {
	result, err := time.ParseDuration(cast.String(value))
	if err == nil {
		*d = Duration(result)
	}
	return err
}

func (d Duration) MarshalJSON() ([]byte, error) {
	return []byte("\"" + time.Duration(d).String() + "\""), nil
}

func (d Duration) MarshalYAML() (interface{}, error) {
	return time.Duration(d).String(), nil
}

func (d Duration) TimeDuration() time.Duration {
	return time.Duration(d)
}

// Load loads configuration from a file, applies env.json/yaml and any
// command-line --config overlays in order, then overrides with environment
// variables.
func Load(conf interface{}, name string) error {
	var errs []error

	// 1. Load base file
	filename := resolveConfigFilePath(name)
	if filename != "" {
		if err := file.UnmarshalFile(filename, conf); err != nil {
			errs = append(errs, err)
		}
	}

	// 2. Load env.json/yaml if present
	envFile := resolveConfigFilePath("env")
	if envFile != "" {
		if err := loadOverlay(conf, name, envFile); err != nil {
			errs = append(errs, err)
		}
	}

	// 3. Apply explicit overlays in command-line order
	overlays, err := configFilesFromArgs(os.Args[1:])
	if err != nil {
		errs = append(errs, err)
	}
	for _, overlay := range overlays {
		filename := resolveConfigFilePath(overlay)
		if filename == "" {
			errs = append(errs, fmt.Errorf("config overlay not found: %s", overlay))
			continue
		}
		if err = loadOverlay(conf, name, filename); err != nil {
			errs = append(errs, err)
		}
	}

	// 4. Apply Environment Variable overrides
	applyEnvOverrides(conf)

	return errors.Join(errs...)
}

func loadOverlay(conf interface{}, name, filename string) error {
	overlay := map[string]interface{}{}
	if err := file.UnmarshalFile(filename, &overlay); err != nil {
		return err
	}
	if val, ok := overlay[name]; ok {
		cast.Convert(conf, val)
	} else {
		cast.Convert(conf, overlay)
	}
	return nil
}

func configFilesFromArgs(args []string) ([]string, error) {
	files := make([]string, 0, 1)
	var errs []error
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if strings.HasPrefix(arg, "--config=") {
			filename := strings.TrimSpace(strings.TrimPrefix(arg, "--config="))
			if filename == "" {
				errs = append(errs, errors.New("--config requires a file path"))
			} else {
				files = append(files, filename)
			}
			continue
		}
		if arg != "--config" {
			continue
		}
		if i+1 >= len(args) || strings.HasPrefix(args[i+1], "-") {
			errs = append(errs, errors.New("--config requires a file path"))
			continue
		}
		i++
		files = append(files, args[i])
	}
	return files, errors.Join(errs...)
}

// To is a frictionless version of Load. It returns the configuration object of type T.
// If loading fails, it returns the zero value of type T.
func To[T any](name string) T {
	var conf T
	_ = Load(&conf, name)
	return conf
}

func resolveConfigPath(name string) string {
	// If it already has an extension, check directly
	if strings.HasSuffix(name, ".json") || strings.HasSuffix(name, ".yml") || strings.HasSuffix(name, ".yaml") {
		if file.Exists(name) {
			return name
		}
		return ""
	}

	// Otherwise, check for .yml, .yaml, then .json
	if file.Exists(name + ".yml") {
		return name + ".yml"
	}
	if file.Exists(name + ".yaml") {
		return name + ".yaml"
	}
	if file.Exists(name + ".json") {
		return name + ".json"
	}
	return ""
}

func resolveConfigFilePath(name string) string {
	// Search in: current dir, executable dir, user home dir
	dirs := []string{"."}
	if execPath, err := os.Executable(); err == nil {
		dirs = append(dirs, filepath.Dir(execPath))
	}
	if home, err := os.UserHomeDir(); err == nil {
		dirs = append(dirs, home)
	}

	for _, d := range dirs {
		p := filepath.Join(d, name)
		if f := resolveConfigPath(p); f != "" {
			return f
		}
	}
	return ""
}

func applyEnvOverrides(conf interface{}) {
	envMap := make(map[string]interface{})
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		if len(pair) != 2 {
			continue
		}

		// Convert DB_HOST to {"db": {"host": value}}
		parts := strings.Split(strings.ToLower(pair[0]), "_")
		curr := envMap
		for i := 0; i < len(parts)-1; i++ {
			if _, ok := curr[parts[i]]; !ok {
				curr[parts[i]] = make(map[string]interface{})
			}
			if next, ok := curr[parts[i]].(map[string]interface{}); ok {
				curr = next
			} else {
				break
			}
		}
		curr[parts[len(parts)-1]] = pair[1]
	}

	cast.Convert(conf, envMap)
}
