# url_shortener

## Конфигурационные файлы

Параметры запуска передаются в JSON файлах `server.json` и `postgres.json`. Примеры файлов находятся в `/configs`.

В файле `server.json` указывается прослушиваемый порт и хранилище данных. Если в параметре `postgres` записано значение `false`, сервер будет использовать in-memory хранилище.

Если в параметре `postgres` записано значение `true`, сервер в качетсве хранилища будет использова базу данных PostgreSQL. В файле `postgres.json` должны быть указаны данные для подключения к базе данных. 

## Хранилище данных

Для работы с базой данных приложению требуется таблица с именем `links`

```postgresql
CREATE TABLE links (
    id SERIAL PRIMARY KEY,
    url TEXT NOT NULL UNIQUE 
)
```

## Образ Docker

В `/build` находится файл `docker-compose.yaml` для запуска сервиса в Docker. При запуске будет использован образ из [репозитория](https://hub.docker.com/r/a0x5f/url_shortener) Docker Hub.

## API

Запрос для получения короткой ссылки

```http request
POST localhost:3001/short
Content-Type: multipart/form-data

{
  "url" : "http://ozon.ru"
}
```
Ответ
```json
{
  "url": "https://lnk.dev/AAAAAAAAAB"
}
```

Запрос для получения полной ссылки

```http request
GET http://localhost:3001
    /full
    ?url=https://lnk.dev/AAAAAAAAAB
```
Ответ
```json
{
  "url": "http://ozon.ru"
}
```