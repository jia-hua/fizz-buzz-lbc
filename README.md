# fizz-buzz-lbc

```
"Exercise: Write a simple fizz-buzz REST server. 

The original fizz-buzz consists in writing all numbers from 1 to 100, and just replacing all multiples of 3 by "fizz", all multiples of 5 by "buzz", and all multiples of 15 by "fizzbuzz". The output would look like this: "1,2,fizz,4,buzz,fizz,7,8,fizz,buzz,11,fizz,13,14,fizzbuzz,16,...".

 

Your goal is to implement a web server that will expose a REST API endpoint that: 

Accepts five parameters : three integers int1, int2 and limit, and two strings str1 and str2.

Returns a list of strings with numbers from 1 to limit, where: all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.

 

The server needs to be:

Ready for production

Easy to maintain by other developers

- Add a statistics endpoint allowing users to know what the most frequent request has been. 

This endpoint should:

- Accept no parameter

- Return the parameters corresponding to the most used request, as well as the number of hits for this request"


```

# How to use the project

## Launch

This project requires `docker` or `go`. 

If you have docker: 

```bash
make build && make run
```

Note that for production, you will probably want to use docker to create an image, push it to your image repository and deploy it on your infrastructure from the pipeline for example. 

If you want to launch natively with the `go` command: 

```bash
go run cmd/server/main.go
```

## Query endpoint

```bash
curl "http://localhost:8080/fizzbuzz?limit=20&fizzNumber=3&buzzNumber=5&fizzString=fizz&buzzString=buzz"
```

## Metrics

Metrics are exposed through an endpoint: 

```bash
curl "http://localhost:8080/metrics"
```

They are of the `prometheus` format, so they can be scrapped by it, and visualised from `grafana` for example. Those are in the `docker-compose` (but require some conf to actually display the dashboard). 

The metrics of the fizzbuzz endpoint are a counter vector: 

```
fizzbuzz_fizzbuzz_metrics{buzz="5",buzzString="buzz",fizz="3",fizzString="fizz",limit="20"} 1
```

# Developper notes

## Architecture

The architecture of the code is inspired by `clean architecture` (which is in my opinion very similar to the concept of the classic approach of layered application, with separation of concerns). If a component needs some dependencies, it should be injected rather than being directly imported, so that it's more easily testable, and if for example, the dependency changes (ex: instead of mariaDB, you want to switch to PostgreSQL, or instead of HTTP REST, you switch to gRPC), the impact will be minimal thanks to the minimal coupling. 

Plus, the flow should be revolving around the data: you want the most basic functions as possible (ex: having functions that do pure business logic given parameter data). 

I made some adaptations/short-circuit, nothing is perfect, but it can be changed afterward as the refactor shouldn't be too cumbersome. 

```
.
├── cmd
│   └── server
├── docs
├── http
│   └── rest
├── internal
│   └── metrics
├── pkg
│   └── fizzbuzz
└── useCase
```

The idea is to have this kind of callstack (see in `docs/components.puml` for a sequence diagram): 
1. controller (inside the `http/rest`)
1. usecase (inside the `usecase`)
1. service (inside `internal` or `pkg`)
1. repository

### Controller

Nothing special, it's the entrypoint of the webservice and it map parameters to types and delegates to the corresponding use case. Note that here, normally you want to inject the use case. For the fizzbuzz use case, as this code doesn't require true dependencies (no need to call other web services, busses, databases), I directly imported it, but if I it would requires dependencies, I would inject it in a similar way as the `metric service`. 

Plus, I used the same structure `ComputeFizzBuzzRequest` for the `controller` and the `usecase`: normally you want to separate those 2 types, and having a `mapper function`, as it's 2 differents layers. I didn't did it here, for the sake of simplicity, as the use case is simple enough. Plus, if the use case is changed to be more complex in the future, it's easy to separate the types afterward when required. 

### Use case

The use case is the very reason you want to develop your service: normally, you will have as much use cases as endpoints. A use case "compose" `services` functions to achieve its goal. In this project, services are in `internal` (private services that shouldn't be imported by other projects) and `pkg` (those are public). Also, note that services should NOT call other local services: service composition should be done in usecase. If the use case is simple enough, it will map directly to the underlying service (for example, if your use case is to update a user infos, the use case will do nothing except delegating to the user service directly. But if your use case is to update a user and notify with a bus message, you will probably compose it here). 

Here the use case is to compute fizzbuzz, so I created a `fizzbuzz service`. Note that I flattened parameters here, rather than giving it an object containing the values. It's because I didn't want to create a model type and a mapper function, for the sake of simplicity. But basically, everywhere I give flattened parameters, it's to convert from external dependencies types to the model type (the `model type` is the type that should be used everywhere in the usecase/services. This type is abstracted from technical stuff. It's only when you want to convert to persistence, or json or other external dependencies for example, that you need to map to external types). 

For error handling, normally, I would map `internal error` to `controller error`: basically I want my usecases or services to throw internal errors, but the controller should call a function that map those errors to http errors. I didn't did it here, for simplicity. 

### Service

The service is responsible of its domain, and its contract abstract how it do its job. What I call a `service`, is for example, you want a function to update a user: this is a service. And to achieve it, it can call a storage webservice, or a repo for example. You can see a service as a "sub-usecase": to prepare a coffee, you need to take a mug, grind coffee beans, pour hot water, etc: every of those "sub-action" of this use case is a service. 

Note that if your webservice is a basic CRUD webservice (updateUser, createUser, getUser deleteUser), the usecase layer is not required in my opinion, as there will be no real service composition. In fact, in this case you will probably have only one "user service" calling internally its "user repo", but if there is some "high level" use case like "getUserInfosAndBillHistory", you will probably need a usecase layer, as the notion of "bill" is another domain, and you will compose user infos with its bills. 


### Repository

It's very similar to a service, except it's dedicated to call a persistence layer. The role of the repo is to abstract how you manipulate/query data from your persistence layer. Note that some times, it's hard to separate the repository of the service requiring it: some requests needs transactions, and to pass along a chain of functions the `transaction object`, with some business logic between those functions, and thus, breaking boundaries. This can be mitigated with some interfaces to wrap stuffs, doing ping/pong between the service and its repo, but at the expense of less obvious transaction, and sometimes, less optimized queries (but as this nasty problem is contained inside the service boundary, it will not propagate to the rest of the code). We don't have repo in this project as this use case doesn't require it.

## Whats missing

 - generated openAPI contract to avoid having bad contract
 - expose openAPI contract through the webservice
 - configuration file. Also, note that the main file has the default port `8080`, it's still decent, since it will be overriden with docker (but in real life app, it would still be configurable, just in case, despite being overriden later). 