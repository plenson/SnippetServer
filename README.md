SnippetService
==============
The package creates an in-memory snippet storage service.
Snippets are limited to 300 characters in size.


There is a choice between two in-memory databases which can be used, one based on BoltDB, the other on Badger.

To ease the cost of implimentation I've used a wrapper around bolt located:

```
https://github.com/xyproto/simplebolt
```

Similarily for Badger I've used a wrapper around Badger:

```
https://github.com/zippoxer/bow
```


In addition, rudimentary promethus instrumentation has been done, which outputs statistics on how many times a
route has been called. This information can be accessed at :/8000/metrics

You can try the bolt version of the service (running on AWS ECS) by quering the following :

http:/xxx.xxx.xxx.xxx:8000/snippets

http://xxx.xxx.xxx.xxx:8000/metrics

http://xxx.xxx.xxx.xxx:8000/status

see the query section below for format of url for GET,POST, DELETE

You can try the badger version of the service  (running on AWS ECS) by quering the following :

http://xxx.xxx.xxx.xxx:8000/snippets

http://xxx.xxx.xxx.xxx:8000/metrics

http://xxx.xxx.xxx.xxx:8000/status

see the query section below for format of url for GET,POST, DELETE

Note, docker container footprints are much less 15 MB

Cloud Hosting
----------------
The services above are hosted on AWS ECS using a micro instance type. Currently there is no load balancing enabled.
The docker containers used in the solutions are stored in ECR. I did attempt to build and deploy the containers via AWS code pipeline but ran into a deployment issue and so I build them locally and uploaded to ECR. Task definitions were created on 
ECS to deploy containers.

Monitoring & Metrics
--------------------
The service is instrumented (minimally) with Prometheus. In addition, I have setup CloudWatch dashboard and alarms.


Docker Build Instructions
-------------------------

Insure that your $GOPATH is setup in default configuration of $HOME/go.

```
$cd $HOME/go
$mkdir -p src/github.com/plenson
$cd src/github.com/plenson
$git clone https://github.com/plenson/SnippetServer.git

```
Code Considerations
------------------
At the moment there is a little "op smell" on how this works regarding the use of IF DEF's
to differentiate code paths for BOLT vs BOW. These can be found in the following folder/file locations:

```
1.main.go
2.routes/routes.go

To compile the code for a BoltDB database search for #IF BUILD BOLT ... immediately below it uncomment the line.

To compile the code for a Badger database search for #IF BUILD BOW ... immediately below it uncomment the line.

Obviously make sure you don't have both uncommented at the same time.
After these changes you're ready to build the docker images.

Note,I hope to reduce this "op smell shortly"
```

Build
------

```
Bolt Build
$ docker build -f DockerfileBolt -t  golang-docker-snippet-service-bolt  .

Bow Build
$ docker build -f DockerfileBow -t  golang-docker-snippet-service-bow  .

```

Run (Default settings)
-----------------------

```
$ docker volume create --name Data

Bolt
$ docker run -d -it -p 8000:8000 -v /Data golang-docker-snippet-service-bolt

Bow
$ docker run -d -it -p 8000:8000 -v /Data golang-docker-snippet-service-bow

```

The server will listen on http://localhost:8000

Query
------
To query use your favorite tool ... Postman is my choice.
In the postman interface try the following queries:

Get All Snippets
-----------------

```
 [Get]   http://localhost:8000/snippets
```

 You should see something like:

 ```
 1 {"items":2}
 2 {"id":"cd2a6c-xxxxxxxxxxxxxxx","text":"Little lamb","shared":true}
 3 {"id":"455353-xxxxxxxxxxxxxxx","text":"Mary had a little lamb.","shared":true}
```

Get A Specific Snippet
------------------------

```
 [Get]   http://localhost:8000/snippet/cd2a6c-xxxxxxxxxxxxxxx
```

 You should see something like:

```
 1 {"id":"cddfadasc-xxxxxxxxxxxxxxx","shared":true}
```

Delete A Specific Snippet
------------------------

```
 [DELETE]   http://localhost:8000/snippet/cd2a6c-xxxxxxxxxxxxxxx
```

 You should see something like:

```
 1 {"id":"cddfadasc-xxxxxxxxxxxxxxx","shared":true}
```

If you run the Get All Snippets above again you should see the  snippet has been deleted.

Create A Specific Snippet
------------------------

```
 [POST]   http://localhost:8000/snippet/
```

In the body portion add you snippet in a json format:

```
{
  "Text":"My Snippet text"
}
```

You should see something like:

```
 1 {"id":"cddfadasc-xxxxxxxxxxxxxxx","shared":true}
```

If you run the Get All Snippets above again you should see the new snippet has been added.

Check Status of Service
------------------------

```
 [Get]   http://localhost:8000/status
```

 You should see something like:

```

Api is up and running

```

Get Metrics
------------------------

```
 [Get]   http://localhost:8000/metrics
```

You should see something like:

```
# HELP del_handler_total Del Handler requested.
# TYPE del_handler_total counter
del_handler_total 0
# HELP get_handler_total Get Handler requested.
# TYPE get_handler_total counter
get_handler_total 0
# HELP getall_handler_total GetAll Handler requested.
# TYPE getall_handler_total counter
getall_handler_total 102003
# HELP go_gc_duration_seconds A summary of the GC invocation durations.
# TYPE go_gc_duration_seconds summary
go_gc_duration_seconds{quantile="0"} 0
go_gc_duration_seconds{quantile="0.25"} 0
go_gc_duration_seconds{quantile="0.5"} 0
go_gc_duration_seconds{quantile="0.75"} 0
go_gc_duration_seconds{quantile="1"} 0.0899937
go_gc_duration_seconds_sum 0.1130529
go_gc_duration_seconds_count 126
# HELP go_goroutines Number of goroutines that currently exist.
.
.
.
```

TODO
----

Unit test


API-GATEWAY

I've looked at setting up an API gateway in front of the service allowing me to vary the backend.


Sharding

The database could be sharded allowing it to scale out.


