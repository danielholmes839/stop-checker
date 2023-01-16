# stop-checker.com

## Commands

Backend

```
go run ./cmd/server/main.go --config=dev

go get github.com/99designs/gqlgen
go run github.com/99designs/gqlgen generate

docker build -t stop-checker:latest .
docker run -p 3001:3001 stop-checker:latest --config=dev
```

Go test commands

```
go test ./...
go test -bench=. -benchtime=20s ./...
go test ./features/travel/v3/... -run=XXX -bench=. -benchtime=20s -cpuprofile cpu.pprof
go tool pprof -http=:8080 cpux.pprof
```

Frontend

```
npm install
npm run start
npm run build
npm run codegen
```

# Deployment

```
docker build -t stop-checker:latest .

docker network create stop-checker-network

# Backend deployment
docker run --network stop-checker-network --name stop-checker-backend -v "C:\..\backend\data:/app/data" -p 5003:3000 -d -m 1500mb --oom-kill-disable stop-checker:latest --config=prod

# OSRM deployment
docker run --network stop-checker-network --name osrm -d -t -v "C:\..\backend\data:/data" ghcr.io/project-osrm/osrm-backend osrm-routed --algorithm mld /data/osrm/ottawa.osrm
```

docker build -t stop-checker:latest .
docker run --network stop-checker-network --name stop-checker-backend -v "ENTER-YOUR-PATH/backend/data:/app/data" -p 5003:3000 -d stop-checker:latest --config=dev
docker run --network stop-checker-network --name osrm -d -t -v "ENTER-YOUR-PATH/backend/data:/data" ghcr.io/project-osrm/osrm-backend osrm-routed --algorithm mld /data/osrm/ottawa.osrm

/root/stopchecker2-update/data

docker run --network stop-checker-network --name stop-checker-backend -v "C:\Users\danie\Desktop\stopchecker2\backend\data:/app/data" -p 5003:3000 -d -m 1500mb --oom-kill-disable stop-checker:latest --config=prod

# OSRM

https://extract.bbbike.org/

```
// Run to create OSRM indexes
docker run --rm -t -v "C:\..\backend\data:/data" ghcr.io/project-osrm/osrm-backend osrm-extract -p /opt/foot.lua /data/osrm/ottawa.pbf
docker run --rm -t -v "C:\..\backend\data:/data" ghcr.io/project-osrm/osrm-backend osrm-partition /data/osrm/ottawa.osrm
docker run --rm -t -v "C:\..\backend\data:/data" ghcr.io/project-osrm/osrm-backend osrm-customize /data/osrm/ottawa.osrm

// Run OSRM
docker run -d -t -p 5000:5000 -v "C:\..\backend\data:/data" ghcr.io/project-osrm/osrm-backend osrm-routed --algorithm mld /data/osrm/ottawa.osrm
```
