# Grill
---

Grill your application.

## Motivation
* * *
1. Reduce the overload of writing Functional Tests at the same time better their quality.
2. Functional Tests should test the behaviour of your application, without any knowledge of its internals(Behaviour Driven Testing/Black Box Testing).
3. Functional Tests should be declarative, so its easy to read and understand them.
4. Provide an easy way to mock external dependencies and setup infra components.


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

2) Grill has it own test case runner which takes in a slice of testcases and runs them.
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
```
TODO
```

## Why write functional tests
* * *
Ans:
![umbrella](https://media.tenor.com/images/74be340020f6b91b66065b51abae7a76/tenor.gif)

