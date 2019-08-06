# pingtraq

A small daemon and cli that registers "pings" from devices of your choice.

## ptd

```sh
# Start server
$ ptd -d foo.db -p 8080
```

## ptctl

```sh
# Add ping endpoint
$ ptctl add foo -d foo.db

# List ping endpoints
$ ptctl list -d foo.db

# Get records from endpoint
$ ptcl get foo -d foo.db
```
