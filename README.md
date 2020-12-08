# fuzz-linearizabilty
Fuzzing input for linearizability in distributed systems.

## Testing RQLite 

1. Download [rqlite](https://github.com/rqlite/rqlite)

On Linux:
```
curl -L https://github.com/rqlite/rqlite/releases/download/v5.6.0/rqlite-v5.6.0-linux-amd64.tar.gz -o rqlite-v5.6.0-linux-amd64.tar.gz
tar xvfz rqlite-v5.6.0-linux-amd64.tar.gz
```
1. To start a single node 
```
cd rqlite-v5.6.0-linux-amd64
./rqlited -http-addr localhost:4001 -raft-addr localhost:4002 ~/node.1
```
1.  Multiple nodes
```
./rqlited -http-addr localhost:4003 -raft-addr localhost:4004 -join http://localhost:4001 ~/node.2
./rqlited -http-addr localhost:4005 -raft-addr localhost:4006 -join http://localhost:4001 ~/node.3
./rqlited -http-addr localhost:4007 -raft-addr localhost:4008 -join http://localhost:4001 ~/node.4
./rqlited -http-addr localhost:4009 -raft-addr localhost:4010 -join http://localhost:4001 ~/node.5
```
1. Now you can start testing for linearizabilty, See `main.go` for an example. 