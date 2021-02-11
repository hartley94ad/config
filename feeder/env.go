package feeder

import (
	"bufio"
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// Env is a feeder that feeds using a single env file.
type Env struct {
	Path string
}

// Feed returns all the content.
func (e *Env) Feed() (map[string]interface{}, error) {
	m := make(map[string]interface{})

	values, err := load(e.Path)
	if err != nil {
		return nil, err
	}

	for k, v := range values {
		m[standardize(k)] = get(k, v)
	}

	return m, nil
}

// standardize updates config key (e.g. APP_NAME to  app.name)
func standardize(k string) string {
	return strings.Replace(strings.ToLower(k), "_", ".", -1)
}

// load reads the given env file and extracts the variables as a string map
func load(filename string) (map[string]string, error) {
	path, err := filepath.Abs(filename)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filepath.Clean(path))
	if err != nil {
		return nil, err
	}
	defer file.Close()

	variables, err := read(file)
	if err != nil {
		return nil, err
	}

	return variables, nil
}

// read opens the given env file and extracts the variables as a string map
func read(file io.Reader) (map[string]string, error) {
	items := map[string]string{}
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		key, value, err := parse(scanner.Text())
		if err != nil {
			return nil, err
		}

		if key != "" {
			items[key] = value
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return items, nil
}

// parse extracts the key/value from the given line
func parse(line string) (string, string, error) {
	ln := strings.TrimSpace(line)

	if len(ln) == 0 {
		return "", "", nil
	}

	if ln[0] == '#' {
		return "", "", nil
	}

	s := strings.Index(ln, "=")
	if s == -1 {
		return "", "", errors.New("Invalid line: " + ln)
	}

	k := strings.TrimSpace(ln[:s])
	v := strings.TrimSpace(ln[s+1:])

	return k, v, nil
}

// get fetches variable from OS if exist and return fallback otherwise.
func get(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}

	return fallback
}
