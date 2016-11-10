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

A "remote" is a remote api that is interested in tracking a local state stream. A remote should implement the `ReporterRemoteService` as declared in the `remote` sub-package of this repo.

A remote is configured in reporter with the following parameters:

 - Endpoint: how to get to the remote.

That's it. Internally the reporter will, for each remote, call the GetStreamInterest endpoint to ask which streams the remote is interested in receiving. This will return an object describing which streams at what RateConfigs the remote is interested in receiving, as well as an nonce number with the version of this config.

The reporter will then "backfill" the remote stream:

 - Start at the latest timestamp of the remote
 - Write all events in the stream to the remote stream
 - From then onward, keep the remote's stream in sync with the writer.

Then, in another goroutine:

 - Connect to the remote (use an internal complex connection state / backoff mechanism)
 - Start at the oldest unsynced stream entry, and push it to the remote, and continue as such till the latest. For each successful push, remove the entry from the local stream.
 - As states are added to the push queue, redo the above.
 - Each push call returns the nonce of the configuration for the remote, and if this has changed, the reporter will check the config endpoint again, grab the new config, apply it to the `EntryEmitter` and all future pushes will have the new config applied.

This will allow a remote to achieve the following:

 - Receive a stream of state change entries with the desired rate config.
 - Stay synced with the latest state.
 - Pull data / history all the way since a desired starting point.

The remote will need to know some identity information about the reporter in order to make an educated decision about what streams it's interested in receiving. This identifier, in FuseCloud's case, would be something like `plane_1` - the hostname of the machine. The reporter on default will use the hostname of the machine for this identifier, but otherwise can use an identifier given on the command line.

Remote Authentication
=====================

This will not be initially implemented, but added down the line.

When calling GetRemoteConfig the reporter gives an identifier. This identifier needs to be authenticated to make sure reporters can't be forged.

API
===

The "reporter" api is intended for local components on a system communicating with the central reporter instance.

The "view" api is intended for consumers of state information and state history. It is shared between historian and reporter, and allows:

 - View state at a given time
 - View latest state
 - View history of a state, with any desired rate configuration (server-side batching of events)
 - Tail latest state (receive live updates)
 - Request history between two time points at varying resolutions.

Bounded State History
=====================

View api used by the "windows" algorithm in the remote-state-stream.

Request contains:

 - Requested middle timestamp
 - Mode (currently just snapshot bound)

Response messages:

 - START_BOUND with early bound. If null, terminates call here.
 - END_BOUND with end bound. If null, is a tailing call.
 - INITIAL_SET with each known entry.
 - TAIL with nil to signal end of initial set
 - TAIL with each new entry, if needed until...
 - END_BOUND with ending snapshot after tailing.

Whereas the state history call requests a history with a specific rate config, and re-builds each entry, bounded state history call sticks strictly to the entry stream from the backend. This way the data in the browser will exactly match the data in the backend, but can be fetched in partials.
