package bot

import "context"

type Message struct {
	Context context.Context
	Msg     string
	Matches []string
}
