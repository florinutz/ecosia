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
