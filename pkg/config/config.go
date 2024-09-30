package config

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"

	utilyaml "k8s.io/apimachinery/pkg/util/yaml"
)

// Load loads the given path into the context.
func Load(ctx context.Context, path ...string) ([]any, error) {
	var raws []json.RawMessage

	for _, p := range path {
		r, err := loadRawConfig(p)
		if err != nil {
			return nil, err
		}
		raws = append(raws, r...)
	}

	var meta any
	objs := make([]any, 0, len(raws))
	for _, raw := range raws {
		err := json.Unmarshal(raw, &meta)
		if err != nil {
			slog.Error("Unmarshal faild", slog.String("raw", string(raw)), slog.String("err", err.Error()))
			continue
		}
		objs = append(objs, meta)
	}

	return objs, nil
}

func loadRawConfig(p string) ([]json.RawMessage, error) {
	var raws []json.RawMessage
	file, err := os.Open(p)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = file.Close()
	}()
	decoder := utilyaml.NewYAMLToJSONDecoder(file)
	for {
		var raw json.RawMessage
		err = decoder.Decode(&raw)
		if err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, err
		}
		raws = append(raws, raw)
	}
	return raws, nil
}

type MetaObject interface {
	GetName() string
}
