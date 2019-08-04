example for yamux + serial + grpc
=================================

# use case
```
.
├── basic
│   ├── grpc             # basic example for grpc
│   ├── serial           # basic example for serial
│   └── yamux            # basic example for yamux
├── yamux_serial         # example for yamux + serial
└── yamux_serial_grpc    # example for yamux + serial + grpc(WIP)
```

# manage vendor
```
$ govendor init
$ govendor add +external
```