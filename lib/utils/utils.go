package utils

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/url"
	"time"

	"gitlab.com/snap-clickstaff/torque-go/config"
	"gitlab.com/snap-clickstaff/torque-go/lib/meta"
)

func JsonObjectLoadReader(io io.Reader) meta.O {
	ioBytes, err := ioutil.ReadAll(io)
	if err != nil {
		panic(err)
	}

	var jsonData meta.O
	err = json.Unmarshal(ioBytes, &jsonData)
	if err != nil {
		panic(err)
	}

	return jsonData
}

func JsonModelLoadReader(io io.Reader, model interface{}) interface{} {
	ioBytes, err := ioutil.ReadAll(io)
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(ioBytes, &model)
	if err != nil {
		panic(err)
	}

	return model
}

func MapGetDefault(dataMap map[string]interface{}, key string, default_ interface{}) interface{} {
	value, ok := dataMap[key]
	if !ok {
		return default_
	}

	return value
}

func QueryGetDefault(query url.Values, key string, default_ string) string {
	_, ok := query[key]
	if !ok {
		return default_
	}

	return query.Get(key)
}

func GetEnvCacheDuration(timeCache time.Duration) time.Duration {
	if config.Test {
		timeCache = timeCache / 5
	}
	return timeCache
}
