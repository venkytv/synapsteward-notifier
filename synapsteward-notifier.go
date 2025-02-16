package main

import (
	"context"
	"encoding/json"
	"log"
	"log/slog"
	"os"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/urfave/cli/v3"
)

const (
	NATS_DEFAULT_URL      = "nats://localhost:4222"
	NATS_DEFAULT_STREAM   = "notifications"
	NATS_DEFAULT_CONSUMER = "synapsteward-notifier"
)

type NATSAlert struct {
	Title   string `json:"title"`
	Message string `json:"message"`
}

func main() {

	cmd := &cli.Command{
		Name:  "synapsteward-notifier",
		Usage: "Listen to a NATS stream and send notifications to a Pushover user",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "nats-url",
				Value:   NATS_DEFAULT_URL,
				Usage:   "NATS server URL",
				Sources: cli.EnvVars("NATS_URL"),
			},
			&cli.StringFlag{
				Name:    "nats-stream",
				Value:   NATS_DEFAULT_STREAM,
				Usage:   "NATS stream name",
				Sources: cli.EnvVars("NATS_STREAM"),
			},
			&cli.StringFlag{
				Name:    "nats-consumer",
				Value:   NATS_DEFAULT_CONSUMER,
				Usage:   "NATS consumer name",
				Sources: cli.EnvVars("NATS_CONSUMER"),
			},
			&cli.StringFlag{
				Name:     "pushover-api-token",
				Usage:    "Pushover API token",
				Required: true,
				Sources:  cli.EnvVars("PUSHOVER_API_TOKEN"),
			},
			&cli.StringFlag{
				Name:     "pushover-user-key",
				Usage:    "Pushover user key",
				Required: true,
				Sources:  cli.EnvVars("PUSHOVER_USER_KEY"),
			},
			&cli.BoolFlag{
				Name:    "debug",
				Usage:   "Enable debug logging",
				Sources: cli.EnvVars("DEBUG"),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			natsUrl := cmd.String("nats-url")
			natsStream := cmd.String("nats-stream")
			natsConsumer := cmd.String("nats-consumer")
			pushoverApiToken := cmd.String("pushover-api-token")
			pushoverUserKey := cmd.String("pushover-user-key")
			debug := cmd.Bool("debug")

			logLevel := slog.LevelInfo
			if debug {
				logLevel = slog.LevelDebug
			}
			logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
				Level: logLevel,
			}))

			// Connect to NATS
			logger.Debug("Connecting to NATS", slog.String("url", natsUrl))
			nc, err := nats.Connect(natsUrl)
			if err != nil {
				return err
			}
			defer nc.Close()

			// Connect to JetStream
			logger.Debug("Connecting to JetStream")
			js, err := jetstream.New(nc)
			if err != nil {
				return err
			}

			// Get stream handle
			logger.Debug("Getting stream handle", slog.String("stream", natsStream))
			stream, err := js.Stream(ctx, natsStream)
			if err != nil {
				return err
			}

			// Create a durable consumer
			logger.Debug("Creating durable consumer", slog.String("consumer", natsConsumer))
			consumer, err := stream.CreateOrUpdateConsumer(ctx, jetstream.ConsumerConfig{
				Durable:   natsConsumer,
				AckPolicy: jetstream.AckExplicitPolicy,
			})

			// Create a pushover instance
			logger.Debug("Creating Pushover instance")
			pushover := NewPushover(pushoverApiToken, pushoverUserKey,
				logger.With(slog.String("component", "pushover")))

			// Start listening
			iter, err := consumer.Messages()
			for {
				msg, err := iter.Next()
				logger.Debug("Received message", slog.Any("message", msg))
				if err != nil {
					logger.Error("Error reading message", slog.Any("error", err))
					break
				}

				// Parse JSON message
				var alert NATSAlert
				if err := json.Unmarshal(msg.Data(), &alert); err != nil {
					logger.Error("Error parsing message", slog.Any("error", err))
				} else {
					pushover.Send(alert.Title, alert.Message)
				}
				msg.Ack()
			}

			return nil
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}

}
