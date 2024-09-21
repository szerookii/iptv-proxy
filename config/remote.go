package config

import (
	"encoding/json"
	"fmt"
	"github.com/bytedance/sonic"
)

type RemoteType string

const (
	RemoteTypeXtream RemoteType = "xtream"
	RemoteTypeStb    RemoteType = "stb"
)

type RemoteData interface {
	UnmarshalJSON([]byte) error
}

type Remote struct {
	Type RemoteType `json:"type"`
	Data RemoteData `json:"data"`
}

type XtreamRemote struct {
	URL      string `json:"url"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type StbRemote struct {
	URL        string `json:"url"`
	MacAddress string `json:"mac_address"`
}

func (x *XtreamRemote) UnmarshalJSON(data []byte) error {
	type temp XtreamRemote
	return sonic.Unmarshal(data, (*temp)(x))
}

func (s *StbRemote) UnmarshalJSON(data []byte) error {
	type temp StbRemote
	return sonic.Unmarshal(data, (*temp)(s))
}

func (r *Remote) UnmarshalJSON(data []byte) error {
	var raw struct {
		Type RemoteType      `json:"type"`
		Data json.RawMessage `json:"data"`
	}

	if err := sonic.Unmarshal(data, &raw); err != nil {
		return err
	}

	r.Type = raw.Type

	switch raw.Type {
	case RemoteTypeXtream:
		var xtream XtreamRemote
		if err := sonic.Unmarshal(raw.Data, &xtream); err != nil {
			return err
		}
		r.Data = &xtream
	case RemoteTypeStb:
		var stb StbRemote
		if err := sonic.Unmarshal(raw.Data, &stb); err != nil {
			return err
		}
		r.Data = &stb
	default:
		return fmt.Errorf("unknown remote type: %s", raw.Type)
	}

	return nil
}
