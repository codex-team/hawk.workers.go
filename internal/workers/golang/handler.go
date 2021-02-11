package golang

import "github.com/valyala/fastjson"

var event_type = fastjson.MustParse(`"errors/golang"`)

// Handler modifies incoming messages.
func Handler(v *fastjson.Value) ([]byte, error) {
	v.Set("eventType", event_type)
	v.Set("event", v.Get("payload"))
	v.Del("payload")

	return v.StringBytes()
}
