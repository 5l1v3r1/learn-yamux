example for yamux + serial + grpc
=================================

# use case
```
.
├── basic
│   ├── grpc             # basic example for grpc
│   ├── serial           # basic example for serial
│   └── yamux            # basic example for yamux
├── exec_serial          # remote exec command over serial
├── yamux_serial         # example for yamux over serial
└── yamux_grpc_serial    # example for yamux + grpc over serial
```

# manage vendor
```
$ govendor init
$ govendor add +external
```