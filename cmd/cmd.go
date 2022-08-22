package cmd

import (
	"context"
	"flag"
	"go-http-server-template/pkg/httpserver"
	"io"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

var (
	httpAddr = ":8800" // bind address
)

const (
	httpTimeout = 3 * time.Second // timeouts used to protect the server
)

// run accepts the program arguments and where to send output (default: stdout)
func Run(args []string, _ io.Writer) error {
	var (
		port string
	)

	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	flags.StringVar(&port, "port", "8080", "set port")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	if port != "" {
		httpAddr = ":" + port
	}

	log.Println("Starting Application on port", httpAddr)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	errGrp, egCtx := errgroup.WithContext(ctx)
	s := httpserver.NewServer(httpAddr, httpTimeout, httpTimeout)

	// http server
	errGrp.Go(func() error {
		return s.Start(egCtx)
	})

	// signal handler
	errGrp.Go(func() error {
		return handleSignals(egCtx, cancel)
	})

	return errGrp.Wait()
}

// handleSignals will handle Interrupts or termination signals
func handleSignals(ctx context.Context, cancel context.CancelFunc) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
	select {
	case s := <-sigCh:
		log.Printf("got signal %v, stopping", s)
		cancel()
		return nil
	case <-ctx.Done():
		log.Println("context is done")
		return ctx.Err()
	}
}
