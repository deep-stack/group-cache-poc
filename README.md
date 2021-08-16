# Group Cache Proof of Concept

## Usage

* Ensure `Go 1.16` is installed
* Ensure `Postgres 12` is installed 
* Run
```bash
sudo -u postgres createdb group_cache_poc
```
* Update password of `postgres` user in config.json
* Run 
```bash
go build main.go
```
* Run
```bash
./main -addr=:8080 -pool=http://localhost:8080,http://localhost:8081 &
```
* Run
```bash
./main -addr=:8081 -pool=http://localhost:8081,http://localhost:8080 &
```
* Database Tables will be created automatically
* 2 HTTP servers will start on ports 8080 and 8081

## RESTful API
* get          : GET `kv/get?key=k`
* set (upsert) : POST `kv/set?key=k&value=v`

## Testing
```bash
curl -X POST "localhost:8080/kv/set?key=foo&value=bar"
```

```bash
curl "localhost:8081/kv/get?key=foo"
```

## References
* https://www.mailgun.com/blog/golangs-superior-cache-solution-memcached-redis
* https://github.com/mailgun/groupcache
* https://www.alexedwards.net/blog/organising-database-access
