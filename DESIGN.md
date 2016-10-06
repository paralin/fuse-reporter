Overall Summary
===============

Reporter keeps track of a few things:

 - List of components on the system
 - List of states in these components
 - History of these states.

It does this using the historian state stream package.

Reporter Design
===============

Reporter uses BoltDB as a reliable key/value store.

There is a top level "Component" key which contains a list of known local components. Each local component has its own key, which lists known states. Each state has two keys: a list of snapshot timestamps, and a list of mutation timestamps. Each snapshot / mutation has its own key in the DB.

The "global" bucket contains:

 - `components`: list of all components

Each component has a bucket (named `cmp.{component}`):

 - `component`: component level key, lists states

Each state has a bucket in the component bucket (named `st.{statename}`).

 - `state`: state level key, lists timestamps

Each state has a sub-bucket for snapshots and one for mutations, each with timestamp used as key.

Lazy Loading
============

A individual component state stream is lazy-loaded out of the DB when they are first referenced.
