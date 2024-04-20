Привет!
=======
Если будут проблемы или вопросы по проекту, пиши на https://t.me/Mur466


Инсталляция
===========

0. Клонируем репозиторий git к себе
-----------------------------------
Раз вы читаете readme.md, вероятно вы уже это сделали
Если нет - 
```
git clone https://github.com/Mur466/distribcalc.git
```

Вариант через Docker
--------------------
```
docker compose up -d
```
Создаются 4 контейнера:
storage - база данных postgres
server - орекстратор и веб-сервер
agent1 - агент вычислений 1
agent2 - агент вычислений 2
Веб-морда сервера доступна по адресу http://localhost:8080/

Можно переходить  к пункту "Отправка заданий и получение результатов"


Вариант без докера, вручную
---------------------------


1. Установка Postgresql
-----------------------
Скачиваем дистрибутив Postgres 16.2
https://www.enterprisedb.com/downloads/postgres-postgresql-downloads

Запускаем инсталлятор
В процесссе установки указываем пароль для системного пользователя __postgres__
Для учебного инстанса можно указать такой __пароль__ же как имя - __postgres__ (с одной S на конце).

Стандартный порт - 5432

Проверка, что всё установлено командой: 
```
pg_config --version
```
Если команда не запускается, проверьте что в вашем PATH есть папка "C:\Program Files\PostgreSQL\16\bin" 

Если postgres установили не локально, а на другой хост/порт (или хотите иметь базу на другом хосте) то нужно в дальнейшем в командах подменять 
localhost на имя вашего хоста, 5432 на номер нужного порта



2. Создание базы данных и объектов в ней
----------------------------------------
Созаем базу distribcalc комадой 
```
createdb -h localhost -p 5432 -U postgres distribcalc
```

После запуска в ответ на запрос вводим вышеуказанный пароль 

Создаем таблицы запуском sql-скрипта
Запускайте команду ниже из папки \distribcalc\database или укажите полный путь к файлу \distribcalc\database\dbmigrate.sql
```
psql -f dbmigrate.sql -U postgres -h localhost -p 5432 -d distribcalc
```

3. Получение необходимых пакетов go

Выполняем в папке distribcalc
```
go mod download
```
После этого все необходимые пакеты должны быть установлены

На всякий случай список команд для получения пакетов вручную
```
go get github.com/gin-gonic/gin
go get -u go.uber.org/zap
go get github.com/jackc/pgx/v5/pgxpool
go get github.com/johncgriffin/overflow
```


Запуск
======
Запуск сервера
```
cd cd cmd\server
go run main.go 
```

С параметрами по-умолчанию веб-морда сервера доступна по адресу http://localhost:8080/

Запуск агентов:
```
cd cd cmd\agent
go run main.go 
```

Отправка заданий и получение результатов
========================================

Поддержимаются простые арифметические операции над целыми числами и скобки
Округление при некратном делении целых чисел происходит по правилам go

Примеры выражений:

- 4*3/(2+1) результат 4
- -1-1 результат -2
- -(3+4) результат -7
- (1+1)+(2+2)+(3+3) результат 12, параллельный расчет до 3 выражений
- 9223372036854775802+9223372036854775802 результат overflow
- 4/(2-2) результат  division by zero


Подавать задания и получать можно как через web-морду, так и через API endpoint /calculate-expression


Логин
```
curl -v --location --request POST http://localhost:8080/signin --header "Content-Type: application/json" --data-raw "{     \"username\": \"user2\",     \"password\": \"2222\" }"
```

Отправка задания от клиента серверу через Curl (формат для Windows, нужно все вложенные кавычки экранировать обратным слешем)
Нужно подменить значение Authorization: Bearer на JWT-токен, полученный при выполнении логина
```
curl http://localhost:8080/calculate-expression  --request "POST" --data "{\"Expr\": \"(1+2)/3*4\", \"ext_id\": \"SomeUniqId6543\"}" --include --header "Content-Type: application/json" --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1MTg3MDMsImlhdCI6MTcxMzUxNTEwMywibmJmIjoxNzEzNTE1MTAzLCJ1c2VybmFtZSI6InVzZXIyIn0.dYEdYm4rKQ5j4SDcfvybOWHOExz_em1iHJvjmQVxP9A"

curl http://localhost:8080/calculate-expression  --request "POST" --data "{\"Expr\": \"9223372036854775802+9223372036854775802 \", \"ext_id\": \"SomeUniqId1214618\"}"  --include --header "Content-Type: application/json" --header "Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3MTM1MTg3MDMsImlhdCI6MTcxMzUxNTEwMywibmJmIjoxNzEzNTE1MTAzLCJ1c2VybmFtZSI6InVzZXIyIn0.dYEdYm4rKQ5j4SDcfvybOWHOExz_em1iHJvjmQVxP9A"
```
ext_id - уникальный идентификатор. При повторной отправке выражения с тем же идентификатором возвращается статус ранее переданного выражения
Если его не передать, то уникальность не контролируется. При этом видеть результат можно через web-морду.
Но для получения результата через API ext_id - обязателен.

 

Более подробное описание
========================
Смотри файл description.md

