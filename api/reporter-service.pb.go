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
	ListStatesRequest
	ListStatesResponse
	GetStateRequest
	GetStateResponse
	ListRemotesRequest
	ListRemotesResponse
	CreateRemoteRequest
	CreateRemoteResponse
	RemoteContext
	StateContext
	StateReport
	RemoteList
	StateList
	StateListComponent
	StateListState
	StateQuery
*/
package api

import proto "github.com/golang/protobuf/proto"
import fmt "fmt"
import math "math"
import _ "github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis/google/api"
import stream "github.com/fuserobotics/statestream"

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
	Context *StateContext `protobuf:"bytes,1,opt,name=context" json:"context,omitempty"`
	Report  *StateReport  `protobuf:"bytes,2,opt,name=report" json:"report,omitempty"`
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

func (m *RecordStateRequest) GetReport() *StateReport {
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

type ListStatesRequest struct {
}

func (m *ListStatesRequest) Reset()                    { *m = ListStatesRequest{} }
func (m *ListStatesRequest) String() string            { return proto.CompactTextString(m) }
func (*ListStatesRequest) ProtoMessage()               {}
func (*ListStatesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{4} }

type ListStatesResponse struct {
	List *StateList `protobuf:"bytes,1,opt,name=list" json:"list,omitempty"`
}

func (m *ListStatesResponse) Reset()                    { *m = ListStatesResponse{} }
func (m *ListStatesResponse) String() string            { return proto.CompactTextString(m) }
func (*ListStatesResponse) ProtoMessage()               {}
func (*ListStatesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{5} }

func (m *ListStatesResponse) GetList() *StateList {
	if m != nil {
		return m.List
	}
	return nil
}

type GetStateRequest struct {
	Context *StateContext `protobuf:"bytes,1,opt,name=context" json:"context,omitempty"`
	Query   *StateQuery   `protobuf:"bytes,2,opt,name=query" json:"query,omitempty"`
}

func (m *GetStateRequest) Reset()                    { *m = GetStateRequest{} }
func (m *GetStateRequest) String() string            { return proto.CompactTextString(m) }
func (*GetStateRequest) ProtoMessage()               {}
func (*GetStateRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{6} }

func (m *GetStateRequest) GetContext() *StateContext {
	if m != nil {
		return m.Context
	}
	return nil
}

func (m *GetStateRequest) GetQuery() *StateQuery {
	if m != nil {
		return m.Query
	}
	return nil
}

type GetStateResponse struct {
	State *StateReport `protobuf:"bytes,1,opt,name=state" json:"state,omitempty"`
}

func (m *GetStateResponse) Reset()                    { *m = GetStateResponse{} }
func (m *GetStateResponse) String() string            { return proto.CompactTextString(m) }
func (*GetStateResponse) ProtoMessage()               {}
func (*GetStateResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{7} }

func (m *GetStateResponse) GetState() *StateReport {
	if m != nil {
		return m.State
	}
	return nil
}

type ListRemotesRequest struct {
}

func (m *ListRemotesRequest) Reset()                    { *m = ListRemotesRequest{} }
func (m *ListRemotesRequest) String() string            { return proto.CompactTextString(m) }
func (*ListRemotesRequest) ProtoMessage()               {}
func (*ListRemotesRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{8} }

type ListRemotesResponse struct {
	List *RemoteList `protobuf:"bytes,1,opt,name=list" json:"list,omitempty"`
}

func (m *ListRemotesResponse) Reset()                    { *m = ListRemotesResponse{} }
func (m *ListRemotesResponse) String() string            { return proto.CompactTextString(m) }
func (*ListRemotesResponse) ProtoMessage()               {}
func (*ListRemotesResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{9} }

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
func (*CreateRemoteRequest) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{10} }

func (m *CreateRemoteRequest) GetContext() *RemoteContext {
	if m != nil {
		return m.Context
	}
	return nil
}

type CreateRemoteResponse struct {
}

func (m *CreateRemoteResponse) Reset()                    { *m = CreateRemoteResponse{} }
func (m *CreateRemoteResponse) String() string            { return proto.CompactTextString(m) }
func (*CreateRemoteResponse) ProtoMessage()               {}
func (*CreateRemoteResponse) Descriptor() ([]byte, []int) { return fileDescriptor0, []int{11} }

func init() {
	proto.RegisterType((*RegisterStateRequest)(nil), "api.RegisterStateRequest")
	proto.RegisterType((*RegisterStateResponse)(nil), "api.RegisterStateResponse")
	proto.RegisterType((*RecordStateRequest)(nil), "api.RecordStateRequest")
	proto.RegisterType((*RecordStateResponse)(nil), "api.RecordStateResponse")
	proto.RegisterType((*ListStatesRequest)(nil), "api.ListStatesRequest")
	proto.RegisterType((*ListStatesResponse)(nil), "api.ListStatesResponse")
	proto.RegisterType((*GetStateRequest)(nil), "api.GetStateRequest")
	proto.RegisterType((*GetStateResponse)(nil), "api.GetStateResponse")
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
const _ = grpc.SupportPackageIsVersion3

// Client API for ReporterService service

type ReporterServiceClient interface {
	RegisterState(ctx context.Context, in *RegisterStateRequest, opts ...grpc.CallOption) (*RegisterStateResponse, error)
	RecordState(ctx context.Context, in *RecordStateRequest, opts ...grpc.CallOption) (*RecordStateResponse, error)
	ListStates(ctx context.Context, in *ListStatesRequest, opts ...grpc.CallOption) (*ListStatesResponse, error)
	GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*GetStateResponse, error)
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

func (c *reporterServiceClient) ListStates(ctx context.Context, in *ListStatesRequest, opts ...grpc.CallOption) (*ListStatesResponse, error) {
	out := new(ListStatesResponse)
	err := grpc.Invoke(ctx, "/api.ReporterService/ListStates", in, out, c.cc, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *reporterServiceClient) GetState(ctx context.Context, in *GetStateRequest, opts ...grpc.CallOption) (*GetStateResponse, error) {
	out := new(GetStateResponse)
	err := grpc.Invoke(ctx, "/api.ReporterService/GetState", in, out, c.cc, opts...)
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
	ListStates(context.Context, *ListStatesRequest) (*ListStatesResponse, error)
	GetState(context.Context, *GetStateRequest) (*GetStateResponse, error)
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

func _ReporterService_ListStates_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListStatesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReporterServiceServer).ListStates(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ReporterService/ListStates",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReporterServiceServer).ListStates(ctx, req.(*ListStatesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _ReporterService_GetState_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetStateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ReporterServiceServer).GetState(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/api.ReporterService/GetState",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ReporterServiceServer).GetState(ctx, req.(*GetStateRequest))
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
			MethodName: "ListStates",
			Handler:    _ReporterService_ListStates_Handler,
		},
		{
			MethodName: "GetState",
			Handler:    _ReporterService_GetState_Handler,
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
	Metadata: fileDescriptor0,
}

func init() {
	proto.RegisterFile("github.com/fuserobotics/reporter/api/reporter-service.proto", fileDescriptor0)
}

var fileDescriptor0 = []byte{
	// 612 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x09, 0x6e, 0x88, 0x02, 0xff, 0x9c, 0x54, 0xcb, 0x6e, 0xd3, 0x40,
	0x14, 0x55, 0x0a, 0x2d, 0xe5, 0xf6, 0x91, 0xf6, 0x26, 0xe9, 0xc3, 0x42, 0x08, 0x0d, 0x02, 0x55,
	0x3c, 0xec, 0x92, 0x2c, 0xa8, 0xc2, 0x06, 0x29, 0x0b, 0x58, 0x20, 0x24, 0xa6, 0xdd, 0x47, 0x8e,
	0x33, 0x0d, 0x23, 0x1c, 0x8f, 0x33, 0x33, 0xa9, 0x8a, 0x50, 0x37, 0xfd, 0x05, 0x24, 0xfe, 0x88,
	0x2f, 0xe0, 0x17, 0xf8, 0x10, 0xe4, 0x99, 0x49, 0x6c, 0xc7, 0x61, 0x91, 0x2e, 0xef, 0xf3, 0x9c,
	0x33, 0x3e, 0xd7, 0xf0, 0x6e, 0xc4, 0xf5, 0xd7, 0xe9, 0xc0, 0x8f, 0xc4, 0x38, 0xb8, 0x9c, 0x2a,
	0x26, 0xc5, 0x40, 0x68, 0x1e, 0xa9, 0x40, 0xb2, 0x54, 0x48, 0xcd, 0x64, 0x10, 0xa6, 0x7c, 0x1e,
	0xbc, 0x56, 0x4c, 0x5e, 0xf1, 0x88, 0xf9, 0xa9, 0x14, 0x5a, 0xe0, 0xbd, 0x30, 0xe5, 0xde, 0xa3,
	0x91, 0x10, 0xa3, 0x98, 0x99, 0xbe, 0x30, 0x49, 0x84, 0x0e, 0x35, 0x17, 0x89, 0xb2, 0x2d, 0x5e,
	0x67, 0xa5, 0xfd, 0x6e, 0xe8, 0xf4, 0x7f, 0x43, 0x4a, 0x87, 0x9a, 0x29, 0x2d, 0x59, 0x38, 0x0e,
	0x22, 0x91, 0x5c, 0xf2, 0x91, 0x9d, 0x20, 0xd7, 0xd0, 0xa4, 0x6c, 0xc4, 0x95, 0x66, 0xf2, 0x3c,
	0xeb, 0xa1, 0x6c, 0x32, 0x65, 0x4a, 0xe3, 0x4b, 0x78, 0x10, 0x89, 0x44, 0xb3, 0x6b, 0x7d, 0x54,
	0x7b, 0x52, 0x3b, 0xd9, 0x6a, 0xef, 0xfb, 0x61, 0xca, 0x7d, 0xd3, 0xd3, 0xb3, 0x05, 0x3a, 0xeb,
	0xc0, 0x0e, 0xec, 0xd8, 0xdd, 0x7d, 0xbb, 0xfb, 0x68, 0xcd, 0x8c, 0xec, 0xfa, 0x36, 0xeb, 0xf7,
	0x4c, 0x96, 0x6e, 0xdb, 0xd0, 0x46, 0xe4, 0x10, 0x5a, 0x0b, 0xc8, 0x2a, 0x15, 0x89, 0x62, 0xe4,
	0x1b, 0x20, 0x65, 0x91, 0x90, 0xc3, 0xbb, 0x13, 0x3a, 0x81, 0x0d, 0xfb, 0x32, 0x8e, 0xc9, 0x5e,
	0xde, 0x4b, 0x4d, 0x9e, 0xba, 0x3a, 0x69, 0x41, 0xa3, 0x04, 0xe6, 0x38, 0x34, 0x60, 0xff, 0x13,
	0x57, 0xda, 0x24, 0x95, 0xa3, 0x40, 0xce, 0x00, 0x8b, 0x49, 0xdb, 0x8a, 0x04, 0xee, 0xc7, 0x5c,
	0xcd, 0x58, 0xed, 0xe6, 0x48, 0x59, 0x2f, 0x35, 0x35, 0xc2, 0xa0, 0xfe, 0x81, 0xe9, 0xbb, 0xeb,
	0x79, 0x06, 0xeb, 0x93, 0x29, 0x93, 0xdf, 0x9d, 0x9c, 0x7a, 0xde, 0xfa, 0x25, 0x4b, 0x53, 0x5b,
	0x25, 0x5d, 0xd8, 0xcb, 0x61, 0x1c, 0xbd, 0xe7, 0xb0, 0x6e, 0x3e, 0xbe, 0x43, 0xa9, 0xbe, 0x84,
	0x2d, 0x93, 0xa6, 0x15, 0x47, 0xd9, 0x58, 0x14, 0x24, 0x77, 0xa1, 0x51, 0xca, 0xba, 0xa5, 0x4f,
	0x4b, 0x9a, 0x2d, 0x1d, 0xdb, 0x53, 0x10, 0xdd, 0x87, 0x46, 0x4f, 0x32, 0x03, 0x94, 0x55, 0x66,
	0xc2, 0x5f, 0x2d, 0x0a, 0xc7, 0xc2, 0x78, 0x45, 0xb9, 0x07, 0x9b, 0x2c, 0x19, 0xa6, 0x82, 0x27,
	0xf6, 0x5b, 0x3e, 0xa4, 0xf3, 0x98, 0x1c, 0x40, 0xb3, 0x0c, 0x60, 0xd9, 0xb5, 0x7f, 0xaf, 0x43,
	0x9d, 0xba, 0xc3, 0x38, 0xb7, 0x77, 0x87, 0xbf, 0x6a, 0xb0, 0x53, 0xb2, 0x1b, 0x1e, 0x3b, 0xd8,
	0xaa, 0xf9, 0x3d, 0x6f, 0x59, 0xc9, 0x39, 0xe3, 0xf3, 0xed, 0x9f, 0xbf, 0x3f, 0xd7, 0x3e, 0x92,
	0xb7, 0xc1, 0xd5, 0x9b, 0x20, 0x16, 0x51, 0x18, 0x07, 0x3f, 0x1c, 0xd9, 0xec, 0xec, 0x52, 0x91,
	0xb0, 0x44, 0xdf, 0xe4, 0x39, 0xf3, 0xba, 0x7d, 0x3e, 0xbc, 0x09, 0x22, 0xc3, 0xb5, 0x5b, 0x3e,
	0x15, 0xbc, 0xad, 0xc1, 0x56, 0xc1, 0x81, 0x78, 0xe8, 0xb0, 0x17, 0x0f, 0xc0, 0x3b, 0xaa, 0x16,
	0x1c, 0xa5, 0xf7, 0x86, 0x52, 0x97, 0x9c, 0xae, 0x4a, 0xa9, 0xeb, 0xae, 0x00, 0x2f, 0x00, 0x72,
	0x67, 0xe3, 0x81, 0x41, 0xaa, 0xf8, 0xdf, 0x3b, 0xac, 0xe4, 0x1d, 0x81, 0x96, 0x21, 0x50, 0xc7,
	0x9d, 0x9c, 0x40, 0x18, 0xc7, 0x38, 0x81, 0xcd, 0x99, 0x1d, 0xb1, 0x69, 0x66, 0x17, 0x8e, 0xc0,
	0x6b, 0x2d, 0x64, 0xdd, 0xbe, 0x33, 0xb3, 0xaf, 0x8d, 0x2b, 0x0b, 0xc2, 0x0b, 0xd8, 0x2a, 0xf8,
	0x15, 0x73, 0xc6, 0x65, 0x5f, 0xbb, 0xc7, 0x5c, 0x62, 0x6d, 0x82, 0x06, 0x7b, 0x1b, 0x21, 0xc3,
	0x96, 0xa6, 0x88, 0x53, 0xd8, 0x2e, 0x1a, 0x0d, 0xed, 0xf4, 0x12, 0x73, 0x7b, 0xc7, 0x4b, 0x2a,
	0x6e, 0xf1, 0xa9, 0x59, 0xfc, 0x82, 0x3c, 0xce, 0x17, 0xe7, 0x0a, 0x6c, 0x6c, 0xbe, 0xc9, 0xdc,
	0xdf, 0x83, 0x0d, 0xf3, 0x8b, 0xee, 0xfc, 0x0b, 0x00, 0x00, 0xff, 0xff, 0xa7, 0x99, 0x93, 0x65,
	0x6b, 0x06, 0x00, 0x00,
}
