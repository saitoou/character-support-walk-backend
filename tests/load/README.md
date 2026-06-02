# character-support-walk-backend 負荷テスト結果まとめ

## 実施環境

* Backend: Go + Echo
* DB: PostgreSQL (Docker Compose)
* 負荷試験ツール: k6
* 実行環境: ローカル
* 対象API:

  * POST /auth/dev-login
  * GET /home

---

# 1. Smoke Test

## 目的

* API疎通確認
* JWT取得確認
* 認証付きAPIアクセス確認

## シナリオ

```txt
POST /auth/dev-login
↓
GET /home
```

## k6設定

```js
vus: 1
iterations: 1
```

## 結果

```txt
checks_succeeded: 100%
http_req_failed: 0%
```

### レスポンス

```txt
avg: 17.43ms
p95: 27.88ms
```

## 確認できたこと

* /auth/dev-login 正常動作
* access_token取得成功
* JWT認証成功
* /home 正常レスポンス
* 基本疎通問題なし

---

# 2. Load Test (Rate Limit 有効)

## 目的

* GET /home に対する並列アクセス確認
* RateLimiter動作確認

## RateLimiter設定

```go
middleware.NewRateLimiterMemoryStore(20)
```

### 意味

```txt
1IPあたり 20req/sec
```

---

## k6設定

```js
stages: [
    { duration: "30s", target: 10 },
    { duration: "1m", target: 30 },
    { duration: "30s", target: 50 },
    { duration: "30s", target: 0 },
]

sleep(0.2)
```

## 結果

```txt
平均RPS: 約104 req/sec
```

### レスポンス

```txt
p95: 24.42ms
```

### 失敗率

```txt
http_req_failed: 81.65%
```

### ステータス

```txt
429 Too Many Requests
```

## 原因

* k6は localhost 単一IPからアクセス
* 1IPあたり20req/sec制限に到達
* RateLimiter正常発火

## 確認できたこと

* RateLimiterが正常に動作
* DoS/連打制御が有効
* 高負荷時も高速に429返却

---

# 3. Load Test (Rate Limit 緩和後)

## 目的

* API自体の性能確認
* DB/Repository/JWT middleware性能確認

## RateLimiter設定変更

```go
middleware.NewRateLimiterMemoryStore(1000)
```

---

## k6設定

```js
stages: [
    { duration: "30s", target: 10 },
    { duration: "1m", target: 30 },
    { duration: "30s", target: 50 },
    { duration: "30s", target: 0 },
]

sleep(0.2)
```

## 結果

```txt
最大VU: 50
平均RPS: 約102 req/sec
http_req_failed: 0%
checks_succeeded: 100%
```

### レスポンス

```txt
avg: 12.43ms
p90: 21.38ms
p95: 26.66ms
max: 121.9ms
```

## 確認できたこと

* GET /home は100req/sec程度で安定
* JWT middleware問題なし
* DB SELECT系は安定
* Echo middleware問題なし
* Repository層問題なし
* PostgreSQLローカル構成で十分高速

---

# 4. Walk Flow Load Test (10VU)

## 目的

* 実ユーザー行動に近いフロー負荷試験
* SELECT / INSERT / UPDATE 混在時の確認

## シナリオ

```txt
POST /auth/dev-login
↓
GET /home
↓
GET /walk-options
↓
POST /walks
↓
PATCH /walks/:walkID/complete
```

## k6設定

```js
stages: [
    { duration: "30s", target: 5 },
    { duration: "1m", target: 10 },
    { duration: "30s", target: 0 },
]

sleep(1)
```

## 結果

```txt
最大VU: 10
完了フロー数: 647
HTTPリクエスト数: 2589
平均RPS: 約21.5 req/sec
http_req_failed: 0%
checks_succeeded: 100%
```

### レスポンス

```txt
avg: 6.29ms
p90: 17.77ms
p95: 19.51ms
max: 33.7ms
```

## 確認できたこと

* GET /home 正常
* GET /walk-options 正常
* POST /walks 正常
* walk_id返却正常
* PATCH /walks/:walkID/complete 正常
* INSERT/UPDATE含むフローでも安定
* transaction含む書き込み系でも高速

---

# 5. Walk Flow Load Test (100VU)

## 目的

* 高並列環境での実ユーザーフロー負荷試験
* JWT / middleware / DB書き込み / transaction を含む性能確認

## シナリオ

```txt
POST /auth/dev-login
↓
GET /home
↓
GET /walk-options
↓
POST /walks
↓
PATCH /walks/:walkID/complete
```

## k6設定

```js
stages: [
    { duration: "30s", target: 30 },
    { duration: "1m", target: 100 },
    { duration: "30s", target: 0 },
]

sleep(1)
```

## 結果

```txt
最大VU: 100
完了フロー数: 5,596
HTTPリクエスト数: 22,385
平均RPS: 約185.6 req/sec
http_req_failed: 0%
checks_succeeded: 100%
```

### レスポンス

```txt
avg: 12.73ms
p90: 32.03ms
p95: 42.52ms
max: 125.22ms
```

## 確認できたこと

* SELECT / INSERT / UPDATE 混在でも安定
* JWT middleware問題なし
* Echo middleware問題なし
* Repository層問題なし
* PostgreSQL書き込み性能問題なし
* 100VU / 約185RPSでも failure 0%
* 高並列環境でも walk flow 完走可能

---

# 現状の所感

## 良かった点

```txt
- 100VUでも安定
- 約185RPSでもfailure 0%
- p95 約43ms
- JWT認証込みでも高速
- INSERT/UPDATE込みでも安定
- middleware込みでも安定
- RateLimiter正常動作確認済み
```

## まだ未確認な点

```txt
- DB connection pool限界
- Cloud環境での性能
- Redisなど外部依存追加時
- Push通知処理追加後
- 長時間耐久試験
- p99分析
- DBデータ大量増加時
```
