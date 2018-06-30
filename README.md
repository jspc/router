# Router

Router provides a proxy to docker containers exposing http/ https endpoints. It receives a request, picks out the `Host` header and tries to find a container matching that name. If it can find one it'll look for a set of env vars to determine how to proxy to that container.


## Usage

There are two environment variables a container can set for router to be able to route to it.

| Variable        | Purpose                         | Optional | Default |
|-----------------|---------------------------------|----------|---------|
| `ROUTER_PORT`   | Port on which to direct traffic | `false`  |         |
| `ROUTER_SCHEME` | URL Scheme to connect with      | `true`   | `http`  |


```bash
$ docker run -d -P -e ROUTER_PORT=80 --name my-site nginx
$ docker run -e DOCKER_API_VERSION=1.37 -v /var/run/docker.sock:/var/run/docker.sock -p8080:8080 -ti jspc/router
```

The nginx container can now be accessed as per:

```bash
$ curl -v -H 'Host: my-site.internal' http://localhost:8080                                                                       * Rebuilt URL to: http://localhost:8080/
*   Trying ::1...
* TCP_NODELAY set
* Connected to localhost (::1) port 8080 (#0)
> GET / HTTP/1.1
> Host: my-site.internal
> User-Agent: curl/7.54.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Accept-Ranges: bytes
< Connection: keep-alive
< Content-Length: 612
< Content-Type: text/html
< Date: Sat, 30 Jun 2018 21:34:01 GMT
< Etag: "5b167b52-264"
< Last-Modified: Tue, 05 Jun 2018 12:00:18 GMT
< Server: nginx/1.15.0
<
<!DOCTYPE html>
<html>
<head>
<title>Welcome to nginx!</title>
```
