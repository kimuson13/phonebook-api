# 電話帳API
## 仕様
- 電話帳のデータの所得
- データのIDによる絞り込み
- データの追加
- データの編集
- データの削除

## demo
### 電話帳データの獲得
#### コマンド
`curl http://localhost:8080/api/phonebooks | jq .`
#### 結果
```
{
    "id": 1,
    "name": "kimu",
    "phone": "09012345678"
}
{
    "id": 2,
    "name": "kimu2",
    "phone": "09098765432"
}
```
### データのIDによる絞り込み
#### コマンド
`curl http://localhost:8080/api/phonebooks/2 | jq .`
#### 結果
```
{
    "id": 2,
    "name": "kimu2",
    "phone": "09098765432"
}
```
### データの追加
#### コマンド
`curl -X POST -d "{\"name\":\"add\",\"phone\":\"09087876565\"}" "http://localhost:8080/api/phonebooks" | jq .`
#### 結果
```
{
    "id": 3,
    "name": "add",
    "phone": "09087876565"
}
```
### データの編集
#### コマンド
`curl -X PUT -d "{\"name\":\"put\",\"phone\":\"08012345678\"}" "http://localhost:8080/api/phonebooks/1" | jq .`
#### 結果
```
{
    "id": 1,
    "name": "put",
    "phone": "09012345678"
}
```
### データの削除
#### コマンド
`curl -X DELETE http://localhost:8080/api/phonebooks/1 | jq .`
#### 結果
```
{
    "id": 0,
    "name": "",
    "phone": ""
}
```

## 今後追加したい機能
- データの検索(名前か電話番号)