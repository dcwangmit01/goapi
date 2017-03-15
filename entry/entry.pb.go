// Code generated by protoc-gen-go.
// source: entry/entry.proto
// DO NOT EDIT!

/*
Package entry is a generated protocol buffer package.

It is generated from these files:
	entry/entry.proto

It has these top-level messages:
	EchoMessage
*/
package entry

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"

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

type EchoMessage struct {
	Value string `protobuf:"bytes,1,opt,name=value" json:"value,omitempty"`
}

func (m *EchoMessage) Reset()                    { *m = EchoMessage{} }
func (m *EchoMessage) String() string            { return proto.CompactTextString(m) }
func (*EchoMessage) ProtoMessage()               {}
func (*EchoMessage) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *EchoMessage) GetValue() string {
	if m != nil {
		return m.Value
	}
	return ""
}

func init() {
	proto.RegisterType((*EchoMessage)(nil), "entry.EchoMessage")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for Entry service

type EntryClient interface {
	Echo(ctx context.Context, in *EchoMessage, opts ...grpc.CallOption) (*EchoMessage, error)
}

type entryClient struct {
	cc *grpc.ClientConn
}

func NewEntryClient(cc *grpc.ClientConn) EntryClient {
	return &entryClient{cc}
}

func (c *entryClient) Echo(ctx context.Context, in *EchoMessage, opts ...grpc.CallOption) (*EchoMessage, error) {
	out := new(EchoMessage)
	err := grpc.Invoke(ctx, "/entry.Entry/Echo", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for Entry service

type EntryServer interface {
	Echo(context.Context, *EchoMessage) (*EchoMessage, error)
}

func RegisterEntryServer(s *grpc.Server, srv EntryServer) {
	s.RegisterService(&_Entry_serviceDesc, srv)
}

func _Entry_Echo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(EchoMessage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(EntryServer).Echo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entry.Entry/Echo",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(EntryServer).Echo(ctx, req.(*EchoMessage))
	}
	return interceptor(ctx, in, info, handler)
}

var _Entry_serviceDesc = grpc.ServiceDesc{
	ServiceName: "entry.Entry",
	HandlerType: (*EntryServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Echo",
			Handler:    _Entry_Echo_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "entry/entry.proto",
}

func init() { proto.RegisterFile("entry/entry.proto", fileDescriptor0) }

var fileDescriptor0 = []byte{
	// 160 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0xe2, 0x12, 0x4c, 0xcd, 0x2b, 0x29,
	0xaa, 0xd4, 0x07, 0x93, 0x7a, 0x05, 0x45, 0xf9, 0x25, 0xf9, 0x42, 0xac, 0x60, 0x8e, 0x94, 0x4c,
	0x7a, 0x7e, 0x7e, 0x7a, 0x4e, 0xaa, 0x7e, 0x62, 0x41, 0xa6, 0x7e, 0x62, 0x5e, 0x5e, 0x7e, 0x49,
	0x62, 0x49, 0x66, 0x7e, 0x5e, 0x31, 0x44, 0x91, 0x92, 0x32, 0x17, 0xb7, 0x6b, 0x72, 0x46, 0xbe,
	0x6f, 0x6a, 0x71, 0x71, 0x62, 0x7a, 0xaa, 0x90, 0x08, 0x17, 0x6b, 0x59, 0x62, 0x4e, 0x69, 0xaa,
	0x04, 0xa3, 0x02, 0xa3, 0x06, 0x67, 0x10, 0x84, 0x63, 0x14, 0xc2, 0xc5, 0xea, 0x0a, 0x32, 0x4b,
	0xc8, 0x9b, 0x8b, 0x05, 0xa4, 0x5a, 0x48, 0x48, 0x0f, 0x62, 0x11, 0x92, 0x56, 0x29, 0x2c, 0x62,
	0x4a, 0xd2, 0x4d, 0x97, 0x9f, 0x4c, 0x66, 0x12, 0x55, 0x12, 0xd0, 0x2f, 0x33, 0xd4, 0x4f, 0xad,
	0x48, 0xcc, 0x2d, 0xc8, 0x49, 0xd5, 0x4f, 0x4d, 0xce, 0xc8, 0xb7, 0x62, 0xd4, 0x4a, 0x62, 0x03,
	0xbb, 0xc0, 0x18, 0x10, 0x00, 0x00, 0xff, 0xff, 0x22, 0x15, 0x49, 0xe0, 0xbb, 0x00, 0x00, 0x00,
}
