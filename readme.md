# GroupFi Proxy API

## url 
```
http://localhost:3010
```

## public apis
### POST /msg/receive, receive the latest message 
#### params
```json
{
    "account": "",
    "groupId": "",
    "message": ""
}
```
#### response
```json
{
    "result": true
}
```
or error
```json
{
    "result"  : false,
    "err-code": 1,
    "err-msg" : ""
}
```
