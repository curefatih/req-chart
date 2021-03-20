# ReqChart

## 1. İçerik

### 1.1 Kafka

`docker-compose: zookeeper` & `docker-compose: kafka`

Kafka için docker compose üzerinde bitnami tarafından hazırlanmış ve güncel tutulan [docker imajı](https://hub.docker.com/r/bitnami/kafka) kullanıldı.

### 1.2 Kafka Consumer

`docker-compose: consumer`

Golang ile consumer implemente edildi. 
Ana paket olarak `confluentinc/confluent-kafka-go` kullanıldı  

### 1.3 Arayüz

`docker-compose: ui`

[Socket.io](https://socket.io) ve [Chart.js](https://www.chartjs.org/) ile arayüz hazırlandı. İşlevsellik ön planda olması için görüntü kaygısı güdülmedi.

🌟 `localhost:80` üzerinde çalışmakta.


### 1.4 REST API
`docker-compose: rest_api`

Golang `mux` ve `confluentinc/confluent-kafka-go` ile birlikte kullanıldı. Rest API, kafka producer ile birlikte kullanıldı.
Log kayıtları `logfile` adlı dosyası üzerinde bulunmakta.

### 1.5 WebSocket Server
`docker-compose: ws_server`

Node.js kullanıldı. UI beslemek için socket.io server için hazırlandı. Veri tabanı üzerindeki değişiklikler veritabanı üzerindeki trigger ile (PostgreSQL'de bulunan `notification` ile) dinlenerek değişiklikler socket io ile arayüze gönderildi.

### 1.6 Database
`docker-compose: db`

PostgreSQL imajı üzerine veri akışını yakalamak için ön tanımlı olacak şekilde trigger eklendi. (initdb.sql dosyasından incelenebilir.) 


## 2. Çalıştırma

🌟 Çalıştırmak için terminal (linux ya da macos) üzerinde proje dosyasının içerisinde bulunan docker-compose.yml dosyasıyla aynı dizindeyken `docker-compose up` komutu çalıştırılması yeterli olacaktır.


🌟 Örnek istekler:

- POST isteği `curl --location --request POST 'localhost:8080/'`
- GET isteği `curl --location --request GET 'localhost:8080/'`
- PUT isteği `curl --location --request PUT 'localhost:8080/'` 
- DELETE isteği `curl --location --request DELETE 'localhost:8080/'`

## Anahtar Kodu

```
gAAAAABgUNh0e2AcWHxi8aa5h6nf7Fg1QtSxScOciKo9xXl5M1iu21r4oMMwXnYm70it_Pm6-Cce0VrQSwNTjRpz_u0qNHL00hfZ8Ujk4vlZgZiCzvbJTrosNUJ3s3ftmm0mGh_Z97leoQt-RCokRal9vJgYWYpzHwn2EcouR5DpeHDWBRs0tOC7noNz2frdH7gd5Lom9ipH
```
