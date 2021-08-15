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

## RESTful API
* get          : GET `/Get?key=k`
* set (upsert) : POST `/Set?key=k&value=v`

## Testing
* `curl -X POST "localhost:8080/Set?key=foo&value=bar"`
* `curl "localhost:8081/Get?key=foo"`

## References
* https://www.mailgun.com/blog/golangs-superior-cache-solution-memcached-redis
* https://github.com/mailgun/groupcache
* https://www.alexedwards.net/blog/organising-database-access
