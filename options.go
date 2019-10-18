package bot

type Option func(bot *Bot) error

func WithBrain(brain Brain) Option {
	return func(bot *Bot) error {
		bot.Brain = brain
		return nil
	}
}
