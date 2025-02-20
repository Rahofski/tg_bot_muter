bot.Handle("/help", func(c telebot.Context) error {
		return c.Send( "Привет! Бот обладает всего одной командной - напиши в чат `/kaban` и узнай, как он работает:)")
	})