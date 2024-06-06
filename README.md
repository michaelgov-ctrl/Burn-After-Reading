# burn_after_reading
This is a simple http web app for sharing messages stored in memory that are deleted after being accessed.

```
endpoints:
    GET  /v1/healthcheck    - generic healthcheck handler
    GET  /v1/messages/:uuid - retrieve message with given uuid
    POST /v1/messages       - post message content
    GET  /metrics           - prometheus metrics endpoint
```


the api accepts a json message 
```
message posting schema is {"content": "$message"} 
return message schema is {"uuid": "$uuid", "content": "$message"}
```



a posted message returns the uuid for the message and a location header to access the message
```
    PowerShell:
        POST:
            Invoke-WebRequest -Uri http://$host:$port/v1/messages -Method Post -Body '{"content": "hello"}'

        GET:
            Invoke-WebRequest -Uri http://$host:$port/v1/messages?uuid=$uuid

        Example:
            $resp = Invoke-WebRequest -Uri http://$host:$port/v1/messages -Method Post -Body '{"content": "hello"}'
            Invoke-WebRequest -Uri ("http://$host:$port"+$resp.Headers.Location)
```


```
    Bash:
        POST:
            curl -d '{"content": "hello"}' -X POST http://$host:$port/v1/messages

        GET:
            curl -i -H "Accept: application/json" http://$host:$port/v1/messages?uuid=$uuid
```

setup for code changes
```
cd existing_repo
git remote add origin https://gitlab.clarkinc.biz/mgovernanti/burn_after_reading.git
git branch -M main
git push -uf origin main
```