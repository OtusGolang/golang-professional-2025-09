package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/redis/go-redis/v9"
)

func main() {
	ctx := context.Background()
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Адрес сервера Redis
		DB:   0,                // Номер базы данных (по умолчанию 0)
	})
	// Закрытие клиента Redis
	defer client.Close()

	const queueName = "myqueue"

	for i := 0; i < 3; i++ {
		message := fmt.Sprintf("Сообщение для обработки %d", i)

		// Отправка сообщения в очередь
		err := client.LPush(ctx, queueName, message).Err()
		if err != nil {
			panic(err)
		}
	}

	for i := 0; i < 3; i++ {
		result, errPop := client.LPop(ctx, queueName).Result()
		if errPop != nil && !errors.Is(errPop, redis.Nil) {
			panic(errPop)
		}

		if result != "" {
			fmt.Println(result)
		}
	}

}
