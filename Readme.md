# one-shot-url

# これは何？

短縮URL を生成する Webアプリケーションです．

# 仕組み

API に送られた url と紐づく code アプリケーション内で生成し，その情報をデータベースに保存しています．

code の生成には [ksuids](https://github.com/segmentio/ksuid)を使用しています.

# 使い方

1. git clone します．

```
$ git clone https://github.com/Issei0804-ie/who-is-in-lab.git
```

2. env ファイルを作成します．

```
$ cp .env-sample .env
```

3. docker compose でコンテナを立ち上げます．

```
$ docker compose up
(もしくは)
$ docker-compose up
```

4. API サーバーがデプロイされたので，curl等で code を作成できます.

```
$ curl -X POST -i -d '{"url":"https://example.com"}'  localhost:8080/short
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 09 May 2022 03:00:03 GMT
Content-Length: 24

{"short_url":"28uSuqZd"}%
```

5. 勿論, code からコード化された URL を復元することもできます．

```
$ curl -i localhost:8080/28uSuqZd
HTTP/1.1 200 OK
Content-Type: application/json; charset=utf-8
Date: Mon, 09 May 2022 03:08:16 GMT
Content-Length: 38

{"message":"https://example.com"}%
```