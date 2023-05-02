# A Place For Your Thoughts
Данное веб-приложение позволяет регистрироваться и авторизоваться пользователям, чтобы добавлять новые записи в блоге с поддержкой Markdown. Неавторизированные пользователи имеют доступ только к чтению постов<br>
Реализован RESTful API с использованием PostgreSQL.<br>
Приложение в процессе покрытия тестами (табличные, интеграционные и с использованием моков).<br>

## TLS и переменные окружения
Для корректной работы необходимо сгенерировать TLS-сертификат, а затем создать .env файл, содержащий следующие переменные окружения:
- **POSTGRES_USER** - имя пользователя БД, от которого создаются и наделяются правами пользователи, а так же создаются базы данных (по умолчанию в PostgreSQL POSTGRES_USER=postgres)
- **POSTGRES_PASSWORD** - пароль для POSTGRES_USER (опционально)
- **APP_DB** - название основной базы данных
- **APP_DB_USER** - с этим именем будет создан пользователь с правами SELECT, INSERT, UPDATE, DELETE в основной базе данных
- **APP_DB_PASSWORD** - пароль для APP_DB_USER
- **TEST_DB** - название тестовой базы данных
- **TEST_DB_USER** - с этим именем будет создан пользователь со всеми правами на тестовую базу данных
- **TEST_DB_PASSWORD** - пароль для TEST_DB_USER
- **TLS_CERT** - путь до файла cert.pem
- **TLS_KEY** - путь до файла key.pem<br>
Сгенерировать сертификат можно с помощью следующей команды:<br>
`go run /usr/local/go/src/crypto/tls/generate_cert.go --rsa-bits=2048 --host=localhost`

Пример .env файла
```
POSTGRES_USER=postgres
POSTGRES_PASSWORD=postgres
APP_DB=blog
APP_DB_USER=user
APP_DB_PASSWORD=user
TEST_DB=test
TEST_DB_USER=test
TEST_DB_PASSWORD=test
TLS_CERT=./cert.pem
TLS_KEY=./key.pem
```

## Запуск веб-приложения
### С использованием docker

Команда `make all` выполняет следующее:<br>
- С помощью `docker-compose up --detach` поднимается контейнер с PostgreSQL базами данных и пользователями с необходимыми привилегиями (используются переменные окружения из .env файла). 
- Собирается и запускается сервер

### Флаги приложения
- -addr - адрес веб-приложения в формате host:port (по-умолчанию, на https://localhost:4000). Например, `./server -addr=:4000`
- -dbURL - адрес базы данных в формате `postgres://<user>:<password>@localhost:5432/<database_name>`
- -tls-cert - путь до публичного ключа (например, ./cert.pem)
- -tls-key - путь до приватного ключа (например, ./key.pem) 

## Тесты

Для запуска тестов достаточно команды `make test`
