Reporter
========

Reporter takes current state reports from components on a machine and combines them into a single state tree.

See the [historian repo](https://github.com/FuseRobotics/Historian) for documentation on how this works with Historian.

Mechanism
=========

Reporter accepts state reports from components on the local system. Each component usually maps to a docker container.

This is somewhat similar to a ROS topic, but is designed to be aggregated and streamed to upstream storage locations.

As an example, here is a possible tree of state reports:

```
├── flight_controller
│   ├── pose
│   └── target
└── sensor_rx
    ├── rx_state
    ├── sensor_123
    ├── sensor_221
    └── sensor_xxx
```

Where `sensor_rx.sensor_233` contains:

```json
{
  "rx_timestamp": "1475439400",
  "temperature": "30.0"
}
```

State objects can be any valid JSON object, which is an object: `{}` or `null`.

A component can set a state object. Reporter will determine what changed from the last version of the state, and record a time-series transaction log of changes to the state. It will then, depending on configuration, attempt to report this transaction log to one or more remotes. It will keep track of the current status of each remote (what the remote knows about the state).

Remotes
=======

A "remote" is a `historian` instance.

API
===

States are reported and stored as MessagePack binaries (to save space, rather than use JSON).
