_This project has been created as part of the 42 curriculum by anacharp, frrenaux._

# The_Answer_Protocol

## Description

## Instructions

## Architecture
section explaining your server design choices (dispatcher/router vs inline handling, concurrency model, etc.)

## Protocole Implementation

### Commands (current state)

| Command | Success response | Error responses | Events pushed |
|---|---|---|---|
| *(on connect)* | `OK hello proto=1` | | |
| any command before `CONNECT` | | `ERR 100 YOU_MUST_BE_CONNECTED` | |
| `CONNECT <name>` | `OK connected` | `ERR 10 YOU_HAVE_A_NAME`, `ERR 09 THIS_IS_NOT_A_NAME`, `ERR 201 NAME_IN_USE` | none |
| `QUIT` | `OK bye` | | none |
| `LOOK` | `OK {"room": {"Id":"loc.square","Name":"Village Square","Description":"A bustling square.","Exits":{"east":"loc.market","north":"loc.forest"},"Players":["arbre"],"Items":["item.herbs","item.herbs"],"NPCs":["npc.village_guard"]}}` | | |
| `MOVE <direction>` | `OK room=<room-id>` | `ERR 301 NO_EXIT` | none |
| `TAKE <item id or name>` | `OK taken=<item-id>` | `ERR 07 ITEM_NOT_OBTAINABLE`, `ERR 404 ITEM_NOT_FOUND` | none | done |
| `DROP <item id or name>` | `OK dropped=<item-id>` | `ERR 404 ITEM_NOT_FOUND` | none |
| `INVENTORY` | `OK [item.a item.b]` | | |
| `TALK <npc id or name>` | `OK <dialogue line>` | `ERR 404 NPC_NOT_FOUND` | |
| `ATTACK <npc id or name>` | `OK {'attacker_hp': N, 'target_hp': N, 'damage': N, 'status': S}` | `ERR 404 NPC_NOT_FOUND`, `ERR 405 NPC_NOT_HOSTILE` | none |
| `STATUS` | `OK {'hp': N, 'max_hp': N, 'status': S}` | | | done |
| `WHO` | `OK players=<n>` | | |
| `CHAT GLOBAL <msg>` | `OK` (sender) | `ERR 06 NO_CHAT_SPECIFIED`, `ERR 08 MISSING_MESSAGE` | `EVT GLOBAL CHAT <username> <message>` to others |
| `CHAT ROOM <msg>` | `OK` (sender) | `ERR 06 NO_CHAT_SPECIFIED`, `ERR 08 MISSING_MESSAGE` | `EVT ROOM CHAT <username> <message>` to others in room |
| `CHAT GROUP <msg>` | `OK` (sender) | `ERR 401 NOT_IN_GROUP`, `ERR 06 NO_CHAT_SPECIFIED`, `ERR 08 MISSING_MESSAGE` | `EVT GROUP CHAT <username> <message>` to group members |
| `GROUP CREATE` | `OK group=<id>` | `ERR 402 ALREADY_IN_GROUP` | none | done |
| `GROUP INVITE <user>` | `OK` | `ERR 03 CANT_FIND_USER`, `ERR 402 ALREADY_IN_GROUP` | `EVT GROUP INVITE` to the invited user |
| `GROUP JOIN <leader>` | `OK group=<id>` | `ERR 402 ALREADY_IN_GROUP`, `ERR 03 CANT_FIND_LEADER`, `ERR 04 THIS_USER_IS_NOT_A_LEADER`, `ERR 05 REQUIRE_INVITATION` | `EVT GROUP JOIN` to members |
| `GROUP LEAVE` | `OK` | `ERR 401 NOT_IN_GROUP` | `EVT GROUP LEAVE` to remaining members |
| `QUEST`, `QUESTS` | | | | not implemented |
| UNKNOWN COMMAND | | `ERR 01 UNKNOWN_COMMAND`|

### Error codes

| Code | Name | Source |
|------|------|--------|
| 201 | NAME_IN_USE | RFC |
| 301 | NO_EXIT | RFC |
| 401 | NOT_IN_GROUP | RFC |
| 402 | ALREADY_IN_GROUP | RFC |
| 404 | ITEM_NOT_FOUND / NPC_NOT_FOUND | RFC |
| 405 | NPC_NOT_HOSTILE | RFC |
| 01  | UNKNOWN_COMMAND | **custom (not in RFC)** |
| 02 | YOU_MUST_BE_CONNECTED | **custom (not in RFC)** |
| 03 | CANT_FIND_USER | **custom (not in RFC)** |
| 04 | THIS_USER_IS_NOT_A_LEADER | **custom (not in RFC)** |
| 05 | REQUIRE_INVITATION | **custom (not in RFC)** |
| 06 | NO_CHAT_SPECIFIED | **custom (not in RFC)** |
| 07 | ITEM_NOT_OBTAINABLE | **custom (not in RFC)** |
| 08 | MISSING_MESSAGE | **custom (not in RFC)** |
| 09 | THIS_IS_NOT_A_NAME | **custom (not in RFC)** |
| 10 | YOU_HAVE_A_NAME | **custom (not in RFC)** |

### Events

| Event | Format |
|-------|--------|
| Group invitation | `EVT GROUP INVITE <leader>` |
| Group join | `EVT GROUP JOIN <username>` |
| Group leave | `EVT GROUP LEAVE <username>` |
| Chat | `EVT <SCOPE> CHAT <username> <message>` |
| Room presence | `EVT ROOM PRESENCE ENTER/LEAVE <username>` |

## Combat System
section describing your turn-based combat mechanics, damage formulas, initiative order, and additional combat commands (DEFEND, FLEE, etc.).

## Quest System
section explaining your quest progression mechanics, completion validation, and reward systems.

## World Design
section describing your world layout, room connections, NPC roles, and item distribution

## Server Logging
section documenting your logging implementation, including log format, event types, output destinations, and how to monitor server behavior and detect abuse patterns.

## Group Contributions
section clearly indicating each team member’s responsibilities and contributions to different components (server, CLI client, GUI client, world design, etc.).

## Building and Running
section with detailed instructions for your chosen building tool and how to run each component (server, CLI client, GUI client).

## Testing
section explaining how to test the multiplayer functionality, combat system, and quest mechanics.

## Resources

### Documentation
#### GoLang
- [Cours: Apprendre GoLang](https://devopssec.fr/category/apprendre-golang)
- [Go by Example](https://gobyexample.com/)
- [A Tour of Go](https://fr-go-dev.shuijingwanwq.com/tour/flowcontrol/8)
- [Les fondamentaux de Go - Organisation d'un projet](https://go-start.vercel.app/fr/chapter-1-introduction/1-4-project-layout/)
- [Les fondamentaux de Go - Commandes essentielles](https://go-start.vercel.app/fr/chapter-1-introduction/1-5-essential-commands/)

#### Server
- [Qu’est-ce que le protocole TCP/IP de contrôle de la transmission ?](https://www.fortinet.com/fr/resources/cyberglossary/tcp-ip)
- [Build a Blazing-Fast TCP Server in Go: A Practical Guide](https://dev.to/jones_charles_ad50858dbc0/build-a-blazing-fast-tcp-server-in-go-a-practical-guide-29d)
- [Go by Example: TCP Server](https://gobyexample.com/tcp-server)
- [Les fondamentaux de Go - Serveur TCP](https://go-start.vercel.app/fr/chapter-10-cli-mastery-build-a-publishable-task-manager/10-17-tcp-server/)
- [Package net](https://pkg.go.dev/net)
- [JSON and Go](https://go.dev/blog/json)
- [Switch](https://go.dev/tour/flowcontrol/9)
- [Goroutines](https://go.dev/tour/concurrency/1)
- [Channels](https://go.dev/tour/concurrency/2)
- [Buffered Channels](https://go.dev/tour/concurrency/3)

#### Client
- [What is the net.DialTimeout function in Golang?](https://www.educative.io/answers/what-is-the-netdialtimeout-function-in-golang)

#### GUI
- [Package Fyne](https://docs.fyne.io/)

### AI Usage
