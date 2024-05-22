include pkg/common/envs/.env


startDatabase:
     docker-compose up -d

stopDatabase:
     docker-compose stop

migrateUP:
      migrate -path pkg/common/db/migration/ -database $(POSTGRES_CONNECTION_STRING) -verbose up

migrateDOWN:
      migrate -path pkg/common/db/migration/ -database $(POSTGRES_CONNECTION_STRING) -verbose down

sqlcGenerate:
     sqlc generate
     
gqlGenerate:
     go run -mod=mod github.com/99designs/gqlgen generate   
server:
     go run cmd/main.go


test:
     go test ./pkg/test -v