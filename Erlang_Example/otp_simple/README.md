```
erlc server.erl name_server.erl
```

```
erl

server:start(name_server, name_server).
name_server:add(erlang,"foo").
name_server.whereis(erlang).
```
