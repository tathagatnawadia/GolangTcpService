### Notes

* To create entities of the system namely the hub, client, command and message
* Identify the falling/offline tcp connections and mark them inactive
* Same user_id may get assigned to multiple connections, lets avoid that by some mutex
* Have a loop mechanism to hold the users tcp connection to let him enter the 3 commands I am supposed to design
* Dead clients appearing in the active list - gotta do something about it

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

### Hub setup
```bash
$go run main.go
```
### Client connection
```bash
$nc localhost 6666
```
### Client usage
```bash
$IDENTIFY
$LIST
$RELAY
```
### Unittest 
```bash
$todo
```
