// Code generated by protoc-gen-go. DO NOT EDIT.
// source: infinity.proto

package common

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"

import (
	context "golang.org/x/net/context"
	grpc "google.golang.org/grpc"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion2 // please upgrade the proto package

type DeleteRecordingResponse struct {
	Status               StatusCode `protobuf:"varint,1,opt,name=status,proto3,enum=common.StatusCode" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *DeleteRecordingResponse) Reset()         { *m = DeleteRecordingResponse{} }
func (m *DeleteRecordingResponse) String() string { return proto.CompactTextString(m) }
func (*DeleteRecordingResponse) ProtoMessage()    {}
func (*DeleteRecordingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_infinity_98a14005dffa65b2, []int{0}
}
func (m *DeleteRecordingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DeleteRecordingResponse.Unmarshal(m, b)
}
func (m *DeleteRecordingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DeleteRecordingResponse.Marshal(b, m, deterministic)
}
func (dst *DeleteRecordingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DeleteRecordingResponse.Merge(dst, src)
}
func (m *DeleteRecordingResponse) XXX_Size() int {
	return xxx_messageInfo_DeleteRecordingResponse.Size(m)
}
func (m *DeleteRecordingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DeleteRecordingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DeleteRecordingResponse proto.InternalMessageInfo

func (m *DeleteRecordingResponse) GetStatus() StatusCode {
	if m != nil {
		return m.Status
	}
	return StatusCode_Ok
}

type UpsertRecordingResponse struct {
	Status               StatusCode `protobuf:"varint,1,opt,name=status,proto3,enum=common.StatusCode" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *UpsertRecordingResponse) Reset()         { *m = UpsertRecordingResponse{} }
func (m *UpsertRecordingResponse) String() string { return proto.CompactTextString(m) }
func (*UpsertRecordingResponse) ProtoMessage()    {}
func (*UpsertRecordingResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_infinity_98a14005dffa65b2, []int{1}
}
func (m *UpsertRecordingResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpsertRecordingResponse.Unmarshal(m, b)
}
func (m *UpsertRecordingResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpsertRecordingResponse.Marshal(b, m, deterministic)
}
func (dst *UpsertRecordingResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpsertRecordingResponse.Merge(dst, src)
}
func (m *UpsertRecordingResponse) XXX_Size() int {
	return xxx_messageInfo_UpsertRecordingResponse.Size(m)
}
func (m *UpsertRecordingResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_UpsertRecordingResponse.DiscardUnknown(m)
}

var xxx_messageInfo_UpsertRecordingResponse proto.InternalMessageInfo

func (m *UpsertRecordingResponse) GetStatus() StatusCode {
	if m != nil {
		return m.Status
	}
	return StatusCode_Ok
}

type Token struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *Token) Reset()         { *m = Token{} }
func (m *Token) String() string { return proto.CompactTextString(m) }
func (*Token) ProtoMessage()    {}
func (*Token) Descriptor() ([]byte, []int) {
	return fileDescriptor_infinity_98a14005dffa65b2, []int{2}
}
func (m *Token) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Token.Unmarshal(m, b)
}
func (m *Token) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Token.Marshal(b, m, deterministic)
}
func (dst *Token) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Token.Merge(dst, src)
}
func (m *Token) XXX_Size() int {
	return xxx_messageInfo_Token.Size(m)
}
func (m *Token) XXX_DiscardUnknown() {
	xxx_messageInfo_Token.DiscardUnknown(m)
}

var xxx_messageInfo_Token proto.InternalMessageInfo

func (m *Token) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

type IncDeviceCountResponse struct {
	Status               StatusCode `protobuf:"varint,1,opt,name=status,proto3,enum=common.StatusCode" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *IncDeviceCountResponse) Reset()         { *m = IncDeviceCountResponse{} }
func (m *IncDeviceCountResponse) String() string { return proto.CompactTextString(m) }
func (*IncDeviceCountResponse) ProtoMessage()    {}
func (*IncDeviceCountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_infinity_98a14005dffa65b2, []int{3}
}
func (m *IncDeviceCountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_IncDeviceCountResponse.Unmarshal(m, b)
}
func (m *IncDeviceCountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_IncDeviceCountResponse.Marshal(b, m, deterministic)
}
func (dst *IncDeviceCountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_IncDeviceCountResponse.Merge(dst, src)
}
func (m *IncDeviceCountResponse) XXX_Size() int {
	return xxx_messageInfo_IncDeviceCountResponse.Size(m)
}
func (m *IncDeviceCountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_IncDeviceCountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_IncDeviceCountResponse proto.InternalMessageInfo

func (m *IncDeviceCountResponse) GetStatus() StatusCode {
	if m != nil {
		return m.Status
	}
	return StatusCode_Ok
}

type DecDeviceCountResponse struct {
	Status               StatusCode `protobuf:"varint,1,opt,name=status,proto3,enum=common.StatusCode" json:"status,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *DecDeviceCountResponse) Reset()         { *m = DecDeviceCountResponse{} }
func (m *DecDeviceCountResponse) String() string { return proto.CompactTextString(m) }
func (*DecDeviceCountResponse) ProtoMessage()    {}
func (*DecDeviceCountResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_infinity_98a14005dffa65b2, []int{4}
}
func (m *DecDeviceCountResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DecDeviceCountResponse.Unmarshal(m, b)
}
func (m *DecDeviceCountResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DecDeviceCountResponse.Marshal(b, m, deterministic)
}
func (dst *DecDeviceCountResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DecDeviceCountResponse.Merge(dst, src)
}
func (m *DecDeviceCountResponse) XXX_Size() int {
	return xxx_messageInfo_DecDeviceCountResponse.Size(m)
}
func (m *DecDeviceCountResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DecDeviceCountResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DecDeviceCountResponse proto.InternalMessageInfo

func (m *DecDeviceCountResponse) GetStatus() StatusCode {
	if m != nil {
		return m.Status
	}
	return StatusCode_Ok
}

type RemainingBandwidthResponse struct {
	Status               StatusCode `protobuf:"varint,1,opt,name=status,proto3,enum=common.StatusCode" json:"status,omitempty"`
	Remaining            uint64     `protobuf:"varint,2,opt,name=Remaining,proto3" json:"Remaining,omitempty"`
	XXX_NoUnkeyedLiteral struct{}   `json:"-"`
	XXX_unrecognized     []byte     `json:"-"`
	XXX_sizecache        int32      `json:"-"`
}

func (m *RemainingBandwidthResponse) Reset()         { *m = RemainingBandwidthResponse{} }
func (m *RemainingBandwidthResponse) String() string { return proto.CompactTextString(m) }
func (*RemainingBandwidthResponse) ProtoMessage()    {}
func (*RemainingBandwidthResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_infinity_98a14005dffa65b2, []int{5}
}
func (m *RemainingBandwidthResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_RemainingBandwidthResponse.Unmarshal(m, b)
}
func (m *RemainingBandwidthResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_RemainingBandwidthResponse.Marshal(b, m, deterministic)
}
func (dst *RemainingBandwidthResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_RemainingBandwidthResponse.Merge(dst, src)
}
func (m *RemainingBandwidthResponse) XXX_Size() int {
	return xxx_messageInfo_RemainingBandwidthResponse.Size(m)
}
func (m *RemainingBandwidthResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_RemainingBandwidthResponse.DiscardUnknown(m)
}

var xxx_messageInfo_RemainingBandwidthResponse proto.InternalMessageInfo

func (m *RemainingBandwidthResponse) GetStatus() StatusCode {
	if m != nil {
		return m.Status
	}
	return StatusCode_Ok
}

func (m *RemainingBandwidthResponse) GetRemaining() uint64 {
	if m != nil {
		return m.Remaining
	}
	return 0
}

type ConsumedBandwidthRequest struct {
	Token                string   `protobuf:"bytes,1,opt,name=token,proto3" json:"token,omitempty"`
	RecordingId          uint64   `protobuf:"varint,2,opt,name=recordingId,proto3" json:"recordingId,omitempty"`
	Amount               uint64   `protobuf:"varint,3,opt,name=amount,proto3" json:"amount,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ConsumedBandwidthRequest) Reset()         { *m = ConsumedBandwidthRequest{} }
func (m *ConsumedBandwidthRequest) String() string { return proto.CompactTextString(m) }
func (*ConsumedBandwidthRequest) ProtoMessage()    {}
func (*ConsumedBandwidthRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_infinity_98a14005dffa65b2, []int{6}
}
func (m *ConsumedBandwidthRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ConsumedBandwidthRequest.Unmarshal(m, b)
}
func (m *ConsumedBandwidthRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ConsumedBandwidthRequest.Marshal(b, m, deterministic)
}
func (dst *ConsumedBandwidthRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ConsumedBandwidthRequest.Merge(dst, src)
}
func (m *ConsumedBandwidthRequest) XXX_Size() int {
	return xxx_messageInfo_ConsumedBandwidthRequest.Size(m)
}
func (m *ConsumedBandwidthRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_ConsumedBandwidthRequest.DiscardUnknown(m)
}

var xxx_messageInfo_ConsumedBandwidthRequest proto.InternalMessageInfo

func (m *ConsumedBandwidthRequest) GetToken() string {
	if m != nil {
		return m.Token
	}
	return ""
}

func (m *ConsumedBandwidthRequest) GetRecordingId() uint64 {
	if m != nil {
		return m.RecordingId
	}
	return 0
}

func (m *ConsumedBandwidthRequest) GetAmount() uint64 {
	if m != nil {
		return m.Amount
	}
	return 0
}

func init() {
	proto.RegisterType((*DeleteRecordingResponse)(nil), "common.DeleteRecordingResponse")
	proto.RegisterType((*UpsertRecordingResponse)(nil), "common.UpsertRecordingResponse")
	proto.RegisterType((*Token)(nil), "common.Token")
	proto.RegisterType((*IncDeviceCountResponse)(nil), "common.IncDeviceCountResponse")
	proto.RegisterType((*DecDeviceCountResponse)(nil), "common.DecDeviceCountResponse")
	proto.RegisterType((*RemainingBandwidthResponse)(nil), "common.RemainingBandwidthResponse")
	proto.RegisterType((*ConsumedBandwidthRequest)(nil), "common.ConsumedBandwidthRequest")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// BandwidthTrackingClient is the client API for BandwidthTracking service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type BandwidthTrackingClient interface {
	GetRemainingBandwidth(ctx context.Context, in *Token, opts ...grpc.CallOption) (*RemainingBandwidthResponse, error)
	AddConsumedBandwidth(ctx context.Context, in *ConsumedBandwidthRequest, opts ...grpc.CallOption) (*RemainingBandwidthResponse, error)
}

type bandwidthTrackingClient struct {
	cc *grpc.ClientConn
}

func NewBandwidthTrackingClient(cc *grpc.ClientConn) BandwidthTrackingClient {
	return &bandwidthTrackingClient{cc}
}

func (c *bandwidthTrackingClient) GetRemainingBandwidth(ctx context.Context, in *Token, opts ...grpc.CallOption) (*RemainingBandwidthResponse, error) {
	out := new(RemainingBandwidthResponse)
	err := c.cc.Invoke(ctx, "/common.BandwidthTracking/GetRemainingBandwidth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *bandwidthTrackingClient) AddConsumedBandwidth(ctx context.Context, in *ConsumedBandwidthRequest, opts ...grpc.CallOption) (*RemainingBandwidthResponse, error) {
	out := new(RemainingBandwidthResponse)
	err := c.cc.Invoke(ctx, "/common.BandwidthTracking/AddConsumedBandwidth", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BandwidthTrackingServer is the server API for BandwidthTracking service.
type BandwidthTrackingServer interface {
	GetRemainingBandwidth(context.Context, *Token) (*RemainingBandwidthResponse, error)
	AddConsumedBandwidth(context.Context, *ConsumedBandwidthRequest) (*RemainingBandwidthResponse, error)
}

func RegisterBandwidthTrackingServer(s *grpc.Server, srv BandwidthTrackingServer) {
	s.RegisterService(&_BandwidthTracking_serviceDesc, srv)
}

func _BandwidthTracking_GetRemainingBandwidth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BandwidthTrackingServer).GetRemainingBandwidth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.BandwidthTracking/GetRemainingBandwidth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BandwidthTrackingServer).GetRemainingBandwidth(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

func _BandwidthTracking_AddConsumedBandwidth_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ConsumedBandwidthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BandwidthTrackingServer).AddConsumedBandwidth(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.BandwidthTracking/AddConsumedBandwidth",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BandwidthTrackingServer).AddConsumedBandwidth(ctx, req.(*ConsumedBandwidthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _BandwidthTracking_serviceDesc = grpc.ServiceDesc{
	ServiceName: "common.BandwidthTracking",
	HandlerType: (*BandwidthTrackingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetRemainingBandwidth",
			Handler:    _BandwidthTracking_GetRemainingBandwidth_Handler,
		},
		{
			MethodName: "AddConsumedBandwidth",
			Handler:    _BandwidthTracking_AddConsumedBandwidth_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infinity.proto",
}

// DeviceTrackingClient is the client API for DeviceTracking service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type DeviceTrackingClient interface {
	IncDeviceCount(ctx context.Context, in *Token, opts ...grpc.CallOption) (*IncDeviceCountResponse, error)
	DecDeviceCount(ctx context.Context, in *Token, opts ...grpc.CallOption) (*DecDeviceCountResponse, error)
}

type deviceTrackingClient struct {
	cc *grpc.ClientConn
}

func NewDeviceTrackingClient(cc *grpc.ClientConn) DeviceTrackingClient {
	return &deviceTrackingClient{cc}
}

func (c *deviceTrackingClient) IncDeviceCount(ctx context.Context, in *Token, opts ...grpc.CallOption) (*IncDeviceCountResponse, error) {
	out := new(IncDeviceCountResponse)
	err := c.cc.Invoke(ctx, "/common.DeviceTracking/IncDeviceCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *deviceTrackingClient) DecDeviceCount(ctx context.Context, in *Token, opts ...grpc.CallOption) (*DecDeviceCountResponse, error) {
	out := new(DecDeviceCountResponse)
	err := c.cc.Invoke(ctx, "/common.DeviceTracking/DecDeviceCount", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// DeviceTrackingServer is the server API for DeviceTracking service.
type DeviceTrackingServer interface {
	IncDeviceCount(context.Context, *Token) (*IncDeviceCountResponse, error)
	DecDeviceCount(context.Context, *Token) (*DecDeviceCountResponse, error)
}

func RegisterDeviceTrackingServer(s *grpc.Server, srv DeviceTrackingServer) {
	s.RegisterService(&_DeviceTracking_serviceDesc, srv)
}

func _DeviceTracking_IncDeviceCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceTrackingServer).IncDeviceCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.DeviceTracking/IncDeviceCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceTrackingServer).IncDeviceCount(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

func _DeviceTracking_DecDeviceCount_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Token)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(DeviceTrackingServer).DecDeviceCount(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.DeviceTracking/DecDeviceCount",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(DeviceTrackingServer).DecDeviceCount(ctx, req.(*Token))
	}
	return interceptor(ctx, in, info, handler)
}

var _DeviceTracking_serviceDesc = grpc.ServiceDesc{
	ServiceName: "common.DeviceTracking",
	HandlerType: (*DeviceTrackingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "IncDeviceCount",
			Handler:    _DeviceTracking_IncDeviceCount_Handler,
		},
		{
			MethodName: "DecDeviceCount",
			Handler:    _DeviceTracking_DecDeviceCount_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infinity.proto",
}

// ContentClient is the client API for Content service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type ContentClient interface {
	DeleteRecording(ctx context.Context, in *RecordingIdentifier, opts ...grpc.CallOption) (*DeleteRecordingResponse, error)
	UpsertRecording(ctx context.Context, in *RecordingIdentifier, opts ...grpc.CallOption) (*UpsertRecordingResponse, error)
}

type contentClient struct {
	cc *grpc.ClientConn
}

func NewContentClient(cc *grpc.ClientConn) ContentClient {
	return &contentClient{cc}
}

func (c *contentClient) DeleteRecording(ctx context.Context, in *RecordingIdentifier, opts ...grpc.CallOption) (*DeleteRecordingResponse, error) {
	out := new(DeleteRecordingResponse)
	err := c.cc.Invoke(ctx, "/common.Content/DeleteRecording", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *contentClient) UpsertRecording(ctx context.Context, in *RecordingIdentifier, opts ...grpc.CallOption) (*UpsertRecordingResponse, error) {
	out := new(UpsertRecordingResponse)
	err := c.cc.Invoke(ctx, "/common.Content/UpsertRecording", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ContentServer is the server API for Content service.
type ContentServer interface {
	DeleteRecording(context.Context, *RecordingIdentifier) (*DeleteRecordingResponse, error)
	UpsertRecording(context.Context, *RecordingIdentifier) (*UpsertRecordingResponse, error)
}

func RegisterContentServer(s *grpc.Server, srv ContentServer) {
	s.RegisterService(&_Content_serviceDesc, srv)
}

func _Content_DeleteRecording_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordingIdentifier)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServer).DeleteRecording(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.Content/DeleteRecording",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServer).DeleteRecording(ctx, req.(*RecordingIdentifier))
	}
	return interceptor(ctx, in, info, handler)
}

func _Content_UpsertRecording_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordingIdentifier)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ContentServer).UpsertRecording(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/common.Content/UpsertRecording",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ContentServer).UpsertRecording(ctx, req.(*RecordingIdentifier))
	}
	return interceptor(ctx, in, info, handler)
}

var _Content_serviceDesc = grpc.ServiceDesc{
	ServiceName: "common.Content",
	HandlerType: (*ContentServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "DeleteRecording",
			Handler:    _Content_DeleteRecording_Handler,
		},
		{
			MethodName: "UpsertRecording",
			Handler:    _Content_UpsertRecording_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "infinity.proto",
}

func init() { proto.RegisterFile("infinity.proto", fileDescriptor_infinity_98a14005dffa65b2) }

var fileDescriptor_infinity_98a14005dffa65b2 = []byte{
	// 387 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xa4, 0x94, 0xcd, 0x6e, 0x9b, 0x40,
	0x10, 0xc7, 0x45, 0x5b, 0x53, 0x79, 0xda, 0x62, 0x75, 0xe5, 0xba, 0x88, 0x7e, 0x21, 0x4e, 0x56,
	0x0f, 0x3e, 0xd0, 0x07, 0xb0, 0x5a, 0xa8, 0x2a, 0xf7, 0x62, 0x89, 0x38, 0xb7, 0x5c, 0x08, 0x3b,
	0x76, 0x36, 0x0e, 0xb3, 0x0e, 0x2c, 0x89, 0xf2, 0x1a, 0x79, 0x8d, 0x3c, 0x43, 0xde, 0x2d, 0xe2,
	0xd3, 0x8e, 0x3f, 0x22, 0xcb, 0x3e, 0xc1, 0x0c, 0xcc, 0x6f, 0xf6, 0x3f, 0x1f, 0x0b, 0x86, 0xa0,
	0xa9, 0x20, 0xa1, 0xee, 0x06, 0x8b, 0x44, 0x2a, 0xc9, 0xf4, 0x48, 0xc6, 0xb1, 0x24, 0xeb, 0x7d,
	0xf9, 0x2c, 0xbd, 0xce, 0x5f, 0xf8, 0xec, 0xe3, 0x15, 0x2a, 0x0c, 0x30, 0x92, 0x09, 0x17, 0x34,
	0x0b, 0x30, 0x5d, 0x48, 0x4a, 0x91, 0xfd, 0x04, 0x3d, 0x55, 0xa1, 0xca, 0x52, 0x53, 0xb3, 0xb5,
	0xbe, 0xe1, 0xb2, 0x41, 0x15, 0x79, 0x52, 0x78, 0x3d, 0xc9, 0x31, 0xa8, 0xfe, 0xc8, 0x31, 0xa7,
	0x8b, 0x14, 0x13, 0x75, 0x1c, 0xe6, 0x1b, 0xb4, 0x26, 0x72, 0x8e, 0xc4, 0xba, 0xd0, 0x52, 0xf9,
	0x4b, 0x11, 0xd3, 0x0e, 0x4a, 0xc3, 0xf1, 0xa1, 0x37, 0xa2, 0xc8, 0xc7, 0x1b, 0x11, 0xa1, 0x27,
	0x33, 0x52, 0x07, 0x25, 0xf1, 0xa1, 0xe7, 0xe3, 0xd1, 0x94, 0x29, 0x58, 0x01, 0xc6, 0xa1, 0x20,
	0x41, 0xb3, 0x3f, 0x21, 0xf1, 0x5b, 0xc1, 0xd5, 0xc5, 0x21, 0x24, 0xf6, 0x15, 0xda, 0x0d, 0xc9,
	0x7c, 0x65, 0x6b, 0xfd, 0x37, 0xc1, 0xd2, 0xe1, 0x5c, 0x82, 0xe9, 0x49, 0x4a, 0xb3, 0x18, 0xf9,
	0x4a, 0x9a, 0xeb, 0x0c, 0x53, 0xb5, 0xbd, 0x4a, 0xcc, 0x86, 0x77, 0x49, 0xdd, 0x85, 0x11, 0xaf,
	0x88, 0xab, 0x2e, 0xd6, 0x03, 0x3d, 0x8c, 0x73, 0xe5, 0xe6, 0xeb, 0xe2, 0x63, 0x65, 0xb9, 0x8f,
	0x1a, 0x7c, 0x6c, 0x92, 0x4c, 0x92, 0x30, 0x9a, 0x0b, 0x9a, 0xb1, 0xff, 0xf0, 0xe9, 0x1f, 0xaa,
	0x4d, 0xb1, 0xec, 0x43, 0x2d, 0xaa, 0xe8, 0x99, 0xe5, 0xd4, 0xe6, 0x0b, 0x75, 0x39, 0x83, 0xee,
	0x6f, 0xce, 0x37, 0x04, 0x31, 0xbb, 0x8e, 0xdd, 0xa5, 0x75, 0x1f, 0xba, 0x7b, 0xaf, 0x81, 0x51,
	0xf6, 0xb5, 0x39, 0xfc, 0x10, 0x8c, 0xe7, 0x23, 0xb3, 0x7e, 0xea, 0xef, 0xb5, 0xb9, 0x63, 0xb2,
	0x86, 0x39, 0x72, 0x2f, 0xc0, 0xf6, 0xa1, 0x72, 0x1f, 0x34, 0x78, 0xeb, 0x49, 0x52, 0x48, 0x8a,
	0x8d, 0xa1, 0xb3, 0xb6, 0x6d, 0xec, 0xcb, 0x52, 0x57, 0xd3, 0x20, 0x24, 0x25, 0xa6, 0x02, 0x13,
	0xeb, 0xc7, 0x92, 0xbd, 0x7d, 0x47, 0xc7, 0xd0, 0x59, 0xdb, 0xbb, 0x3d, 0x81, 0x3b, 0xb6, 0xf5,
	0x5c, 0x2f, 0xae, 0x85, 0x5f, 0x4f, 0x01, 0x00, 0x00, 0xff, 0xff, 0xa7, 0x27, 0x97, 0x58, 0x3e,
	0x04, 0x00, 0x00,
}