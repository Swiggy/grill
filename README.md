# Grill - A test framework for functional validation of golang microservices
## Background
The behavior of any system should be tested only from the user's perspective, as if its done without any knowledge of its internal implementation, i.e. the System Under Test(SUT) should be a black box for the test. In the test we call the public API of the SUT, validate the response and any outgoing requests/messages.

This decouples the test from the actual system and allows us to change the implementation without any change in the test cases.

All tests, be it unit tests, service level tests or end to end integration tests, should follow the same methodology, only the SUT changes in all cases.

![test](https://github.com/Swiggy/grill/blob/media/testing.png?raw=true)

## Motivation
With rapid adoption of microservices architecture at Swiggy, we realized the limitation of unit tests. Especially for I/O bound microservice which 
typically engage in data transformation, interaction with downstream dependencies and infrastructural components (e.g. database, queues, 
cache etc.) to implement business logic, unit tests have limited value. It was more important to test the API surface of the microservice 
with all of its infrastructure components such that the complete behaviour of the microservice can be validated. In this type of test, the 
microservices along with its infrastructure dependencies in installed locally and the test client invokes the public APIs and validates output via 
both API response and via storage layer validation(e.g. checking the output of a queue). We coined a term for this type of test - Service Level Testing (SLT). 

While individual microservices could develop their own test framework along with mocking, container management for infra components, we developed a utility 
that drastically reduced the cost of writing SLT test cases. Once the tool was widely adopted within Swiggy, we decided to open source 
this tool to give back to the community and solicit more contributions.

## Grill
Grill is a SLT framework which extends the above principle and provides a declarative way for writing service level tests. The core framework defines a TestCase struct 
with interfaces for stubs, assertions and cleaners.

```
type Stub interface {
	Stub() error
}

type Assertion interface {
	Assert() error
}

type Cleaner interface {
	Clean() error
}

type TestCase struct {
	Name       string
	Stubs      []Stub              // Setup all mocks, stubs etc. Inject them as required. 
	Action     func() interface{}  // Call the public API of the SUT (http call, kafka publish etc) 
	Assertions []Assertion         // Validate the response received, outgoing requests made. 
	Cleaners   []Cleaner           // Undo all the mocks and reset the state for the next test.

}
```

The core framework can be used for any type of tests, we just have to provide required implementations of stubs, assertions and cleaners.

But since the library was built for aiding in service level tests, it provides helpers(stubs, assertions, cleaners implementations) for commonly used infra dependencies like dynamodb, redis, grpc/http downstreams, kafka etc out of the box.

It uses testcontainers-go underneath for dockers.

## Supported Features

| Grill           | Stubs                                 | Assertions  | Cleaners                     |
|-----------------|---------------------------------------|---|------------------------------|
| HTTP (wiremock) | Stub, StubFromJSON, StubFromFile      | AssertCount  | ResetAllStubs                |
| GRPC            | Stub                                  | AssertCount | ResetAllStubs                |
| DynamoDB        | CreateTable, SeedDataFromFile, PutItem | AssertScanCount, AssertItem  | DeleteTable, DeleteItem      |
| Kafka           | CreateTopics                          | AssertCount, AssertMessageCount | DeleteTopics                 |
| Redis           | SelectDB, Set                         | AssertValue | FlushDB                      |
| Clustered Redis | Set                          | AssertValue | FlushDB                      |
| Mysql           | CreateTable, SeedFromCSVFile          | AssertCount | DeleteTable                  |
| S3 (minio)      | CreateBucket, UploadFile              | AssertFileExists | DeleteBucket, DeleteAllFiles |
| Tile38          | SetObject                             | AssertObject  | FlushDB                      |
| Consul          | SeedFromCSVFile, Set                  | AssertValue | DeleteAllKeys                |
| SQS             | CreateQueue                           | AssertCount, AssertMessageCount | DeleteQueues                 |

## Getting Started
Check the [Wiki page](https://github.com/Swiggy/grill/wiki) for getting started with Grill

## Why write functional tests at all ??
* * *
Ans:
![umbrella](https://media.tenor.com/images/74be340020f6b91b66065b51abae7a76/tenor.gif)
