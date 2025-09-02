# gomon

GO Monitoring Service for VM. The service collects CPU and memory usage of the
host and exposes a `/stats` HTTP endpoint. The endpoint returns the maximum CPU
percentage and memory usage (in MB) observed over the last 5 minutes, 1 hour and
24 hours in the following format:

```
{
  "cpu": [max5m, max1h, max24h],
  "mem": [max5m, max1h, max24h]
}
```

Set the polling interval (seconds, decimals allowed) with `READ_TICKER_TIME_SEC` and the server port with
`STATS_PORT` in the environment. Each instance must define `VM_NAME`.

## Master/worker mode

If `MASTER_NODE=true`, the instance accepts `POST /stats` requests from other
nodes. Requests must include `Authorization: Basic <base64(MASTER_KEY)>` and a
JSON body:

```
{
  "name": "RNG",
  "cpu": [..],
  "mem": [..]
}
```

Workers can forward their statistics by setting `MASTER_SEND=true` along with
`MASTER_IP`, `MASTER_PORT` and `MASTER_KEY`. Use `MASTER_SEND_INTERVAL_MIN` to
control how often workers post their stats in minutes (defaults to `0.5`, i.e.
every 30 seconds). The master aggregates all received statistics and serves
them in a single JSON object keyed by VM name at `/stats`.

