go-todo-app
===========

A simplistic demo web application in golang to manage TODO list using a JSON-RPC API.
This is a demo project/learning aid! Don't use anywhere near production.

Architecture
------------


./go-todo-app
```
├── api          types used in the public API layer, to decouple from the internal representation
│   └── v1       current version
├── cmd          app entry point. Keep minimal!
├── config       configuration processing, from flags, files...
├── controller   orchestration layer, decodes/encodes object from API, manipulates internal objects
├── ledger       high level data store, deals with objects (e.g. Todo)
├── middleware   utilities to inject in the HTTP handling to augment it
├── model        internal data types definitions, including their operations
└── store        durable data store, bytestream oriented
    └── fake     fake, non durable, data store to be used in testing
```

Please look at godocs of packages, functions, types for more details

Limitations
-----------

Due to the demo nature of the project:
- The JSON-RPC is minimal, because this project is meant for demo purposes
- The routes are not very REST-ish nor especially clean
- bytestream encoding is not versioned
- No object is thread safe (no locking).

License
-------

MIT
