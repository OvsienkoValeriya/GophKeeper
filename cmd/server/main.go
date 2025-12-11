// Package main — точка входа сервера GophKeeper.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
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
	// Флаги командной строки
	addr := flag.String("addr", ":3200", "адрес сервера")
	showVersion := flag.Bool("version", false, "показать версию")
	flag.Parse()

	if *showVersion {
		fmt.Printf("GophKeeper Server\nVersion: %s\nBuild Time: %s\nCommit: %s\n",
			Version, BuildTime, Commit)
		return
	}

	// Настройка логирования
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	log.Printf("Starting GophKeeper server v%s", Version)

	// Создание контекста с отменой
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Обработка сигналов завершения
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	// Создание слушателя
	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	log.Printf("Server listening on %s", *addr)

	// Горутина для graceful shutdown
	go func() {
		sig := <-sigCh
		log.Printf("Received signal: %v", sig)
		cancel()
		lis.Close()
	}()

	// TODO: Инициализация gRPC сервера и хранилища
	// server := grpc.NewServer(...)
	// pb.RegisterGophKeeperServer(server, ...)

	<-ctx.Done()
	log.Println("Server stopped")
}

