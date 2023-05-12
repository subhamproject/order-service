package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"context"

	"github.com/gin-gonic/gin"

	"os/signal"
	"syscall"

	"github.com/subhamproject/order-service/ordermgr"
	"github.com/subhamproject/order-service/otelsvc"
	"go.mongodb.org/mongo-driver/mongo"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel/trace"
)

func main() {
	var client *mongo.Client
	var ctx context.Context
	var cFund context.CancelFunc

	wg := sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		fmt.Println("initializing connection mongo...")
		//init mogno db
		client, ctx, cFund, _ = ordermgr.InitMongoDB()
	}()

	wg.Wait()

	log.Printf("initializing otel connection...")
	otelUrl := ordermgr.GetEnvParam("OTEL_COLLECTOR_URL", "localhost:4317")
	otelShutdown := otelsvc.InitTracerProvider(otelUrl)

	r := gin.Default()
	f := func(req *http.Request) bool { return req.URL.Path != "/health" }
	r.Use(otelgin.Middleware("order-service", otelgin.WithFilter(f)))
	r.GET("/health", GetServiceHealthHandler)
	r.GET("/order", GetUserOrderHandler)
	r.POST("/order", CreateUserOrderHandler)

	serverPort := ordermgr.GetEnvParam("SERVICE_PORT", "8081")
	srv := &http.Server{
		Addr:    ":" + serverPort,
		Handler: r,
	}

	// Initializing the server in a goroutine so that
	// it won't block the graceful shutdown handling below
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	//clode mongo driver
	ordermgr.CloseMongoDB(client, ctx, cFund)

	//close otel conn
	otelShutdown()

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	log.Println("Server exiting")
}

func GetUserOrderHandler(c *gin.Context) {
	span := trace.SpanFromContext(c)
	fmt.Println("received user order request headers", c.Request.Header)
	if span.SpanContext().IsValid() {
		fmt.Println("SpanId", span.SpanContext().SpanID())

		if span.SpanContext().HasTraceID() {
			fmt.Println("--TraceId", span.SpanContext().TraceID())
		}
	}

	defer span.End()

	id := c.Query("userId")
	teacher, err := ordermgr.GetUserOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable retrieve order"})
		return
	}
	c.JSON(http.StatusOK, teacher)
}

func CreateUserOrderHandler(c *gin.Context) {
	span := trace.SpanFromContext(c)
	fmt.Println("received user order request headers", c.Request.Header)
	if span.SpanContext().IsValid() {
		fmt.Println("SpanId", span.SpanContext().SpanID())

		if span.SpanContext().HasTraceID() {
			fmt.Println("--TraceId", span.SpanContext().TraceID())
		}
	}
	defer span.End()

	id := c.Query("userId")
	err := ordermgr.CreateUserOrder(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "unable create new user order, error: " + err.Error()})
		return
	}
	c.JSON(http.StatusOK, "success!")
}

func GetServiceHealthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "I'm Healthly")
}
