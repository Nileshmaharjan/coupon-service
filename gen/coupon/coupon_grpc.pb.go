// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.20.3
// source: coupon.proto

package coupon

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// CouponServiceClient is the client API for CouponService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CouponServiceClient interface {
	CreateCampaign(ctx context.Context, in *CreateCampaignRequest, opts ...grpc.CallOption) (*CreateCampaignResponse, error)
	GetCampaign(ctx context.Context, in *GetCampaignRequest, opts ...grpc.CallOption) (*GetCampaignResponse, error)
	IssueCoupon(ctx context.Context, in *IssueCouponRequest, opts ...grpc.CallOption) (*IssueCouponResponse, error)
}

type couponServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCouponServiceClient(cc grpc.ClientConnInterface) CouponServiceClient {
	return &couponServiceClient{cc}
}

func (c *couponServiceClient) CreateCampaign(ctx context.Context, in *CreateCampaignRequest, opts ...grpc.CallOption) (*CreateCampaignResponse, error) {
	out := new(CreateCampaignResponse)
	err := c.cc.Invoke(ctx, "/coupon.CouponService/CreateCampaign", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponServiceClient) GetCampaign(ctx context.Context, in *GetCampaignRequest, opts ...grpc.CallOption) (*GetCampaignResponse, error) {
	out := new(GetCampaignResponse)
	err := c.cc.Invoke(ctx, "/coupon.CouponService/GetCampaign", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *couponServiceClient) IssueCoupon(ctx context.Context, in *IssueCouponRequest, opts ...grpc.CallOption) (*IssueCouponResponse, error) {
	out := new(IssueCouponResponse)
	err := c.cc.Invoke(ctx, "/coupon.CouponService/IssueCoupon", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CouponServiceServer is the server API for CouponService service.
// All implementations must embed UnimplementedCouponServiceServer
// for forward compatibility
type CouponServiceServer interface {
	CreateCampaign(context.Context, *CreateCampaignRequest) (*CreateCampaignResponse, error)
	GetCampaign(context.Context, *GetCampaignRequest) (*GetCampaignResponse, error)
	IssueCoupon(context.Context, *IssueCouponRequest) (*IssueCouponResponse, error)
	mustEmbedUnimplementedCouponServiceServer()
}

// UnimplementedCouponServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCouponServiceServer struct {
}

func (UnimplementedCouponServiceServer) CreateCampaign(context.Context, *CreateCampaignRequest) (*CreateCampaignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCampaign not implemented")
}
func (UnimplementedCouponServiceServer) GetCampaign(context.Context, *GetCampaignRequest) (*GetCampaignResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCampaign not implemented")
}
func (UnimplementedCouponServiceServer) IssueCoupon(context.Context, *IssueCouponRequest) (*IssueCouponResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method IssueCoupon not implemented")
}
func (UnimplementedCouponServiceServer) mustEmbedUnimplementedCouponServiceServer() {}

// UnsafeCouponServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CouponServiceServer will
// result in compilation errors.
type UnsafeCouponServiceServer interface {
	mustEmbedUnimplementedCouponServiceServer()
}

func RegisterCouponServiceServer(s grpc.ServiceRegistrar, srv CouponServiceServer) {
	s.RegisterService(&CouponService_ServiceDesc, srv)
}

func _CouponService_CreateCampaign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCampaignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponServiceServer).CreateCampaign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/coupon.CouponService/CreateCampaign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponServiceServer).CreateCampaign(ctx, req.(*CreateCampaignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CouponService_GetCampaign_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCampaignRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponServiceServer).GetCampaign(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/coupon.CouponService/GetCampaign",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponServiceServer).GetCampaign(ctx, req.(*GetCampaignRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CouponService_IssueCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IssueCouponRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CouponServiceServer).IssueCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/coupon.CouponService/IssueCoupon",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CouponServiceServer).IssueCoupon(ctx, req.(*IssueCouponRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CouponService_ServiceDesc is the grpc.ServiceDesc for CouponService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CouponService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "coupon.CouponService",
	HandlerType: (*CouponServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCampaign",
			Handler:    _CouponService_CreateCampaign_Handler,
		},
		{
			MethodName: "GetCampaign",
			Handler:    _CouponService_GetCampaign_Handler,
		},
		{
			MethodName: "IssueCoupon",
			Handler:    _CouponService_IssueCoupon_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "coupon.proto",
}
