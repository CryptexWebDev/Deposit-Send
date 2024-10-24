package uniclient

import (
	"encoding/json"
	"fmt"
)

const (
	JSON_RPC_VERSION_2_0 = "2.0"

	ERROR_CODE_PARSE_ERROR         = -32700
	ERROR_MESSAGE_PARSE_ERROR      = "Parse error"
	ERROR_CODE_INVALID_REQUEST     = -32600
	ERROR_MESSAGE_INVALID_REQUEST  = "invalid request"
	ERROR_CODE_METHOD_NOT_FOUND    = -32601
	ERROR_MESSAGE_METHOD_NOT_FOUND = "method not found"
	ERROR_CODE_SERVER_ERROR        = -32000
	ERROR_MESSAGE_SERVER_ERROR     = "server error"
)

func NewRequest(method string, params interface{}) (req *Request) {
	req = &Request{
		Id:      "1",
		JsonRpc: JSON_RPC_VERSION_2_0,
		Method:  method,
		Params:  params,
	}
	return
}

type Request struct {
	Id      RequestId   `json:"id"`
	JsonRpc string      `json:"jsonrpc"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params,omitempty"`
}

func (r *Request) SetId(id RequestId) {
	r.Id = id
}

func (r *Request) String() string {
	b, _ := json.MarshalIndent(r, "", " ")
	return string(b)
}

type RequestId string

func (id RequestId) String() string {
	return string(id)
}

func (id RequestId) MarshalJSON() ([]byte, error) {
	if id == "" {
		return []byte("0"), nil
	}
	out := fmt.Sprintf("%s", id)
	return []byte(out), nil
}

func (id *RequestId) UnmarshalJSON(data []byte) error {
	var dc []byte
	fc := true
	for _, b := range data {
		if b < 58 && b > 47 {
			if fc && b != 48 {
				fc = false
			}
			if !fc {
				dc = append(dc, b)
			}
		}
	}
	if len(dc) == 0 {
		*id = "0"
		return nil
	}
	*id = RequestId(dc)
	return nil
}
