package bot

import (
	"context"
	"fmt"
	"github.com/fraugster/cli"
	"go.uber.org/zap"
)

type Bot struct {
	Context  context.Context
	Name     string
	Adapter  Adapter
	Logger   *zap.Logger
	handlers []responseHandler
}

func New(name string) *Bot {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	ctx := cli.Context()
	return &Bot{
		Context: ctx,
		Name:    name,
		Adapter: NewCLIAdapter(ctx, name),
		Logger:  logger,
	}
}

func (b *Bot) Run() {
	b.Logger.Info("start bot", zap.String("name", b.Name))
	for {
		select {
		case <-b.Context.Done():
			break
		case msg := <-b.Adapter.NextMessage():
			b.handleMessage(msg)
		}
	}
}

func (b *Bot) handleMessage(s string) {
	msg := Message{
		Context: b.Context,
		Msg:     s,
	}

	for _, handler := range b.handlers {
		matches := handler.regexp.FindStringSubmatch(s)
		if len(matches) == 0 {
			continue
		}
		msg.Matches = matches[1:]
		err := handler.run(msg)
		if err != nil {
			b.Logger.Error("failed to handle message", zap.Error(err))
		}
	}
}

func (b *Bot) Respond(msg string, fun RespondFunc) {
	expr := "^(?i)" + msg + "$"
	h, err := newHandler(expr, fun)
	if err != nil {
		b.Logger.Fatal("failed to add Response handler", zap.Error(err))
	}

	b.handlers = append(b.handlers, h)
}

func (b *Bot) Say(msg string, args ...interface{}) {
	if len(args) > 0 {
		msg = fmt.Sprintf(msg, args...)
	}

	err := b.Adapter.Send(msg)
	if err != nil {
		b.Logger.Error("failed to send message", zap.Error(err))
	}
}
