### Сервер

```bash
    # 1. Запускаем сервер
    go run ./cmd/server/main.go
```

### Клиент

```bash
    # 1. Регистрация пользователя
    go run ./cmd/client/main.go register -u testuser2 -p testuser2
    # Проверяем в базе что создался пользователь и пароль захешировался.

    # 2. Логин пользователя
    go run ./cmd/client/main.go login -u testuser2 -p testuser2
    # Убедились, что мы получили сообщение о том, что мастер-ключ не инициализирован

    # 3. Проверяем токены
    cat ~/.gophkeeper/tokens.json
    # Видим user_id, access_token, has_master_key=false
    # Убедились, что refresh_token`а нет

    # 4. Вызываем refresh
    go run ./cmd/client/main.go refresh
    # Убедились, что refresh_token появился

    # 5. Logout пользователя
    go run ./cmd/client/main.go logout
    # Убедились, что файл с токенами удалился

    # 6. Инициализируем мастер-ключ
    # Login
    go run ./cmd/client/main.go init
    # Проверяем, что в базе появились salt и verifier
    # Проверяем в файле токенов флаг has_master_key

    # 7. Добавляем секреты (< 1 Мб)
    go run ./cmd/client/main.go set -n test@gmail.com -v qwerty  -t credentials
    go run ./cmd/client/main.go set -n '1111 2222 3333 4444' -v '123' -t card
    # Проверяем, что в postgres базе появились данные, data - зашифровано

    # 8. Получаем секреты
    go run ./cmd/client/main.go get test@gmail.com
    go run ./cmd/client/main.go get '1111 2222 3333 4444' 

    # 9. Добавляем секреты (> 1 Мб)
    # для примера создадим большой файл
    dd if=/dev/urandom bs=1024 count=1500 of=/tmp/bigtestbinary.bin
    ls -lh /tmp/bigtestbinary.bin

    go run ./cmd/client/main.go set -n "bigbinaryfile" -f /tmp/bigtestbinary.bin -t binary
    # проверяем в postgres базе object key
    # проверяем в minio 

    # Генерируем и загружаем большой текстовый файл
    yes "This is a test line for GophKeeper storage testing. Lorem ipsum dolor sit amet." | head -n 20000 > /tmp/bigtext.txt
    ls -lh /tmp/bigtext.txt
    go run ./cmd/client/main.go set -n "big-text-note" -f /tmp/bigtext.txt -t text

    # 9. Получаем секреты (> 1 Мб)
    go run ./cmd/client/main.go get bigbinaryfile | head -20
    go run ./cmd/client/main.go get big-text-note | head -20


```