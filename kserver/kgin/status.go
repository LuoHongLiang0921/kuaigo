package kgin

import (
	"encoding/json"
	"fmt"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

// EmptyMessage ...
type EmptyMessage struct{}

// Reset ...
func (m *EmptyMessage) Reset() { *m = EmptyMessage{} }

// String ...
func (m *EmptyMessage) String() string { return "{}" }

// ProtoMessage ...
func (*EmptyMessage) ProtoMessage() {}

// GRPCProxyMessage ...
type GRPCProxyMessage struct {
	Error   int           `protobuf:"varint,1,opt,name=error" json:"error"`
	Message string        `protobuf:"bytes,2,opt,name=msg" json:"msg"`
	Data    proto.Message `protobuf:"bytes,3,opt,name=data" json:"data"`
}

// Reset ...
func (m *GRPCProxyMessage) Reset() { *m = GRPCProxyMessage{} }

// String ...
func (m *GRPCProxyMessage) String() string { return proto.CompactTextString(m) }

// ProtoMessage ...
func (*GRPCProxyMessage) ProtoMessage() {}

// MarshalJSONPB ...
func (m *GRPCProxyMessage) MarshalJSONPB(jsb *jsonpb.Marshaler) ([]byte, error) {
	ss, err := jsonpbMarshaler.MarshalToString(m.Data)
	if err != nil {
		return []byte{}, err
	}

	msg := struct {
		Error   int             `json:"error"`
		Message string          `json:"msg"`
		Data    json.RawMessage `json:"data"`
	}{
		Error:   m.Error,
		Message: m.Message,
		Data:    json.RawMessage([]byte(ss)),
	}

	return json.Marshal(msg)
}

var (
	jsonpbMarshaler = jsonpb.Marshaler{
		EmitDefaults: true,
	}
)

func createStatusErr(code uint32, msg string) string {
	return fmt.Sprintf("%d:%s", code, msg)
}
