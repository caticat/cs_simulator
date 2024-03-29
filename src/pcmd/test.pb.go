// Code generated by protoc-gen-go. DO NOT EDIT.
// source: test.proto

/*
Package pcmd is a generated protocol buffer package.

It is generated from these files:
	test.proto

It has these top-level messages:
	ProtoCommandTest
*/
package pcmd

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type EProtoCommand int32

const (
	EProtoCommand_EProtoCommandBegin EProtoCommand = 0
	EProtoCommand_EProtoCommandTest  EProtoCommand = 1
)

var EProtoCommand_name = map[int32]string{
	0: "EProtoCommandBegin",
	1: "EProtoCommandTest",
}
var EProtoCommand_value = map[string]int32{
	"EProtoCommandBegin": 0,
	"EProtoCommandTest":  1,
}

func (x EProtoCommand) String() string {
	return proto.EnumName(EProtoCommand_name, int32(x))
}
func (EProtoCommand) EnumDescriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

type ProtoCommandTest struct {
	Id   uint32 `protobuf:"varint,1,opt,name=id" json:"id,omitempty"`
	Data string `protobuf:"bytes,2,opt,name=data" json:"data,omitempty"`
}

func (m *ProtoCommandTest) Reset()                    { *m = ProtoCommandTest{} }
func (m *ProtoCommandTest) String() string            { return proto.CompactTextString(m) }
func (*ProtoCommandTest) ProtoMessage()               {}
func (*ProtoCommandTest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *ProtoCommandTest) GetId() uint32 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *ProtoCommandTest) GetData() string {
	if m != nil {
		return m.Data
	}
	return ""
}

func init() {
	proto.RegisterType((*ProtoCommandTest)(nil), "pcmd.ProtoCommandTest")
	proto.RegisterEnum("pcmd.EProtoCommand", EProtoCommand_name, EProtoCommand_value)
}

func init() { proto.RegisterFile("test.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 127 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0xe2, 0x2a, 0x49, 0x2d, 0x2e,
	0xd1, 0x2b, 0x28, 0xca, 0x2f, 0xc9, 0x17, 0x62, 0x29, 0x48, 0xce, 0x4d, 0x51, 0x32, 0xe3, 0x12,
	0x08, 0x00, 0x71, 0x9d, 0xf3, 0x73, 0x73, 0x13, 0xf3, 0x52, 0x42, 0x52, 0x8b, 0x4b, 0x84, 0xf8,
	0xb8, 0x98, 0x32, 0x53, 0x24, 0x18, 0x15, 0x18, 0x35, 0x78, 0x83, 0x98, 0x32, 0x53, 0x84, 0x84,
	0xb8, 0x58, 0x52, 0x12, 0x4b, 0x12, 0x25, 0x98, 0x14, 0x18, 0x35, 0x38, 0x83, 0xc0, 0x6c, 0x2d,
	0x3b, 0x2e, 0x5e, 0x57, 0x64, 0x8d, 0x42, 0x62, 0x5c, 0x42, 0x28, 0x02, 0x4e, 0xa9, 0xe9, 0x99,
	0x79, 0x02, 0x0c, 0x42, 0xa2, 0x5c, 0x82, 0xae, 0xe8, 0x36, 0x08, 0x30, 0x26, 0xb1, 0x81, 0x1d,
	0x61, 0x0c, 0x08, 0x00, 0x00, 0xff, 0xff, 0xf8, 0x91, 0x4d, 0xdb, 0x92, 0x00, 0x00, 0x00,
}
