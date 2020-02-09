package runtime

import (
	"encoding/json"
)

func MarshalRequest(args map[string]interface{}, v interface{}) error {
	buf, err := json.Marshal(args)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, &v)
}
