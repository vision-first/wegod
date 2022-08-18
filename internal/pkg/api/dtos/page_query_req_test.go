package dtos

import (
	"encoding/json"
	"testing"
)

func TestSerializePageQueryReq(t *testing.T) {
	j := `{"limit":10, "offset":11}`
	var req PageBuddhaReq
	err := json.Unmarshal([]byte(j), &req)
	if err != nil {
		t.Fatal(err)
	}
	t.Log(req)
}
