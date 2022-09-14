# stop-checker.com

## Commands

Backend
```
go run ./cmd/server/main.go --config=dev
go test ./...
go test -bench=. -benchtime=20s ./...

go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen generate

docker build -t stop-checker:latest .
docker run -p 3001:3001 stop-checker:latest --config=dev
```

Frontend
```
npm install
npm run start
npm run build
npm run codegen
```