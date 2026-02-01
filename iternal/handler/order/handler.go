package order

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/IBM/sarama"
)

type OrderHandler struct {
	orderService orderService
}

func NewOrderHandler(orderService orderService) *OrderHandler {
	return &OrderHandler{
		orderService: orderService,
	}
}

func (h *OrderHandler) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *OrderHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (h *OrderHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		if err := h.HandleMessage(session.Context(), message.Value); err != nil {
			return err
		}
		session.MarkMessage(message, "")
	}
	return nil
}

func (h *OrderHandler) HandleMessage(ctx context.Context, message []byte) error {
	var orderDTO Order
	if err := json.Unmarshal(message, &orderDTO); err != nil {
		return fmt.Errorf("failed to unmarshal order message: %w", err)
	}
	orderModel, err := orderDTO.ToModel()
	if err != nil {
		return err
	}
	if err := h.orderService.ProcessOrder(ctx, orderModel); err != nil {
		return err
	}
	return nil
}
