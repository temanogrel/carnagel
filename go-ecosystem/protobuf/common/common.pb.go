// Code generated by protoc-gen-go. DO NOT EDIT.
// source: common.proto

package common

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

type StatusCode int32

const (
	StatusCode_Ok StatusCode = 0
	// Generic errors
	StatusCode_InvalidSessionToken StatusCode = 10
	StatusCode_RecordingNotFound   StatusCode = 11
	StatusCode_PlanNotFound        StatusCode = 12
	// Device tracking errors
	StatusCode_TooManyActiveDevices StatusCode = 30
	// Delete recording errors
	StatusCode_RecordingHasViews                StatusCode = 50
	StatusCode_RecordingInPremiumUserCollection StatusCode = 51
	// Generic error incase we fucked up
	StatusCode_InternalServerErr StatusCode = 500
)

var StatusCode_name = map[int32]string{
	0:   "Ok",
	10:  "InvalidSessionToken",
	11:  "RecordingNotFound",
	12:  "PlanNotFound",
	30:  "TooManyActiveDevices",
	50:  "RecordingHasViews",
	51:  "RecordingInPremiumUserCollection",
	500: "InternalServerErr",
}
var StatusCode_value = map[string]int32{
	"Ok": 0,
	"InvalidSessionToken":              10,
	"RecordingNotFound":                11,
	"PlanNotFound":                     12,
	"TooManyActiveDevices":             30,
	"RecordingHasViews":                50,
	"RecordingInPremiumUserCollection": 51,
	"InternalServerErr":                500,
}

func (x StatusCode) String() string {
	return proto.EnumName(StatusCode_name, int32(x))
}
func (StatusCode) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_common_d1f4029bdc8239aa, []int{0}
}

type RecordingIdentifier struct {
	Id                   uint64   `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *RecordingIdentifier) Reset()         { *m = RecordingIdentifier{} }
func (m *RecordingIdentifier) String() string { return proto.CompactTextString(m) }
func (*RecordingIdentifier) ProtoMessage()    {}
func (*RecordingIdentifier) Descriptor() ([]byte, []int) {
	return fileDescriptor_common_d1f4029bdc8239aa, []int{0}
}
func (m *RecordingIdentifier) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RecordingIdentifier.Unmarshal(m, b)
}
func (m *RecordingIdentifier) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RecordingIdentifier.Marshal(b, m, deterministic)
}
func (dst *RecordingIdentifier) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RecordingIdentifier.Merge(dst, src)
}
func (m *RecordingIdentifier) XXX_Size() int {
	return xxx_messageInfo_RecordingIdentifier.Size(m)
}
func (m *RecordingIdentifier) XXX_DiscardUnknown() {
	xxx_messageInfo_RecordingIdentifier.DiscardUnknown(m)
}

var xxx_messageInfo_RecordingIdentifier proto.InternalMessageInfo

func (m *RecordingIdentifier) GetId() uint64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func init() {
	proto.RegisterType((*RecordingIdentifier)(nil), "common.RecordingIdentifier")
	proto.RegisterEnum("common.StatusCode", StatusCode_name, StatusCode_value)
}

func init() { proto.RegisterFile("common.proto", fileDescriptor_common_d1f4029bdc8239aa) }

var fileDescriptor_common_d1f4029bdc8239aa = []byte{
	// 239 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x54, 0xcf, 0x4d, 0x4a, 0x03, 0x41,
	0x10, 0x05, 0x60, 0x33, 0xca, 0x2c, 0xca, 0x41, 0x3a, 0x1d, 0x7f, 0xb2, 0x92, 0x20, 0x0a, 0xe2,
	0xc2, 0x85, 0x39, 0x81, 0x44, 0xc5, 0x59, 0xa8, 0x21, 0x13, 0xdd, 0xb7, 0xdd, 0xa5, 0x14, 0xe9,
	0xa9, 0x92, 0xea, 0x9e, 0x11, 0x0f, 0xe8, 0x8d, 0x3c, 0x80, 0xa0, 0x10, 0x71, 0xf9, 0xde, 0xe2,
	0xe3, 0x3d, 0xa8, 0xbc, 0xb4, 0xad, 0xf0, 0xf9, 0x9b, 0x4a, 0x16, 0x5b, 0xfe, 0xa6, 0xa3, 0x13,
	0x18, 0x2d, 0xd0, 0x8b, 0x06, 0xe2, 0xd7, 0x3a, 0x20, 0x67, 0x7a, 0x21, 0x54, 0xbb, 0x03, 0x05,
	0x85, 0xf1, 0x60, 0x32, 0x38, 0xdd, 0x5a, 0x14, 0x14, 0xce, 0x3e, 0x07, 0x00, 0x4d, 0x76, 0xb9,
	0x4b, 0x33, 0x09, 0x68, 0x4b, 0x28, 0x1e, 0x56, 0x66, 0xc3, 0x1e, 0xc0, 0xa8, 0xe6, 0xde, 0x45,
	0x0a, 0x0d, 0xa6, 0x44, 0xc2, 0x4b, 0x59, 0x21, 0x1b, 0xb0, 0x7b, 0x30, 0x5c, 0xb3, 0xf7, 0x92,
	0x6f, 0xa4, 0xe3, 0x60, 0xb6, 0xad, 0x81, 0x6a, 0x1e, 0x1d, 0xaf, 0x9b, 0xca, 0x8e, 0x61, 0x77,
	0x29, 0x72, 0xe7, 0xf8, 0xe3, 0xd2, 0x67, 0xea, 0xf1, 0x0a, 0x7b, 0xf2, 0x98, 0xcc, 0xe1, 0x3f,
	0xe2, 0xd6, 0xa5, 0x27, 0xc2, 0xf7, 0x64, 0x2e, 0xec, 0x31, 0x4c, 0xfe, 0x06, 0xf3, 0x5c, 0xb1,
	0xa5, 0xae, 0x7d, 0x4c, 0xa8, 0x33, 0x89, 0x11, 0x7d, 0x26, 0x61, 0x33, 0xb5, 0xfb, 0x30, 0xac,
	0x39, 0xa3, 0xb2, 0x8b, 0x0d, 0x6a, 0x8f, 0x7a, 0xad, 0x6a, 0xbe, 0x36, 0x9f, 0xcb, 0x9f, 0xf7,
	0xd3, 0xef, 0x00, 0x00, 0x00, 0xff, 0xff, 0x4a, 0x51, 0xe3, 0x23, 0x0d, 0x01, 0x00, 0x00,
}
