1. Go to https://jwt.io/ to generate a jwt token with "username" in payload, "secret" as the secret string
2. 

```
go run main.go
```

3. curl -v -H "Authorization: Bearer {your jwt token}" http://localhost:3001/ping
