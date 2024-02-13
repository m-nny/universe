package jsoncache

import (
	"encoding/json"
	"log"
	"os"
	"path"
)

const cacheDir = "./data/cache"

func SetValue[Data any](key string, data Data) error {
	dir := path.Dir(cacheDir)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	dataBuf, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path.Join(cacheDir, key+".json"), dataBuf, os.ModePerm)
}

func GetValue[Data any](key string) (Data, error) {
	var data Data
	dataBuf, err := os.ReadFile(path.Join(cacheDir, key+".json"))
	if err != nil {
		return data, err
	}

	if err := json.Unmarshal(dataBuf, &data); err != nil {
		return data, err
	}
	return data, nil
}

func CachedExec[Data any](key string, fn func() (Data, error)) (Data, error) {
	cachedVal, err := GetValue[Data](key)
	if err != nil {
		log.Printf("[cache] miss on key %s: %v", key, err)
	} else if err == nil {
		log.Printf("[cache] hit on key %s", key)
		return cachedVal, nil
	}
	val, err := fn()
	if err != nil {
		return val, err
	}
	if err := SetValue(key, val); err != nil {
		return val, err
	}
	log.Printf("[cache] saved for key %s", key)
	return val, nil
}
