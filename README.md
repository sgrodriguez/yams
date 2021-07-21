# yams
## Yet Another Mock Server
### Motivation
I want to use something similar to [Postman mock server](https://learning.postman.com/docs/designing-and-developing-your-api/mocking-data/setting-up-mock/) but self hosted and without sending data to thirds party nor depends on their infra.

### Featurses

* In memory
* Load Mocks through config or in runtime via POST /mock 
* Delay Execution of mocks
* Small

### Installation

```
go install github.com/sgrodriguez/yams@master;
yams config.toml
```
or 
```
git clone github.com/sgrodriguez/yams;
make
```

### How To Use

You can load mocks via toml conf
```toml
[[Mocks]]
  Name = "example-foo"
  DelayInSeconds = 0
  [Mocks.Request]
    Resource = "/foo"
    Method = "GET"
  [Mocks.Response]
    HTTPStatusCode = 200
    Body = "Hi Foo"
[[Mocks]]
  Name = "example-var"
  DelayInSeconds = 5
  [Mocks.Request]
    Resource = "/foo/var"
    Method = "PUT"
  [Mocks.Response]
    HTTPStatusCode = 500
    Body = "Bye Bar"
```
or in run time
```
curl --location --request POST 'localhost:8080/mock' \
--header 'Content-Type: application/json' \
--data-raw '{
    "request": {
        "resource": "/foo",
        "method": "PUT"
    },
    "response": {
        "body": "{\"money\":\"us\"}",
        "http_status_code": 404
    },
    "delay_in_seconds": 0,
    "name": "money"
}'
```
Use the Foo mock resource
```
curl --location --request PUT 'localhost:8080/foo'
```

List all available mocks
```
curl --location --request GET 'localhost:8080/mocks'
```
#### Config
[Example_Config](https://github.com/sgrodriguez/yams/example_config.toml)


