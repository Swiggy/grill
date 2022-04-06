**Grill**
---
---
Grill your application.

## Motivation
* * *
* Reduce the overload of writing Functional Tests at the same time better their quality.
* Functional Tests should test the behaviour of your application, without any knowledge of its internals (Behaviour Driven/ Black Box Tests).
* Functional Tests should be declarative, so it's easy to read, understand and review them.
* Provide an easy way to mock external dependencies and setup infra components.


## How it Works
* * *
* It defines a testcase using a list of stubs, assertions and cleaners and an Action method to invoke the public api of your application.
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

* Grill has it own test case runner which takes in a slice of testcases and runs them.
```	
tests := []grill.TestCases{}
grill.Run(t, tests)
```

## Features
* * *
* Grill provides Inbuilt helpers(stubs,assertions,cleaners) and initializers for most of the infra dependencies we use.
* For external/upstream/downstream services it provides mocking utilities for http and grpc.


| Grill                    | Stubs                                  | Assertions  | Cleaners                     |
|--------------------------|----------------------------------------|---|------------------------------|
| HTTP (wiremock)          | Stub, StubFromJSON, StubFromFile       | AssertCount  | ResetAllStubs                |
| GRPC                     | Stub                                   | AssertCount | ResetAllStubs                |
| DynamoDB                 | CreateTable, SeedDataFromFile, PutItem | AssertScanCount, AssertItem  | DeleteTable, DeleteItem      |
| Kafka                    | CreateTopics                           | AssertCount, AssertMessageCount | DeleteTopics                 |
| Redis                    | SelectDB, Set                          | AssertValue | FlushDB                      |
| Mysql                    | CreateTable, SeedFromCSVFile           | AssertCount | DeleteTable                  |
| S3 (minio)               | CreateBucket, UploadFile               | AssertFileExists | DeleteBucket, DeleteAllFiles |
| Tile38                   | SetObject                              | AssertObject  | FlushDB                      |
| Consul                   | SeedFromCSVFile, Set                   | AssertValue | DeleteAllKeys                |
| SQS                      | CreateQueue                            | AssertCount, AssertMessageCount | DeleteQueues                 |


 
## Usage 
* * *
##### Download
```
go get github.com/swiggy-private/grill
```

##### Actions

* Use Actions to call the public API of your application. eg - http request, grpc method, kafka produce.
* Return the output to assert on.
* Use `grill.ActionOutput(out ...interface)` to return multiple outputs.
* Grill has an in build assertion `grill.AssertOutput(outputs ...interface)` to compare the output using reflect.DeepEqual. Use grill.Any to skip a particular field.

```
action := func() interface{} {
    res, err := http.Get("http://www.google.com")
    return grill.ActionOutput(res, err)
}

grill.AssertOutput(grill.Any, nil)
```

##### Starting/Stopping a Grill
* To Start/Stop one Grill 
```
grl := grillHTTP.HTTP{}
err := grl.Start(context.Background())
err := grl.Stop(context.Background())
``` 
* To Start/Stop Multiple Grills 
```
grills := Grills{
    HTTP:   &grillHTTP.HTTP{},
    GRPC:   &grillGRPC.GRPC{},
    Dynamo: &grillDynamo.Dynamo{},
}
err := grill.StartAll(context.Background(), grills.HTTP, grills.GRPC, grills.Dynamo)
err := grill.StopAll(context.Background(), grills.HTTP, grills.GRPC, grills.Dynamo)
``` 

##### Writing Tests
```
testStub = grillhttp.Stub{
    Request: grillhttp.Request{Method:"GET",UrlPath:"/test"},
    Response: grillhttp.Response{Status: 200,Body: "PASS"},
}

tests := []grill.TestCase{
    {
        Name: "TestHTTPStub",
        Stubs: []grill.Stub{
            grl.Stub(&testStub),
        },
        Action: func() interface{} {
            res, err := http.Get(fmt.Sprintf("http://localhost:%s/test", grl.Port()))
            if err != nil {
                return err
            }
            if res == nil || res.Body == nil {
                return nil
            }
            defer res.Body.Close()
            got, _ := ioutil.ReadAll(res.Body)
    
            return grill.ActionOutput(string(got), res.StatusCode, err)
        },
        Assertions: []grill.Assertion{
            grill.AssertOutput("PASS", http.StatusOK, nil),
            grl.AssertCount(&testStub.Request, 1),
            
            // Check Items in Database, messages in kafka topics, dp events etc.
        },
        Cleaners: []grill.Cleaner{
            grl.ResetAllStubs(),
        },
    },
}
```
##### Running Tests
* To run a single test use `Run()` method on the testcase.
* To run multiple tests use `grill.Run(t, tests)`.
* To run tests in parallel use `grill.RunParallel(t, tests)`. Only use this if your tests don't share state. 
```
test := grill.TestCase{}
test.Run(t)

tests := []grill.TestCase{test, test, test}
grill.Run(t, tests) 
OR
grill.RunParallel(t, tests)
```


##### Testing Async Flows
* Use `grill.Try(deadline, minSuccess, assertion)` method to test async flows, like kafka, dp etc.
* It fails if an assertion is not successful minSuccess times in the given deadline.
* As a best practice keep minSuccess > 1 to make sure the assertion didn't succeed in an intermediate state.
```
grill.Try(time.Second, 3, grill.AssertOutput("PASS", http.StatusOK, nil))
```

##### Testing Negative Assertions
* To test assertions which should fail, wrap them using `grill.Not(assertion)`.
```
grill.Not(grill.AssertOutput("PASS", http.StatusOK, nil))
```

##### Writing Custom Stubs, Assertions and Cleaners
* Implement the Interface.
* Wrap a function using `grill.StubFunc(fn)`, `grill.AssertionFunc(fn)` or `grill.CleanerFunc(fn)`.
```
grill.AssertFunc(func() error {
    return fmt.Errorf("i will always fail")
}) 

```
## Why write functional tests at all ??
* * *
Ans:
![umbrella](https://media.tenor.com/images/74be340020f6b91b66065b51abae7a76/tenor.gif)

