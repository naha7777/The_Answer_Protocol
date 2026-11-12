_This project has been created as part of the 42 curriculum by anacharp, frrenaux._

# The_Answer_Protocol

## Description

## Instructions

## Architecture
section explaining your server design choices (dispatcher/router vs inline handling, concurrency model, etc.)

## Protocole Implementation

### Commands (current state)

| Command | Success response | Error responses | Events pushed | Status |
|---|---|---|---|---|
| *(on connect)* | `OK hello proto=1` | | | done |
| any command before `CONNECT` | | `ERR 100 YOU_MUST_BE_CONNECTED` | | done |
| `CONNECT <name>` | `OK connected` | `ERR 101 YOU_HAVE_A_NAME`, `ERR 102 THIS_IS_NOT_A_NAME`, `ERR 201 NAME_IN_USE` | none | done |
| `QUIT` | `OK bye` | | none | partial |
| `LOOK` | `OK {"room": {...}}` | | | partial |
| `MOVE <direction>` | `OK room=<room-id>` | `ERR 301 NO_EXIT` | none | partial |
| `TAKE <item id or name>` | `OK taken=<item-id>` | `ERR 115 MISSING_ITEM`, `ERR 116 ITEM_NOT_OBTAINABLE`, `ERR 404 ITEM_NOT_FOUND` | none | done |
| `DROP <item id or name>` | `OK dropped=<item-id>` | `ERR 115 MISSING_ITEM`, `ERR 404 ITEM_NOT_FOUND` | none | done |
| `INVENTORY` | `OK [item.a item.b]` | | | partial |
| `TALK <npc id or name>` | `OK <dialogue line>` | `ERR 404 NPC_NOT_FOUND` | | partial |
| `ATTACK <npc id or name>` | `OK {'attacker_hp': N, 'target_hp': N, 'damage': N, 'status': S}` | `ERR 404 NPC_NOT_FOUND`, `ERR 405 NPC_NOT_HOSTILE` | none | partial |
| `STATUS` | `OK {'hp': N, 'max_hp': N, 'status': S}` | | | done |
| `WHO` | `OK players=<n>` | | | partial |
| `CHAT GLOBAL <msg>` | `OK` (sender) | `ERR 110 NO_CHAT_SPECIFIED` | `<name>: <msg>` to others | partial |
| `CHAT ROOM <msg>` | `OK` (sender) | `ERR 110 NO_CHAT_SPECIFIED` | `<name>: <msg>` to others in room | partial |
| `CHAT GROUP <msg>` | `OK` (sender) | `ERR 401 NOT_IN_GROUP`, `ERR 110 NO_CHAT_SPECIFIED` | `<name>: <msg>` to group members | partial |
| `GROUP CREATE` | `OK group=<id>` | `ERR 402 ALREADY_IN_GROUP` | none | done |
| `GROUP INVITE <user>` | `OK` | `ERR 107 CANT_FIND_USER`, `ERR 402 ALREADY_IN_GROUP` | `EVT GROUP INVITE` to the invited user | partial |
| `GROUP JOIN <leader>` | `OK group=<id>` | `ERR 104 MISSING_ARGUMENT=LEADER_NAME`, `ERR 402 ALREADY_IN_GROUP`, `ERR 107 CANT_FIND_LEADER`, `ERR 108 THIS_USER_IS_NOT_A_LEADER`, `ERR 109 REQUIRE_INVITATION` | `EVT GROUP JOIN` to members | partial |
| `GROUP LEAVE` | `OK` | `ERR 401 NOT_IN_GROUP` | `EVT GROUP LEAVE` to remaining members | partial |
| `QUEST`, `QUESTS`, `COMPLETE_QUEST` | | | | not implemented |
| `DEFEND`, `FLEE`, `USE_ITEM` | | | | not implemented |
| unknown command | `Unknown command: <cmd>` | | | to fix |

### Error codes

| Code | Name | Source |
|---|---|---|
| 201 | NAME_IN_USE | RFC |
| 301 | NO_EXIT | RFC |
| 401 | NOT_IN_GROUP | RFC |
| 402 | ALREADY_IN_GROUP | RFC |
| 404 | ITEM_NOT_FOUND / NPC_NOT_FOUND | RFC |
| 405 | NPC_NOT_HOSTILE | RFC |
| 100 | YOU_MUST_BE_CONNECTED | **custom (not in RFC)** |
| 107 | CANT_FIND_USER | **custom (not in RFC)** |
| 108 | THIS_USER_IS_NOT_A_LEADER | **custom (not in RFC)** |
| 104 | MISSING_ARGUMENT=LEADER_NAME | **custom (not in RFC)** |
| 102 | THIS_IS_NOT_A_NAME | **custom (not in RFC)** |
| 101 | YOU_HAVE_A_NAME | **custom (not in RFC)** |
| 109 | REQUIRE_INVITATION | **custom (not in RFC)** |
| 110 | NO_CHAT_SPECIFIED | **custom (not in RFC)** |
| 115 | MISSING_ITEM | **custom (not in RFC)** |
| 116 | ITEM_NOT_OBTAINABLE | **custom (not in RFC)** |

### Events

| Event | Format sent | RFC format |
|---|---|---|
| Group invitation | `EVT GROUP INVITE` | `EVT GROUP INVITE <leader>` |
| Group join | `EVT GROUP JOIN` | `EVT GROUP JOIN <username>` |
| Group leave | `EVT GROUP LEAVE` | `EVT GROUP LEAVE <username>` |
| Chat | `<name>: <msg>` | `EVT <SCOPE> CHAT <username> <message>` |
| Room presence | not sent | `EVT ROOM PRESENCE ENTER/LEAVE <username>` |
| Stats | not sent | `EVT STATS players=<n>` |

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
