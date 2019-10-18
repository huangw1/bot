package bot

import (
	"context"
	"fmt"
	"github.com/fraugster/cli"
)

type CLIAdapter struct {
	Prefix   string
	Messages <-chan string
}

func NewCLIAdapter(ctx context.Context, name string) *CLIAdapter {
	return &CLIAdapter{
		Prefix:   fmt.Sprintf("%s >", name),
		Messages: cli.ReadLines(ctx),
	}
}

func (cli *CLIAdapter) NextMessage() <-chan string {
	fmt.Print(cli.Prefix)
	return cli.Messages
}

func (cli *CLIAdapter) Send(msg string) error {
	fmt.Println(msg)
	return nil
}

func (*CLIAdapter) Close() error {
	fmt.Println()
	return nil
}
