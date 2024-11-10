# Go testing bootcamp exercises

Note: there are some failing tests around that need to be fixed. We suggest running only the test that is
part of the exercise by running `go test -run TestName` in order to avoid unnecessary noise.

Note: there are not universal solutions. The proposed one are illustrative. There are many possible correct solutions.

## Basics and coverage

- Write one or more (the more the better) unit tests. If unsure where to start: `api/v1` or `config` can be quite simple;
  for more complex scenarios: `middleware`.
- Check the coverage before/after the tests. You can use the `make cover-vew` makefile target.
- Write one or more integration tests, e.g. tests which involve two or more packages. If unsure where to start: `ledger`
  is probably the simplest case, `controller` provides more (and more complex) opportunities.
- Check how the coverage changed before/after the integration tests.
- Compile the test binaries, and then run them.

Solutions: please checkout the `part1_solution` branch and `git grep solution:part1`


## Subtests / Table tests

- Modify the `TestMerge` test under [model/todo_merge_test.go](model/todo_merge_test.go) to be a table based
test. 
- Try to have shared code initializing one `Todo` item to be merged
- Try to make the test execution parallel
- Run `go test -v -run TestMerge` to see all the tests (and their names)
- Try to use `go test -run TestMerge/xxx` to run one single test

Examples: [Sub tests and table tests examples](https://github.com/fedepaol/gotestbootcamp/tree/main/subteststabletests)


## Test Fixtures / Golden files

- Try to run the rendering tests under [model/todo_render_test.go](model/todo_render_test.go) with `go test -run TestRender`
- Let the test generate the golden files with `go test -run TestRender -update`
- Verify that the generated golden file is what we expect
- Check if the test is now passing by matching the golden file `go test -run TestRender`
- Extend the tests and verify the generated files are valid

Examples: [Fixtures and golden files examples](https://github.com/fedepaol/gotestbootcamp/tree/main/fixturesandgoldenfiles)

## Using ginkgo

- Bootstrap a ginkgo suite:
```
go install github.com/onsi/ginkgo/v2/ginkgo
go get github.com/onsi/gomega/
go mod tidy
mkdir e2e && cd e2e
ginkgo bootstrap .
```
- Run tests using ginkgo: `ginkgo -v ./e2e/...`
- Write and run one or more e2e tests using ginkgo. You can either assume the `go-todo-app` is running, or run it as part of the test suite. Evaluate the pros and cons of each approach.
- Integrate gingko custom matcher in the e2e test(s) you wrote
- Compile the e2e test binary which uses ginkgo, and run it

Solutions: please checkout the `part3_solution` branch and `git grep solution:part3`

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

- Check the [uuid/uuid_test.go](uuid/uuid_test.go) and see why it fails
- Modify the code so that the result of [uuid/uuid_test.go](uuid/uuid_test.go) is deterministic.
- Change the test so that:
    - it verifies that the endpoint is called only once per `NewUUID` call
    - it verifies it can handle uuids of len=3 and len=10
    - it verifies that when the call returns a failure, the function returns an error

Examples under [https://github.com/fedepaol/gotestbootcamp/tree/main/httpserver](https://github.com/fedepaol/gotestbootcamp/tree/main/httpserver)

## Docker Tests

- Extend the `TestWithRedis` test under [store/redis_test.go](store/redis_test.go)
- Fill the emtpy tests
- Add more tests 

# Extras / Stretch goals

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

## Enhancing go testing

- Integrate `go-cmp` in the tests used previously. Make sure to use `cmp.Diff` and `cmp.Equal` Good candidates can be tests for `model` or for `api/v1` (`go get github.com/google/go-cmp`)
- Rewrite existing unit tests to use `testify/assert` (`go get github.com/stretchr/testify`)


