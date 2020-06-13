# Grill
---

Make your applications better.

## Motivation
* * *
1. Reduce the overload of writing Functional Tests at the same time better their quality.
2. Functional Tests test the behaviour of your application, and hence should be declarative in nature.
3. An easy way to mock external dependencies and setup infra components.


## How it Works
* * *
1) It defines a testcase using a list of stubs, assertions and cleaners and an Action method to invoke the public api of your application.
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
	Stubs      []Stub
	Action     func() interface{}
	Assertions []Assertion
	Cleaners   []Cleaner
}
```

2) Grill has it own test case runner which takes in an array of testcases and runs them.
```	
tests := []grill.TestCases{}
grill.Run(tests)
```

## Features
* * *
Grill provides Inbuilt helpers(stubs,assertions,cleaners) and initializers for most of the infra dependencies we use in swiggy. For external services it provides mocking utilities for http and grpc.

General Helpers -

 - [ ] Wiremock(HTTP)
 - [ ] GRPC
 - [ ] DynamoDB
 - [ ] Kafka
 - [ ] Redis
 - [ ] Mysql
 - [ ] Tile38
 - [ ] S3
 - [ ] Consul
 - [ ] Prometheus

 Swiggy  Helpers -
 
 - [ ] Data Platform
 - [ ] Experimentation platform
 
## Usage 
* * *
```
go get bitbucket.org/swigy/grill
```


## Why write functional tests
* * *
Ans:
![umbrella](https://i.makeagif.com/media/10-03-2015/aW9A9X.gif)

