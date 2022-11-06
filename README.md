# Распределенный конфиг

Предоставляет возможность динамически управлять конфигурацией приложений.
Доступ к сервису должен осуществляется с помощью API.
Общение с сервисом  происходит посредством gRPC.
![](/images/diagram.png)

В config_service.proto описан протокол взаимодействия.
Сервис поддерживает следующий набор запросов:

- ### Создание конфига
  
  ```bash
  curl -XPOST -d '{
      "service_name": "managed-k8s",
      "data": {
          "k1": "v1",
          "k2": "v2"
       } 
  }' 'http://localhost:8085/v1/config'
  ```
  
  Вернёт
  
  ```json
  {
  "config": {
      "serviceName": "managed-k8s",
      "data": {
          "k1": "v1",
          "k2": "v2"
      }
  },
  "version": "1",
  "createdAt": "2022-11-06T10:07:30.013235Z"
  }
  ```
  
  Если конфиг с таким именем уже существовал - вернется ошибка  

```json
{
"code": 409,
"message": "Unable to create managed-k8s config: config already exists",
"details": []
}
```

- ### Получение конфига
  
  В таблице конфигов есть поле `relevant`. Если оно установлено, значит эта версия конфига используется приложением. Только одна версия конфига обозначена так. При запросе конфига по имени, без указания версии, возвращается релевантный конфиг. При этом запросе обновляется поле `last_used` конфига, необходимое для определения того, когда он последний раз использовался. Это поле просматривается при удалении.
  
  ```bash
  curl -XGET 'http://localhost:8085/v1/config/managed-k8s'
  ```
  
  Вернёт
  
  ```json
  {
      "config": {
          "serviceName": "managed-k8s",
          "data": {
              "k1": "v1",
              "k2": "v2"
          }
      },
      "version": "1",
      "createdAt": "2022-11-06T10:07:30.013235Z"
  }
  ```

        Или сообщение о том, что конфига с таким именем не существует.

        Также можно получить конкретную версию конфига:

        

* ```bash
  curl -XGET 'http://localhost:8085/v1/config/managed-k8s/1'
  ```
  
  Вернёт
  
  ```json
  {
      "config": {
          "serviceName": "managed-k8s",
          "data": {
              "k1": "v1",
              "k2": "v2"
          }
      },
      "version": "1",
      "createdAt": "2022-11-06T10:07:30.013235Z"
  }
  ```

        Или сообщение о том, что конфига с таким именем и версией не существует.

- ### Обновление конфига
  
  При обновлении конфига у новой версии `relevant` становится равным `true`.  Если конфига с таким именем нет, то вернется сообщение с этой информацией.
  
  ```bash
  curl -XPUT -d '{
      "service_name": "kuber",
      "data": {
          "k3": "v3",
          "k4": "v4"
       } 
  }' 'http://localhost:8085/v1/config/managed-k8s'
  ```

        Вернет        

```json
{
    "config": {
        "serviceName": "managed-k8s",
        "data": {
            "k3": "v3",
            "k4": "v4"
        }
    },
    "version": "2",
    "createdAt": "2022-11-06T10:44:13.104308Z"
}
```

- ### Удаление конфига
  
  Если не указать версию, то удалятся все версии данного конфига. В конфигурации программы можно указать параметр `DELETE_CONFIG_IF_RECENTLY_USED = FALSE` и тогда конфиг не будет удалён, если использовался меньше чем `RECENT_USE_DURATION_DAYS`.
  
  ```bash
  curl -XDELETE 'http://localhost:8085/v1/config/managed-k8s'
  ```
  
  В случае успеха вернёт:
  
  ```json
  {
      "message": "Config was deleted"
  }
  ```
  
  Что будет, если удалить `relevant` версию конфига? Тогда сервис, если ещё есть версии, сделает конфиг с наибольшей версией актуальным.
  
  ```bash
  curl -XDELETE 'http://localhost:8085/v1/config/managed-k8s/2'
  ```
  
  В случае успеха вернёт то же самое.

- ### Получение всех версий конфига
  
  Пока что эта функция доступна только  для вызовов gRPC, из-за необходимости отдельно реализовывать это для REST запросов. Так происходит потому что: `streaming calls are not yet supported in the in-process transport`

- ### Установка актуальной версии конфига
  
  Допустим у нас возникла потребность откатиться к предыдущей версии. Для этого нам будет достаточно сделать так:
  
  ```bash
  curl -XPUT 'http://localhost:8085/v1/config/managed-k8s/2/set_relevant'
  ```
  
  Вернет:
  
  ```json
  {
      "config": {
          "serviceName": "managed-k8s",
          "data": {
              "k3": "v3",
              "k4": "v4"
          }
      },
      "version": "2",
      "createdAt": "2022-11-06T11:16:47.828395Z"
  }
  ```
  
  Теперь при запросе GetConfig мы не указывая версию будет получать нужную нам.

## Запуск

1) Необходимо заполнить файл конфигурации (файл app.env). Парсер сначала посмотрит в .env файлах, затем, если эти значения указаны в переменных среды, он отдаст приоритет им.

   2.1. Для неконтейнеризированного запуска необходимо:
    a) Установить postgresql,
    b) Создать базу данных: 

```bash
createdb your_db_name
```

    и внести данные о базе в makefile (host, user, password, database).

2.2   Чтобы развернуть сервис в докере нужно проделать следующие шаги:   

```bash
docker-compose --env-file app.env up
```

- Войти в контейнер с приложением: 
  
  ```bash
  docker exec -it <ContainerId> bash
  ```

- Установить утилиту migrate: 
  
  ```bash
  curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
  mv migrate.linux-amd64 $GOPATH/bin/migrate
  ```

\-    Выполнить миграцию:

```bash
make migrate_up
```
