## Shortify API

Shortify API is an API for creating short urls links from long urls on the internet.

Shortify uses:
* Gorilla web toolkit [mux](http://www.gorillatoolkit.org/pkg/mux) for the routing.
* [mgo](https://labix.org/mgo) library for data model and interface with mongodb.
* [Prometheus](https://prometheus.io/) time-series database for collecting metrics/monitoring about the routes and the short links being used. Also metrics on golang threads and http routes request times.
* [prometheus-middleware](https://github.com/albertogviana/prometheus-middleware) for collecting metrics about API routes.

### Assumptions:

* Must be highly-scaleable with a focus on performance. Golang.
* Metrics are important and need to be decoupled from the database for reliability.
* The costs of getting metrics about routes should not increase the costs of using the API.
* ShortURLs should be unlimited, validated, and contain random characters.
* ShortURLs should not collide.
* Can have multiple same LongURLs that map to different ShortURLs.

### Architecture Decisions/Tradeoffs

* Golang is a very performant backend API language.
* Mux router is highly performant during lots of web requests.
* Use MongoDB as a distributed central DB (easily deployed into a Kubernetes cluster).
* Use Prometheus as the time-series metrics database.
* Docker-Compose for development and testing.

### Purpose and Implementation

The purpose is to shorten url links using a REST request to the API and then receive a redirect to the long url when accessing the short url via the API.

Routes:

|Method|Route|Data Payload|Response|
|--:|--:|--:|--:|
|POST|'http://localhost:5000/Create'|{'longurl': 'https://github.com/teamhephy/workflow/'}|{'ShortURL': 'haSd2351'}
|GET|'http://localhost:5000/${ShortURL}'| none | redirect |
|GET|'http://localhost:5000/metrics'| none | route to be scraped by prometheus for http metrics |
|GET|'http://localhost:9090'| none | Prometheus UI for Querying Metrics |

### Run the API on local:

```bash
make start TAG=latest ENV=dev
```

Builds the API and then exposed it at http://localhost:5000

### Testing the API on local:

```bash
go test ./api/tests/...
```

### Tearing down API on local:

```bash
make stop TAG=latest ENV=dev
```

### Metrics Prometheus UI Queries

![Image of Prometheus UI](../master/docs/images/prom-ui-query.png?raw=true)

### Metrics Architecture Choices (Implemented):

I started to implement using Golang and MongoDB and saw that the question also asks for metrics on short url access but I ran out of time (4 hours). There are two ways I could attempt to do this:

1. In MongoDB, I can create a new collection that will track the history as an array of the short url access with a timestamp.

Collection of:
{
  short: #{short_url}
  timestamps: [
    {accessTimestamp: ${time.Time}},
    {accessTimestamp: ${time.Time}},
    ...
  ]
}

2. The other better way would be to use Prometheus as a backend + [golang-http-metrics-prometheus](https://github.com/slok/go-http-metrics) middleware to send metrics each time a call is made to retrieve a short url. This is nicer and will scale better than the mongoDB proposed implementation above mainly because the mongodb collections will not grow so huge and metrics aggregation tools are built precisely for solving problems like these.

---------------------------------------------

Original Problem and My Quick Notes:

Programming Exercise
This exercise should be completed in 4 hours or less. The solution must be runnable, and can be written in any programming language. The challenge is to build a HTTP-based RESTful API for managing Short URLs and redirecting clients similar to bit.ly or goo.gl. Be thoughtful that the system must eventually support millions of short urls. Please include a README with documentation on how to build, and run and test the system. Clearly state all assumptions and design decisions in the README.

A Short Url:
1. Has one long url
- Object:
 {
   id: #{mongodb_id}
   short: #{short_url}  <--- short_url also becomes our object path for redirect
   long: #{long_url}
 }
2. Permanent; Once created
- No update methods allowed - only POST and GET and DELETE for now
3. Is Unique; If a long url is added twice it should result in two different short urls.
- No need to search if already exists, create a new one
4. Not easily discoverable; incrementing an already existing short url should have a low
probability of finding a working short url.
- Needs to be random string for url not just incrementing an ID for the path

Your solution must support:
1. Generating a short url from a long url
2. Redirecting a short url to a long url within 10 ms.
3. Listing the number of times a short url has been accessed in the last 24 hours, past
week and all time.
- MongoDB object history/versioning to see how many hits on a collection
4. Persistence (data must survive computer restarts)
- Use MongoDB for persistence.

Shortcuts
1. No authentication is required
2. No html, web UI is required
3. Transport/Serialization format is your choice, but the solution should be testable via curl
4. Anything left unspecified is left to your discretion
