package main

import (
	"context"
	"net/http"
	_ "net/http/pprof"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/db"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/enviroment"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/interceptor/common/retry"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/kafka"
	prometheus "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/metrics"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/middleware/basicauth"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/middleware/metrics"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/middleware/ratelimiter"
	"github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/pkg/log"
	"github.com/IBM/sarama"
	"github.com/go-chi/chi/v5"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	commonServer "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/server/common"
	courierServer "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/server/courier"
	deliveryServer "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/server/delivery"

	courierHandler "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/courier"
	courierRepo "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/repository/courier"
	courierService "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/service/courier"

	deliveryHandler "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/delivery"
	deliveryRepo "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/repository/delivery"
	deliveryService "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/service/delivery"

	orderGateway "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/gateway/order"
	orderHandler "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/handler/order"
	pbOrder "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/proto/order"
	orderService "github.com/Avito-courses/course-go-avito-Grisha1Kadetov/iternal/service/order"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	l := log.NewStdLogger()
	env := enviroment.NewEnvironment(l)
	env.Load()
	pool, err := db.NewPool(ctx, env.GetDSN())
	if err != nil {
		l.Fatal("unable to connect to database", log.NewField("err", err))
	}
	defer pool.Close()

	couriersRepo := courierRepo.NewCourierRepository(pool)
	courierService := courierService.NewCourierService(couriersRepo)
	courierHandler := courierHandler.NewCourierHandler(courierService)

	deliveryRepo := deliveryRepo.NewDeliveryRepository(pool)
	watcher := deliveryService.NewWatcher(deliveryRepo)
	watcher.RunWatcherDelivery(ctx, env.TickDelivery)
	dtf := deliveryService.DefaultDeliveryTimeFactory{}
	deliveryService := deliveryService.NewDeliveryService(deliveryRepo, couriersRepo, &dtf, l)
	deliveryHandler := deliveryHandler.NewDeliveryHandler(deliveryService)

	conn, err := grpc.NewClient(env.OrderServiceHost, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(retry.NewRetryInterceptor(env.RetryAttempts, env.RetrySleep).Interceptor))
	if err != nil {
		l.Error("unable to connect to order-service", log.NewField("err", err))
	}
	defer func() { _ = conn.Close() }()
	pbClient := pbOrder.NewOrdersServiceClient(conn)
	gateway := orderGateway.NewGateway(pbClient)
	orderProcesserFactory := orderService.NewOrderProcessorFactory(deliveryService)
	orderService := orderService.NewOrderService(gateway, orderProcesserFactory)

	config := sarama.NewConfig()
	config.Version = sarama.V3_3_0_0
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetNewest
	kafkaConsumer := kafka.NewGroupConsumer(env.KafkaEnv.GroupID, env.KafkaEnv.Brokers, env.KafkaEnv.Topics, config, l)
	orderHandler := orderHandler.NewOrderHandler(orderService)
	go func() {
		if err := kafkaConsumer.Run(ctx, orderHandler); err != nil {
			l.Error("kafka consumer error", log.NewField("err", err))
		}
	}()

	r := chi.NewRouter()
	prometheus.InitMetrics()
	r.Use(metrics.NewMetricsMiddleware(l), ratelimiter.New(l, env.Rps, env.Burst).Middleware())
	r.Get("/metrics", promhttp.Handler().ServeHTTP)
	r = commonServer.NewRouter(r)
	r = courierServer.NewRouter(courierHandler, r)
	r = deliveryServer.NewRouter(deliveryHandler, r)

	server := &http.Server{
		Addr:    env.Addr + ":" + strconv.Itoa(env.Port),
		Handler: r,
	}

	pprofR := chi.NewRouter()
	pprofR.Use(basicauth.New(env.PprofUsername, env.PprofPassword).Middleware())
	pprofR.Mount("/debug", http.DefaultServeMux)

	pprofServer := &http.Server{
		Addr:    env.PprofAddr,
		Handler: pprofR,
	}

	go func() {
		l.Info("Start up service-courier")
		if err := server.ListenAndServe(); err != nil {
			l.Error("server error", log.NewField("err", err))
			stop()
		}
	}()

	go func() {
		l.Info("Start up pprof")
		if err := pprofServer.ListenAndServe(); err != nil {
			l.Error("server error", log.NewField("err", err))
		}
	}()

	<-ctx.Done()
	l.Info("Shutting down service-courier")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	err = server.Shutdown(shutdownCtx)
	if err != nil {
		l.Error("failed to shutdown server", log.NewField("err", err))
	}
}
