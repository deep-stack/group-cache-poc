# group-cache-poc

## Usage

1. Make sure you have Go and Postgres 12 installed
2. Run `sudo -u postgres createdb grp_cache`
3. Update password of postgres user in config.json
4. Run `go build main.go`
5. Run `./main -addr=:8080 -pool=http://localhost:8080,http://localhost:8081 &`
6. Run `./main -addr=:8081 -pool=http://localhost:8081,http://localhost:8080 &`
7. SQL Tables will be created automatically
8. 2 HTTP servers will start on 8080 and 8081

## REST API
* get          : GET `/get/{key}`
* set (upsert) : POST `/set/{key}/{value}`

## Testing
* `curl -X POST localhost:8080/set/{key}/{value}`
* `curl -X GET localhost:8080/get/{key}`

## References
* https://www.mailgun.com/blog/golangs-superior-cache-solution-memcached-redis/
* https://github.com/mailgun/groupcache
* https://www.alexedwards.net/blog/organising-database-access
