package bot

import (
	"context"
	"fmt"
	"github.com/fraugster/cli"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strings"
)

type Bot struct {
	Context  context.Context
	Name     string
	Brain    Brain
	Adapter  Adapter
	Logger   *zap.Logger
	initErr  error
	handlers []responseHandler
}

func New(name string, options ...Option) *Bot {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}

	ctx := cli.Context()

	bot := &Bot{
		Context: ctx,
		Name:    name,
		Logger:  logger,
	}

	for _, option := range options {
		err := option(bot)
		if err != nil && bot.initErr == nil {
			bot.initErr = err
		}
	}

	if bot.Adapter == nil {
		bot.Adapter = NewCLIAdapter(bot.Context, name)
	}

	if bot.Brain == nil {
		bot.Brain = NewInMemoryBrain()
	}
	return bot
}

func (b *Bot) Run() error {
	if b.initErr != nil {
		return errors.Wrap(b.initErr, "failed to initialize bot")
	}

	b.Logger.Info("start bot", zap.String("name", b.Name))
	for {
		select {
		case msg := <-b.Adapter.NextMessage():
			b.handleMessage(msg)

		case <-b.Context.Done():
			err := b.Adapter.Close()
			b.Logger.Info("bot is shutting down", zap.String("name", b.Name))
			if err != nil {
				b.Logger.Info("error while closing adapter", zap.Error(err))
			}
			return nil
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
	expr := "^" + msg + "$"
	b.RespondRegex(expr, fun)
}

func (b *Bot) RespondRegex(expr string, fun RespondFunc) {
	if expr == "" {
		return
	}

	if expr[0] == '^' {
		if !strings.HasPrefix(expr, "^(?i)") {
			expr = "^(?i)" + expr[1:]
		}
	} else {
		if !strings.HasPrefix(expr, "(?i)") {
			expr = "(?i)" + expr
		}
	}

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
