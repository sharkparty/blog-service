### Initializing a local DB

You will need a mongodb instance running on the default port with a db called `mydb` 
and a collection called `blog`:

```zsh
$ docker run --name mongo-blog -d -p 27017:27017 mongo
$ docker exec -it mongo-blog mongo
$ use mydb
$ db.createCollection("blog")
```

### Testing

In order to run the tests that are VERY MUCH not unit test, run the API `make serve`
and then run the test suite `go test ./server/server_test.go`