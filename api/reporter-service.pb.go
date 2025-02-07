// Code generated by protoc-gen-go.
// source: github.com/fuserobotics/reporter/api/reporter-service.proto
// DO NOT EDIT!

/*
Package api is a generated protocol buffer package.

It is generated from these files:
	github.com/fuserobotics/reporter/api/reporter-service.proto
	github.com/fuserobotics/reporter/api/reporter.proto

It has these top-level messages:
	RegisterStateRequest
	RegisterStateResponse
	RecordStateRequest
	RecordStateResponse
	ListRemotesRequest
	ListRemotesResponse
	CreateRemoteRequest
	CreateRemoteResponse
	RemoteContext
	StateContext
	RemoteList
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"
import stream "github.com/fuserobotics/statestream"
import view "github.com/fuserobotics/reporter/view"

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

// Register or update a state config
type RegisterStateRequest struct {
	Context      *StateContext  `protobuf:"bytes,1,opt,name=context" json:"context,omitempty"`
	StreamConfig *stream.Config `protobuf:"bytes,2,opt,name=stream_config,json=streamConfig" json:"stream_config,omitempty"`
}

func (m *RegisterStateRequest) Reset()                    { *m = RegisterStateRequest{} }
func (m *RegisterStateRequest) String() string            { return proto.CompactTextString(m) }
func (*RegisterStateRequest) ProtoMessage()               {}
func (*RegisterStateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{0} }

func (m *RegisterStateRequest) GetContext() *StateContext {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *RegisterStateRequest) GetStreamConfig() *stream.Config {
	if m != nil {
		return m.StreamConfig
	}
	return nil
}

type RegisterStateResponse struct {
}

func (m *RegisterStateResponse) Reset()                    { *m = RegisterStateResponse{} }
func (m *RegisterStateResponse) String() string            { return proto.CompactTextString(m) }
func (*RegisterStateResponse) ProtoMessage()               {}
func (*RegisterStateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{1} }

type RecordStateRequest struct {
	Context *StateContext     `protobuf:"bytes,1,opt,name=context" json:"context,omitempty"`
	Report  *view.StateReport `protobuf:"bytes,2,opt,name=report" json:"report,omitempty"`
}

func (m *RecordStateRequest) Reset()                    { *m = RecordStateRequest{} }
func (m *RecordStateRequest) String() string            { return proto.CompactTextString(m) }
func (*RecordStateRequest) ProtoMessage()               {}
func (*RecordStateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{2} }

func (m *RecordStateRequest) GetContext() *StateContext {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *RecordStateRequest) GetReport() *view.StateReport {
	if m != nil {
		return m.Report
	}
	return nil
}

type RecordStateResponse struct {
}

func (m *RecordStateResponse) Reset()                    { *m = RecordStateResponse{} }
func (m *RecordStateResponse) String() string            { return proto.CompactTextString(m) }
func (*RecordStateResponse) ProtoMessage()               {}
func (*RecordStateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{3} }

type ListRemotesRequest struct {
}

func (m *ListRemotesRequest) Reset()                    { *m = ListRemotesRequest{} }
func (m *ListRemotesRequest) String() string            { return proto.CompactTextString(m) }
func (*ListRemotesRequest) ProtoMessage()               {}
func (*ListRemotesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type ListRemotesResponse struct {
	List *RemoteList `protobuf:"bytes,1,opt,name=list" json:"list,omitempty"`
}

func (m *ListRemotesResponse) Reset()                    { *m = ListRemotesResponse{} }
func (m *ListRemotesResponse) String() string            { return proto.CompactTextString(m) }
func (*ListRemotesResponse) ProtoMessage()               {}
func (*ListRemotesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ListRemotesResponse) GetList() *RemoteList {
	if m != nil {
		return m.List
	}
	return nil
}

type CreateRemoteRequest struct {
	Context  *RemoteContext `protobuf:"bytes,1,opt,name=context" json:"context,omitempty"`
	Endpoint string         `protobuf:"bytes,2,opt,name=endpoint" json:"endpoint,omitempty"`
}

func (m *CreateRemoteRequest) Reset()                    { *m = CreateRemoteRequest{} }
func (m *CreateRemoteRequest) String() string            { return proto.CompactTextString(m) }
func (*CreateRemoteRequest) ProtoMessage()               {}
func (*CreateRemoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *CreateRemoteRequest) GetContext() *RemoteContext {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *CreateRemoteRequest) GetEndpoint() string {
	if m != nil {
		return m.Endpoint
	}
	return ""
}

type CreateRemoteResponse struct {
}

func (m *CreateRemoteResponse) Reset()                    { *m = CreateRemoteResponse{} }
func (m *CreateRemoteResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateRemoteResponse) ProtoMessage()               {}
func (*CreateRemoteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func init() {
	proto.RegisterType((*RegisterStateRequest)(nil), "api.RegisterStateRequest")
	proto.RegisterType((*RegisterStateResponse)(nil), "api.RegisterStateResponse")
	proto.RegisterType((*RecordStateRequest)(nil), "api.RecordStateRequest")
	proto.RegisterType((*RecordStateResponse)(nil), "api.RecordStateResponse")
	proto.RegisterType((*ListRemotesRequest)(nil), "api.ListRemotesRequest")
	proto.RegisterType((*ListRemotesResponse)(nil), "api.ListRemotesResponse")
	proto.RegisterType((*CreateRemoteRequest)(nil), "api.CreateRemoteRequest")
	proto.RegisterType((*CreateRemoteResponse)(nil), "api.CreateRemoteResponse")
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// Client API for ReporterService service

type ReporterServiceClient interface {
	RegisterState(ctx context.Context, in *RegisterStateRequest, opts ...grpc.CallOption) (*RegisterStateResponse, error)
	RecordState(ctx context.Context, in *RecordStateRequest, opts ...grpc.CallOption) (*RecordStateResponse, error)
	ListRemotes(ctx context.Context, in *ListRemotesRequest, opts ...grpc.CallOption) (*ListRemotesResponse, error)
	CreateRemote(ctx context.Context, in *CreateRemoteRequest, opts ...grpc.CallOption) (*CreateRemoteResponse, error)
}

type reporterServiceClient struct {
	cc *grpc.ClientConn
}

func NewReporterServiceClient(cc *grpc.ClientConn) ReporterServiceClient {
	return &reporterServiceClient{cc}
}

func (c *reporterServiceClient) RegisterState(ctx context.Context, in *RegisterStateRequest, opts ...grpc.CallOption) (*RegisterStateResponse, error) {
	out := new(RegisterStateResponse)
	err := grpc.Invoke(ctx, "/api.ReporterService/RegisterState", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reporterServiceClient) RecordState(ctx context.Context, in *RecordStateRequest, opts ...grpc.CallOption) (*RecordStateResponse, error) {
	out := new(RecordStateResponse)
	err := grpc.Invoke(ctx, "/api.ReporterService/RecordState", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reporterServiceClient) ListRemotes(ctx context.Context, in *ListRemotesRequest, opts ...grpc.CallOption) (*ListRemotesResponse, error) {
	out := new(ListRemotesResponse)
	err := grpc.Invoke(ctx, "/api.ReporterService/ListRemotes", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reporterServiceClient) CreateRemote(ctx context.Context, in *CreateRemoteRequest, opts ...grpc.CallOption) (*CreateRemoteResponse, error) {
	out := new(CreateRemoteResponse)
	err := grpc.Invoke(ctx, "/api.ReporterService/CreateRemote", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// Server API for ReporterService service

type ReporterServiceServer interface {
	RegisterState(context.Context, *RegisterStateRequest) (*RegisterStateResponse, error)
	RecordState(context.Context, *RecordStateRequest) (*RecordStateResponse, error)
	ListRemotes(context.Context, *ListRemotesRequest) (*ListRemotesResponse, error)
	CreateRemote(context.Context, *CreateRemoteRequest) (*CreateRemoteResponse, error)
}

func RegisterReporterServiceServer(s *grpc.Server, srv ReporterServiceServer) {
	s.RegisterService(&_ReporterService_serviceDesc, srv)
}

func _ReporterService_RegisterState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReporterServiceServer).RegisterState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ReporterService/RegisterState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReporterServiceServer).RegisterState(ctx, req.(*RegisterStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReporterService_RecordState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RecordStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReporterServiceServer).RecordState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ReporterService/RecordState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReporterServiceServer).RecordState(ctx, req.(*RecordStateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReporterService_ListRemotes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRemotesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReporterServiceServer).ListRemotes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ReporterService/ListRemotes",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReporterServiceServer).ListRemotes(ctx, req.(*ListRemotesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReporterService_CreateRemote_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRemoteRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReporterServiceServer).CreateRemote(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ReporterService/CreateRemote",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReporterServiceServer).CreateRemote(ctx, req.(*CreateRemoteRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _ReporterService_serviceDesc = grpc.ServiceDesc{
	ServiceName: "api.ReporterService",
	HandlerType: (*ReporterServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "RegisterState",
			Handler:    _ReporterService_RegisterState_Handler,
		},
		{
			MethodName: "RecordState",
			Handler:    _ReporterService_RecordState_Handler,
		},
		{
			MethodName: "ListRemotes",
			Handler:    _ReporterService_ListRemotes_Handler,
		},
		{
			MethodName: "CreateRemote",
			Handler:    _ReporterService_CreateRemote_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "github.com/fuserobotics/reporter/api/reporter-service.proto",
}

func init() {
	proto.RegisterFile("github.com/fuserobotics/reporter/api/reporter-service.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 513 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x93, 0xcb, 0x6e, 0x13, 0x3d,
	0x14, 0xc7, 0x95, 0xf6, 0x53, 0x3f, 0x38, 0x49, 0xa9, 0xea, 0xa4, 0x34, 0x1d, 0x21, 0x84, 0xcc,
	0x86, 0xeb, 0x4c, 0x68, 0x16, 0x48, 0x61, 0x83, 0x94, 0x0d, 0x0b, 0xc4, 0xc2, 0x65, 0x1f, 0x4d,
	0x26, 0xa7, 0x83, 0xa5, 0x89, 0x3d, 0xd8, 0x4e, 0xa8, 0x84, 0xba, 0xe9, 0x2b, 0x20, 0xf1, 0x62,
	0xbc, 0x02, 0x8f, 0xc1, 0x02, 0xcd, 0xb1, 0x73, 0x99, 0x24, 0x08, 0x95, 0x4d, 0xa4, 0xe3, 0x73,
	0xfb, 0xfd, 0xe7, 0xfc, 0x03, 0x6f, 0x72, 0xe9, 0x3e, 0xcd, 0xc6, 0x71, 0xa6, 0xa7, 0xc9, 0xe5,
	0xcc, 0xa2, 0xd1, 0x63, 0xed, 0x64, 0x66, 0x13, 0x83, 0xa5, 0x36, 0x0e, 0x4d, 0x92, 0x96, 0x72,
	0x19, 0xbc, 0xb4, 0x68, 0xe6, 0x32, 0xc3, 0xb8, 0x34, 0xda, 0x69, 0xb6, 0x9f, 0x96, 0x32, 0x7a,
	0x90, 0x6b, 0x9d, 0x17, 0x48, 0x75, 0xa9, 0x52, 0xda, 0xa5, 0x4e, 0x6a, 0x65, 0x7d, 0x49, 0xd4,
	0xfb, 0xd3, 0x7c, 0xeb, 0x52, 0x87, 0xd6, 0x19, 0x4c, 0xa7, 0x49, 0xa6, 0xd5, 0xa5, 0xcc, 0x43,
	0x47, 0xff, 0x56, 0x44, 0x7f, 0x5b, 0xb3, 0x6c, 0x9a, 0x4b, 0xfc, 0x42, 0x3f, 0xbe, 0x83, 0x5f,
	0x41, 0x47, 0x60, 0x2e, 0xad, 0x43, 0x73, 0x51, 0xa1, 0x08, 0xfc, 0x3c, 0x43, 0xeb, 0xd8, 0x73,
	0xf8, 0x3f, 0xd3, 0xca, 0xe1, 0x95, 0xeb, 0x36, 0x1e, 0x35, 0x9e, 0x34, 0xcf, 0x8f, 0xe3, 0xb4,
	0x94, 0x31, 0xd5, 0x0c, 0x7d, 0x42, 0x2c, 0x2a, 0x58, 0x1f, 0x0e, 0xbd, 0x84, 0x91, 0x97, 0xd0,
	0xdd, 0xa3, 0x96, 0x7b, 0xb1, 0x7f, 0x8d, 0x87, 0xf4, 0x2a, 0x5a, 0x3e, 0xf4, 0x11, 0x3f, 0x85,
	0x93, 0x8d, 0xcd, 0xb6, 0xd4, 0xca, 0x22, 0x2f, 0x80, 0x09, 0xcc, 0xb4, 0x99, 0xfc, 0x3b, 0xd0,
	0x53, 0x38, 0xf0, 0x8a, 0x03, 0xc9, 0x71, 0x4c, 0x92, 0xc3, 0xc0, 0x2a, 0x21, 0x42, 0x01, 0x3f,
	0x81, 0x76, 0x6d, 0x5b, 0x80, 0xe8, 0x00, 0x7b, 0x2f, 0xad, 0x13, 0x38, 0xd5, 0x0e, 0x6d, 0x80,
	0xe0, 0x03, 0x68, 0xd7, 0x5e, 0x7d, 0x31, 0x7b, 0x0c, 0xff, 0x15, 0xd2, 0x2e, 0xc0, 0x8e, 0x08,
	0xcc, 0xd7, 0x50, 0x35, 0x25, 0xf9, 0x08, 0xda, 0x43, 0x83, 0xb4, 0xa3, 0xca, 0x2c, 0x74, 0xbd,
	0xd8, 0xd4, 0xc5, 0xd6, 0xda, 0xb7, 0x84, 0x45, 0x70, 0x07, 0xd5, 0xa4, 0xd4, 0x52, 0x79, 0x69,
	0x77, 0xc5, 0x32, 0xe6, 0xf7, 0xa1, 0x53, 0x5f, 0xe0, 0xe9, 0xce, 0x7f, 0xed, 0xc3, 0x91, 0x08,
	0xf7, 0xbf, 0xf0, 0xc6, 0x65, 0xdf, 0x1b, 0x70, 0x58, 0xfb, 0xfa, 0xec, 0x2c, 0xac, 0xdd, 0xf6,
	0x42, 0x14, 0xed, 0x4a, 0x85, 0xef, 0xf4, 0xe1, 0xe6, 0xc7, 0xcf, 0x6f, 0x7b, 0xef, 0xf8, 0xeb,
	0x64, 0xfe, 0x2a, 0x29, 0x74, 0x96, 0x16, 0xc9, 0xd7, 0x00, 0x5b, 0xb9, 0xb0, 0xd4, 0x0a, 0x95,
	0xbb, 0x5e, 0xbd, 0x91, 0xdf, 0x47, 0x72, 0x72, 0x9d, 0x64, 0xc4, 0x3a, 0xa8, 0x3b, 0x87, 0xdd,
	0x34, 0xa0, 0xb9, 0x76, 0x0f, 0x76, 0x1a, 0x76, 0x6f, 0xfa, 0x21, 0xea, 0x6e, 0x27, 0x02, 0xd2,
	0x5b, 0x42, 0x1a, 0xf0, 0xde, 0x6d, 0x91, 0x06, 0xc1, 0x13, 0xec, 0x23, 0x34, 0xd7, 0xce, 0x1c,
	0x18, 0xb6, 0xed, 0x10, 0x18, 0x76, 0x38, 0x82, 0x33, 0x62, 0x68, 0x31, 0xa8, 0x18, 0x0c, 0x25,
	0xd9, 0x0c, 0x5a, 0xeb, 0xf7, 0x61, 0xbe, 0x7b, 0x87, 0x27, 0xa2, 0xb3, 0x1d, 0x99, 0x30, 0xb8,
	0x47, 0x83, 0x9f, 0xf1, 0x87, 0xab, 0xc1, 0x2b, 0x25, 0x3e, 0x26, 0x29, 0x4b, 0x5b, 0x8c, 0x0f,
	0xe8, 0x8f, 0xde, 0xff, 0x1d, 0x00, 0x00, 0xff, 0xff, 0x01, 0xe2, 0x1e, 0x72, 0xe3, 0x04, 0x00,
	0x00,
}
