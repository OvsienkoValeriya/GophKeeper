// Package main — точка входа клиента GophKeeper.
package main

import (
	"flag"
	"fmt"
	"os"
)

var (
	// Version — версия приложения (устанавливается при сборке).
	Version = "dev"
	// BuildTime — время сборки (устанавливается при сборке).
	BuildTime = "unknown"
	// Commit — коммит (устанавливается при сборке).
	Commit = "unknown"
)

func main() {
	// Основные флаги
	showVersion := flag.Bool("version", false, "показать версию")
	serverAddr := flag.String("server", "localhost:3200", "адрес сервера")
	flag.Parse()

	if *showVersion {
		fmt.Printf("GophKeeper Client\nVersion: %s\nBuild Time: %s\nCommit: %s\n",
			Version, BuildTime, Commit)
		return
	}

	args := flag.Args()
	if len(args) == 0 {
		printUsage()
		os.Exit(1)
	}

	// TODO: Подключение к серверу
	_ = serverAddr

	command := args[0]
	switch command {
	case "register":
		fmt.Println("TODO: Implement register command")
	case "login":
		fmt.Println("TODO: Implement login command")
	case "add":
		fmt.Println("TODO: Implement add command")
	case "list":
		fmt.Println("TODO: Implement list command")
	case "get":
		fmt.Println("TODO: Implement get command")
	case "delete":
		fmt.Println("TODO: Implement delete command")
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`GophKeeper — безопасный менеджер паролей

Использование:
  gophkeeper [флаги] <команда> [аргументы]

Команды:
  register    Регистрация нового пользователя
  login       Вход в систему
  add         Добавить новую запись
  list        Список всех записей
  get         Получить запись по ID
  delete      Удалить запись

Флаги:
  -server     Адрес сервера (по умолчанию: localhost:3200)
  -version    Показать версию

Примеры:
  gophkeeper register --email user@example.com --password secret
  gophkeeper login --email user@example.com --password secret
  gophkeeper add password --name "GitHub" --login "user" --password "pass"
  gophkeeper list
  gophkeeper get --id <uuid>`)
}

