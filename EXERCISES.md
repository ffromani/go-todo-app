# Go testing bootcamp exercises

Note: there are some failing tests around that need to be fixed. We suggest running only the test that is
part of the exercise by running `go test -run TestName` in order to avoid unnecessary noise.

## Subtests / Table tests

- Modify the `TestMerge` test under [model/todo_merge_test.go](model/todo_merge_test.go) to be a table based
test. 
- Try to have shared code initializing one `Todo` item to be merged
- Try to run make the test execution parallel
- Run `go test -v` to see all the tests (and their names)
- Try to use `go test -run testname` to run one single test

Examples: [Sub tests and table tests examples](https://github.com/fedepaol/gotestbootcamp/tree/main/subteststabletests)


## Test Fixtures / Golden files

- Try to run the rendering tests under [model/todo_render_test.go](model/todo_render_test.go) with `go test -run TestRender`
- Let the test generate the golden files with `go test -run TestRender -update`
- Verify that the generated golden file is what we expect
- Check if the test is now passing by matching the golden file `go test -run TestRender`
- Extend the tests and verify the generated files are valid

Examples: [Fixtures and golden files examples](https://github.com/fedepaol/gotestbootcamp/tree/main/fixturesandgoldenfiles)

## Benchmarking

- Run the benchmark test under [controller/controller_bench_test.go](controller/controller_bench_test.go) with `go test -run xx -bench . -benchmem` (note the run xx to avoid running the tests too).
- Replace the implementation being instrumented the the one using a Reader (`todoFromRequestReader`) and benchmark again
- Install benchstat with `go install golang.org/x/perf/cmd/benchstat@latest`
- Validate the difference between the two approaches:

```bash
go test -run xx -bench . -benchmem -count 10 > withreader.txt
# replace the implementation back
go test -run xx -bench . -benchmem -count 10 > withstrings.txt
benchstat withstrings.txt withreader.txt
```

Examples under [https://github.com/fedepaol/gotestbootcamp/tree/main/benchmarking](https://github.com/fedepaol/gotestbootcamp/tree/main/benchmarking)

## Dependency injection

### Inject a function

Change the code so that [model/todo_dep_test.go](model/todo_dep_test.go) is deterministic (hint: the `New(title string)` function has a dependency from the `time` package that can be replaced / injected)

### Inject a field of an object

The `TestTodoCreate` function under [controller/controller_dep_test.go](controller/controller_dep_test.go) is not testing the returned uuid. This is because fetching the uuid requires the interaction with a remote service.

- Make the test validate uuid
- Change the code so that uuid can be deterministic
- Add a negative test (when the call to the service fails)
- Let the test validate that the uuid is requested only once per call

Examples under [https://github.com/fedepaol/gotestbootcamp/tree/main/dependency_injection](https://github.com/fedepaol/gotestbootcamp/tree/main/dependency_injection)

## Integration with Http servers

- Modify the code so that the result of [uuid/uuid_test.go](uuid_test.go) is deterministic.
- Change the test so that:
    - it verifies that the endpoint is called only once per `NewUUID` call
    - it verifies it can handle uuids of len=3 and len=10
    - it verifies that when the call returns a failure, the function returns an error

Examples under [https://github.com/fedepaol/gotestbootcamp/tree/main/httpserver](https://github.com/fedepaol/gotestbootcamp/tree/main/httpserver)

## Docker Tests

- Extend the `TestWithRedis` test under [store/redis_test.go](store/redis_test.go)
- Fill the emtpy tests
- Add more tests 

