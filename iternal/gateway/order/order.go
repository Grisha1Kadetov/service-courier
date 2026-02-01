package order

import (
	"context"
	"fmt"
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/model/order"
	pb "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/proto/order"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type Gateway struct {
	client ordersServiceClient
}

func NewGateway(client ordersServiceClient) *Gateway {
	return &Gateway{
		client: client,
	}
}

func (g *Gateway) GetOrders(ctx context.Context, sicne time.Time) ([]order.Order, error) {
	val, err := g.client.GetOrders(ctx, &pb.GetOrdersRequest{From: timestamppb.New(sicne)})
	if err != nil {
		return nil, err
	}
	orders := val.Orders
	res := make([]order.Order, len(orders))
	for i, o := range orders {
		status, ok := order.ValidStatuses[o.Status]
		if !ok {
			return nil, fmt.Errorf("invalid order status: %s", o.Status)
		}

		res[i] = order.Order{
			ID:        o.Id,
			CreatedAt: o.CreatedAt.AsTime(),
			Status:    status,
		}
	}
	return res, nil
}

func (g *Gateway) GetOrderById(ctx context.Context, id string) (order.Order, error) {
	val, err := g.client.GetOrderById(ctx, &pb.GetOrderByIdRequest{Id: id})
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.NotFound {
			return order.Order{}, order.ErrNotFound
		}
		return order.Order{}, err
	}

	status, ok := order.ValidStatuses[val.Order.Status]
	if !ok {
		return order.Order{}, fmt.Errorf("invalid order status: %s", val.Order.Status)
	}

	res := order.Order{
		ID:        val.Order.Id,
		Status:    status,
		CreatedAt: val.Order.CreatedAt.AsTime(),
	}
	return res, nil
}
