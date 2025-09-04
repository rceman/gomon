# gomon

GO Monitoring Service for VM. The service collects CPU, memory and disk usage of the
host and exposes a `/stats` HTTP endpoint. The endpoint returns the maximum
usage observed over the last 1 minute, 5 minutes, 1 hour and 24 hours.

By default the endpoint returns JSON structured by VM name and metric:

```
{
  "api": {
    "cpu":  {"1m": 22, "5m": 15, "1h": 5, "24h": 2},
    "mem":  {"1m": 512, "5m": 512, "1h": 256, "24h": 128},
    "disk": {"1m": 10.5, "5m": 10.5, "1h": 9.8, "24h": 9.8}
  }
}
```

Setting `output_style=short` returns a compact array representation:

```
{
  "cpu":  [max1m, max5m, max1h, max24h],
  "mem":  [max1m, max5m, max1h, max24h],
  "disk": [max1m, max5m, max1h, max24h]
}
```

The endpoint accepts an optional `output_format` query parameter. When set to
`html` it returns a simple table for each VM; otherwise JSON is returned.

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
  "mem": [..],
  "disk": [..]
}
```

Workers can forward their statistics by setting `MASTER_SEND=true` along with
`MASTER_IP`, `MASTER_PORT` and `MASTER_KEY`. Use `MASTER_SEND_INTERVAL_MIN` to
control how often workers post their stats in minutes (defaults to `0.5`, i.e.
every 30 seconds). The master aggregates all received statistics and serves
them in a single JSON object keyed by VM name at `/stats`.

