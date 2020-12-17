#!/bin/sh
xterm -title "Node 1" -hold -e "./rqlited ~/node.1"  &
xterm -title "Node 2" -hold -e "./rqlited -http-addr localhost:4003 -raft-addr localhost:4004 -join http://localhost:4001 ~/node.2"  &
xterm -title "Node 3" -hold -e "./rqlited -http-addr localhost:4005 -raft-addr localhost:4006 -join http://localhost:4001 ~/node.3"  &
xterm -title "Node 4" -hold -e "./rqlited -http-addr localhost:4007 -raft-addr localhost:4008 -join http://localhost:4001 ~/node.4"  &
xterm -title "Node 5" -hold -e "./rqlited -http-addr localhost:4009 -raft-addr localhost:4010 -join http://localhost:4001 ~/node.5"