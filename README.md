### Notes

* To create entities of the system namely the hub, client, command and message
* Identify the falling/offline tcp connections and mark them inactive
* Same user_id may get assigned to multiple connections, lets avoid that by some mutex
* Have a loop mechanism to hold the users tcp connection to let him enter the 3 commands I am supposed to design
* Dead clients appearing in the active list - gotta do something about it
* A bug possiblity - what if the relay message is sent my myself or sent to a reciever multiple times
* Testing in scope - unittest, component and performance benchmarks - may have to restructure classes to make it testable
* Can we do docker the whole project ?
* The relative paths are ugly and bad practise, plan to remove it


### Structure

#### Client
user_id		: This integer id will uniquely identify the user in the network
handler		: net.Conn object which will
incoming	: channel of RelayMessage which is supposed to go to the user
timestamp	: timestamp of creation of the user
active      : dunno now
history		: maybe have this as all the commands executed by that particular user
SendMessage() -> Input: RelayMessage Output: Should add the message to the channel which the user may be listening to

#### Network 
assignNodeAddress() -> Output: user_id (integer)
getAllCurrentActive() -> Input: user_id/connection (string/list)
addClientToNetwork() -> Input: Conn.Net. Output: Success, create a new client object and add it to the network 

#### RelayMessage
message: Stores the message(json/string/bytes/anything)
from: Identity of the user who sent the message
receiptClients : Contains the list of user_id which will receive the message

```bash
.
├── Config
│   └── config.go
│
├── Entities
│   ├── Client.go
│   ├── Network.go
│   └── RelayMessage.go
│
├── README.md
│
├── RelayServer
│   └── RelayServer.go
│
├── Utils
│   └── Utils.go
│
└── main.go
```

### Hub setup
```bash
$go run main.go
```
or 
```bash
$go build
$./relay_solution
```
### Client connection
```bash
$nc localhost 6666
```
### Client usage
```bash
$IDENTIFY
$LIST
$RELAY #Message i want to send, its a lame way but lets do this #1,3,4,10
```
### Unittest 
```bash
$ go test $(go list ./...| grep -v test) -coverprofile cover.out; go tool cover -func cover.out
```
or go to the respective folder
```bash
$ go test -coverprofile cover.out; go tool cover -func cover.out
```
