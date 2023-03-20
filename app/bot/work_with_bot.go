package bot

import (
	"context"
	"fmt"
	tgbotapi "github.com/Syfaro/telegram-bot-api"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
	"strconv"
	"strings"
	z "tg_bot/bd"
)

func Work() {

	bot, err := tgbotapi.NewBotAPI(os.Getenv("BOT_TOKEN"))
	if err != nil {
		log.Panic(err, " - trouble with connect bot_token")
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Panic(err, " - trouble with connect BD")
	}

	bot.Debug = false
	log.Printf("Authorized on account %s", bot.Self.UserName)

	updateConfig := tgbotapi.NewUpdate(0)

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		log.Panic(err, " - trouble with update chat")
	}
	for update := range updates {

		ChatID := strconv.Itoa(int(update.Message.Chat.ID))
		Text := update.Message.Text
		Nick := update.Message.From.UserName
		Command := update.Message.Command()
		message := strings.Split(Text, " ")
		reply := ""

		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, Text)

		switch Command {

		case "start":
			reply = "Привет!\nВведи команду вида <Число [Категория]>, чтобы добавить запись в журнал расходов\n\n" +
				"Число не может быть меньше 1, так же бот учитывает расходы без категории\n\n" +
				"Так же у меня есть команды\n" +
				"/week - выводит сгруппированные по категориям траты за неделю отсортированные по убыванию цены\n" +
				"/month - выводит сгруппированные по категориям траты за месяц отсортированные по убыванию цены\n" +
				"/last N  - список из последних N операций в формате <дата: сумма категория>"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
			bot.Send(msg)

		case "week":
			body, _ := z.Week(conn, ChatID)
			var text []string
			for i, v := range body {
				text = append(text, fmt.Sprintf("Sum: %d\t, Category: %s", v.Price, v.Category))
				if len(body)-1 == i {
					continue
				}
				text = append(text, "\t")
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(text, "\n"))
			bot.Send(msg)
		case "month":
			body, _ := z.Month(conn, ChatID)
			var text []string
			for i, v := range body {
				text = append(text, fmt.Sprintf("Sum: %d\t, Category: %s", v.Price, v.Category))
				if len(body)-1 == i {
					continue
				}
				text = append(text, "\t")
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(text, "\n"))
			bot.Send(msg)
		case "last":
			getN := message[1]
			N, _ := strconv.Atoi(getN)
			body, _ := z.Last(conn, ChatID, N)
			var text []string
			for i, v := range body {
				text = append(text, fmt.Sprintf("Time: %s\t Sum: %d\t, Category: %s", v.Time_add.Time.Format("02-Jan-2006 15:04:54"), v.Price, v.Category))
				if len(body)-1 == i {
					continue
				}
				text = append(text, "\t")
			}
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, strings.Join(text, "\n"))
			bot.Send(msg)
		}

		if price, err := strconv.Atoi(message[0]); err == nil && price > 0 {
			if err != nil {
				reply = "Не верный ввод, введите 1)сума 2)категория"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
			}
			if len(message) == 2 {
				z.WriteBD(conn, ChatID, Nick, message[1], price)
				reply = "успешная заПись"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
			} else if len(message) == 1 {
				z.WriteBD(conn, ChatID, Nick, "NotCategory", price)
				reply = "успешная заПись"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
			} else if len(message) > 2 {
				reply = "Не верный ввод, введите 1)сума 2)категория"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
			}
		} else {
			if message[0] != "/"+Command {
				reply = "Не верный ввод, введите 1)сума 2)категория"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, reply)
				bot.Send(msg)
			}
		}
	}
}
