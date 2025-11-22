package http

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BMokarzel/OTEL/server-a/pkg/logger"
	"github.com/go-chi/chi/v5"
	_ "github.com/gogo/protobuf/protoc-gen-gogo/grpc"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	_ "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	_ "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type Server struct {
	Server *http.Server
	Router *Router
	Otel   *ServerOtel
}

type ServerOpts struct {
	ServerPort          string
	ServerEnableTracing bool
	ServerEnv           string
	ServerName          string
	ServerVersion       string
}

type ServerOtel struct {
	ServerName   string
	CollectorUrl string
}

func InitTracer(ctx context.Context, connectionUrl, serverName string) (func(ctx context.Context) error, error) {
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serverName),
		),
	)
	if err != nil {
		return nil, err
	}

	conn, err := grpc.NewClient(connectionUrl, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	traceExporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithGRPCConn(conn))
	if err != nil {
		return nil, err
	}

	bsp := sdktrace.NewBatchSpanProcessor(traceExporter)

	traceProvider := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	otel.SetTracerProvider(traceProvider)
	otel.SetTextMapPropagator(propagation.TraceContext{})

	return traceProvider.Shutdown, nil
}

func InitMetrics(ctx context.Context, metricsUrl string) (func(ctx context.Context) error, error) {
	exporter, err := otlpmetricgrpc.New(
		ctx,
		otlpmetricgrpc.WithEndpoint(metricsUrl),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithDialOption(grpc.WithTransportCredentials(insecure.NewCredentials())),
	)
	if err != nil {
		return nil, err
	}

	meterProvider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(
		sdkmetric.NewPeriodicReader(exporter),
	))

	otel.SetMeterProvider(meterProvider)

	runtime.Start(
		runtime.WithMinimumReadMemStatsInterval(10*time.Second),
		runtime.WithMeterProvider(meterProvider),
	)

	return meterProvider.Shutdown, nil
}

func NewServer(log *logger.Logger, opts ServerOpts, otelOpts ServerOtel) *Server {

	routerOptions := routerOpts{
		ServerEnableTracing: opts.ServerEnableTracing,
		ServerEnv:           opts.ServerEnv,
		ServerName:          opts.ServerName,
		ServerVersion:       opts.ServerVersion,
	}

	router := NewRouter(log, routerOptions, opts.ServerName)

	server := &http.Server{
		Addr:              fmt.Sprintf(":%s", opts.ServerPort),
		ReadTimeout:       30 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
		ReadHeaderTimeout: 3 * time.Second,
		Handler:           router.Engine,
	}

	return &Server{
		Server: server,
		Router: router,
		Otel: &ServerOtel{
			ServerName:   opts.ServerName,
			CollectorUrl: otelOpts.CollectorUrl,
		},
	}
}

func (s *Server) Init() error {
	serverCtx, serverStopCtx := context.WithCancel(context.Background())

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	tracerShotdown, err := InitTracer(serverCtx, s.Otel.CollectorUrl, s.Otel.ServerName)
	if err != nil {
		return err
	}

	metricsShotdown, err := InitMetrics(serverCtx, s.Otel.CollectorUrl)
	if err != nil {
		return err
	}

	go func() {
		<-sig
		shutdownCtx, _ := context.WithTimeout(serverCtx, 30*time.Second)

		_ = tracerShotdown
		_ = metricsShotdown

		go func() {
			<-shutdownCtx.Done()
			if shutdownCtx.Err() == context.DeadlineExceeded {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		if err := s.Server.Shutdown(shutdownCtx); err != nil {
			log.Fatal(err)
		}
		serverStopCtx()
	}()

	fmt.Printf("Starting server on port %s\n", s.Server.Addr)

	chi.Walk(s.Router.Engine, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		fmt.Printf("[%s]: '%s' has %d middleware\n", method, route, len(middlewares))
		return nil
	})

	if err := s.Server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}

	<-serverCtx.Done()

	return nil
}
