cmd.exe /c start bash.exe -c "./rqlited ~/node.1" & 
cmd.exe /c start bash.exe -c "./rqlited -http-addr localhost:4003 -raft-addr localhost:4004 -join http://localhost:4001 ~/node.2" &
cmd.exe /c start bash.exe -c "./rqlited -http-addr localhost:4005 -raft-addr localhost:4006 -join http://localhost:4001 ~/node.3" &
cmd.exe /c start bash.exe -c "./rqlited -http-addr localhost:4007 -raft-addr localhost:4008 -join http://localhost:4001 ~/node.4" & 
cmd.exe /c start bash.exe -c "./rqlited -http-addr localhost:4009 -raft-addr localhost:4010 -join http://localhost:4001 ~/node.5"
