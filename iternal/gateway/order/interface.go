package order

import (
	"context"

	proto "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/proto/order"
	"google.golang.org/grpc"
)

type ordersServiceClient interface {
	GetOrders(ctx context.Context, in *proto.GetOrdersRequest, opts ...grpc.CallOption) (*proto.GetOrdersResponse, error)
	GetOrderById(ctx context.Context, in *proto.GetOrderByIdRequest, opts ...grpc.CallOption) (*proto.GetOrderByIdResponse, error)
}
