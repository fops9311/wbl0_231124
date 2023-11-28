# wbl0_231124
## HTTP 
* сервер доступен на порту 8090 по умолчанию
* Эндпоинты:

/mem_len - GET JSON MEMCACHE LEN (TOTAL)

/order_pages/last - GET JSON ARRAY OF IDs ON LAST PAGE

/order_pages/:id - GET JSON ARRAY OF IDs ON PAGE {:id}

/order/:id - GET JSON ORDER {:id}


## NATS
* Микросервис слушает топик nats "TESTING"
* В случае получения в по топику данных в формате JSON, они сохраняются в мапе и в базе postgresql.

* Схема конвейера:
```
initmemchan -> |(fin)...

natschan -> msgdatachan -> orderschan   |(fout) -> memchan -> |(fin) -> memwritechan (inmemcacheConsumer)

...|(fout) -> encodechan -> gobchan (databaseConsumer)
```

* natschan - chan *nats.Msg в который транслируются все сообщения из топика
* msgdatachan - chan []byte в который попадают только не пустые payload
* orderschan chan OrderWithKey в который попадают Json.Unmarshal(payload) в тип data.RawOrderData (тип генерируется через go generate инструментом ./scripts/analize_struct). Key для хранения в бд и в мапе генерируется на этом этапе
* fout - fan out разветвление канала orderschan
* gobchan - chan GobWithKey в который попадают gob encoded data.RawOrderData с ключем Key

* initmemchan - chan OrderWithKey в который попадают считанные из postgresql данные при инициализации микросервиса. Данные считываются конкурентно с выполнением микросервиса и он доступен в процессе инициализации данными.
* fin - fan in слияние каналов в memwritechan
* inmemcacheConsumer - узел, сохраняющий данные в кэше из канала memwritechan
* databaseConsumer - узел, сохраняющий данные в postgresql из канала gobchan. В postgresql данные хранятся в таблице со структурой [id text,data bytea], где id - Key, data - gob encoded []byte

## Запуск
go generate

требования: linux,docker,docker-compose

для работы на windows потребуется изменить директивы generate.

## Publisher
go run ./testers/publisher/.  

## STEP 1
[STEP 1](https://github.com/fops9311/wbl0_231124/blob/main/screenshots/1%20order_pages_last.png?raw=true)
## STEP 2
[STEP 2](https://github.com/fops9311/wbl0_231124/blob/main/screenshots/2%20new%20orders%20send.png)
## STEP 3
[STEP 3](https://github.com/fops9311/wbl0_231124/blob/main/screenshots/3%20new%20orders%20on%20server.png)
## STEP 4
[STEP 4](https://github.com/fops9311/wbl0_231124/blob/main/screenshots/4%20new%20orders%20when%20server%20stopped.png)
## STEP 5
[STEP 5](https://github.com/fops9311/wbl0_231124/blob/main/screenshots/5%20new%20orders%20reappear%20when%20started.png)
## STEP 6
[STEP 6](https://github.com/fops9311/wbl0_231124/blob/main/screenshots/6%20individual%20order.png)
## STEP 7
[STEP 7](https://github.com/fops9311/wbl0_231124/blob/main/screenshots/7%20grafana%20orders.png)
