### usage
```bash
go build cmd/trees-server.go

./trees-server -h

./trees-server&

curl http://localhost:8000
curl http://localhost:8000?favoriteTree=
curl http://localhost:8000?favoriteTree=baobab

pkill trees-server # I hope you're on linux :P
```

### tests
```bash
go test  -v -coverpkg=ecosia/tree -race ./...
```

### benchmarks
```bash
go test -v -run=XXX -bench=. ./...
```

### docs
```bash
godoc -http=:6060
```
then go to [http://localhost:6060/pkg/ecosia/tree/](http://localhost:6060/pkg/ecosia/tree/).
