Hi there :wave: Thank you for your interest in Ecosia and for taking the time to solve this assignment.

We understand that assignments can be time consuming and because of that can be unfair to candidates that have less time to spare. Therefore, we have created our assignment so that you ideally shouldn't spend more than an hour. Feel free to make some compromises to keep it short and document further work/improvements you would add.

## The assigment

We ask you to create a simple web server, using only the standard library packages of your choice of python, go or nodejs. It should do the following:

* Runs locally on port 8000 and accepts `GET` requests at the index URL `/`
* It checks that the request has a query parameter called `favoriteTree` with a valid value
* For a successful request, returns a properly encoded HTML document with the following content:

If `favoriteTree` was specified (e.g. a call like `127.0.0.1:8000/?favoriteTree=baobab`):

```
It's nice to know that your favorite tree is a <value of "favoriteTree" from the url> 
```

if not specified (e.g. a call like `127.0.0.1:8000/`):

```
Please tell me your favorite tree
```

## Submission

Please send us an archive (zip, tar, ...) that contains everything you created for the assignment, like the code itself, documentation, tests etc.

We hope you'll enjoy working on this and good luck! :four_leaf_clover:

# Solution

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
