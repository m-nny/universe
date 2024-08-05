package jsoncache

import (
	"encoding/json"
	"log/slog"
	"os"
	"path"
)

var cacheDir = func() string {
	dir := "./data/cache"
	// dir, _ = filepath.Abs(dir)
	return dir
}()

func cacheFile(key string) string {
	return path.Join(cacheDir, key+".json")
}

func setValue[Data any](file string, data Data) error {
	dir := path.Dir(file)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	dataBuf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(file, dataBuf, os.ModePerm)
}

func getValue[Data any](file string) (Data, error) {
	var data Data
	dataBuf, err := os.ReadFile(file)
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(dataBuf, &data); err != nil {
		return data, err
	}
	return data, nil
}

func CachedExec[Data any](key string, fn func() (Data, error)) (Data, error) {
	file := cacheFile(key)
	cachedVal, err := getValue[Data](file)
	if err != nil {
		slog.Debug("jsoncahe.CachedExec(): miss on", "key", key, "err", err)
	} else {
		slog.Debug("jsoncahe.CachedExec(): hit on key", "key", key)
		return cachedVal, nil
	}
	val, err := fn()
	if err != nil {
		return val, err
	}
	if err := setValue(file, val); err != nil {
		return val, err
	}
	slog.Debug("jsoncahe.CachedExec(): saved for", "key", key)
	return val, nil
}
