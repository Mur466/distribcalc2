Инструкция по установке без докера
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
psql -f dbmigrate.sql -U postgres -h localhost -p 5432 -d distribcalc2
```

3. Получение необходимых пакетов go

Выполняем в папке distribcalc
```
go mod tidy
```
После этого все необходимые пакеты должны быть установлены



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
Можно запустить одновременно произвольное количество агентов