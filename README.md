# trace

infrago trace module with multi-connection fanout.

## config

```toml
[trace]
driver = "default"
json = true
buffer = 1024
timeout = "200ms"
sample = 1.0
format = "%time% [%status%] %name% trace=%traceId% span=%spanId% cost=%costMs%ms"
```

or multi-write:

```toml
[trace.file]
driver = "default"

[trace.greptime]
driver = "greptime"
fields = { trace_id = "tid", span_id = "sid", parent_span_id = "psid", timestamp = "ts" }
[trace.greptime.setting]
host = "127.0.0.1"
port = 4001
```

`fields` should be configured on `trace.<conn>.fields` (not under `setting`).
