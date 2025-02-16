package main

import (
	"log/slog"
	"time"

	"gopkg.in/resty.v1"
)

type Pushover struct {
	ApiToken string
	UserKey  string

	client *resty.Client
	logger *slog.Logger
}

type PushoverMessage struct {
	Token   string `json:"token"`
	User    string `json:"user"`
	Title   string `json:"title"`
	Message string `json:"message"`
}

func NewPushover(apiToken, userKey string, logger *slog.Logger) *Pushover {
	client := resty.New()
	client.SetHostURL("https://api.pushover.net/1/messages.json")
	client.SetTimeout(10 * time.Second)
	client.SetHeader("Content-Type", "application/json")

	return &Pushover{
		ApiToken: apiToken,
		UserKey:  userKey,
		client:   client,
		logger:   logger,
	}
}

func (p *Pushover) Send(title, message string) error {
	// Send notification
	msg := &PushoverMessage{
		Token:   p.ApiToken,
		User:    p.UserKey,
		Title:   title,
		Message: message,
	}
	p.logger.Info("Sending pushover message", slog.String("title", title))
	_, err := p.client.R().SetBody(msg).Post("")
	if err != nil {
		p.logger.Error("Error sending pushover message",
			slog.Any("error", err),
			slog.Any("message", msg))
	}
	return err
}
