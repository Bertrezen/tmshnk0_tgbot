package cmd

import (
	"github.com/spf13/cobra"
	"tg_bot/rest_api"
)

func start(f *rest_api.Server) *cobra.Command {
	return &cobra.Command{
		Use:   "start",
		Short: "Start tg web-server",
		Run: func(cmd *cobra.Command, args []string) {
			f.Init()        // Инициируем приложение и устанавливаем соединение с БД
			defer f.Close() // Откладываем закрытие соединения

			f.Start() // Запускаем сервер
		},
	}
}
