package tgbot

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"time"

	"github.com/mymmrac/telego"
	th "github.com/mymmrac/telego/telegohandler"
)

func RunBot(token string) error {
	// Init signal handling
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	// bot
	bot, err := telego.NewBot(token, telego.WithDefaultDebugLogger()) // Убрать потом этот логгер!
	if err != nil {
		return fmt.Errorf("bot creation error: %v", err)
	}
	if err := SetBotMenu(ctx, bot); err != nil {
		return fmt.Errorf("bot commans setting error: %v", err)
	}

	// signal handling
	updates, _ := bot.UpdatesViaLongPolling(ctx, nil)
	bh, _ := th.NewBotHandler(bot, updates)
	InitRouter(bh, bot)
	// init done chan
	done := make(chan struct{}, 1)

	// handle stop signal
	go func() {
		<-ctx.Done()
		fmt.Println("Stopping...")

		stopCtx, stopCancel := context.WithTimeout(context.Background(), time.Second*20)
		defer stopCancel()

	loop:
		for len(updates) > 0 {
			select {
			case <-stopCtx.Done():
				break loop
			case <-time.After(time.Microsecond * 100):
				// Continue
			}
		}
		fmt.Println("Long polling done")

		_ = bh.StopWithContext(stopCtx)
		fmt.Println("Bot handler done")
		done <- struct{}{}
	}()
	go func() { _ = bh.Start() }()
	fmt.Println("Handling updates...")

	<-done
	fmt.Println("Bot stopped gracefully")
	return nil
}
