_This project has been created as part of the 42 curriculum by anacharp, frrenaux and tcolson._

# The_Answer_Protocol

## Description

## Instructions

Server:
```bash
go run ./cmd/server
```

GUI:
```bash
go run ./cmd/gui
```

CLI:
```bash
go run ./cmd/cli
```

## Architecture
section explaining your server design choices (dispatcher/router vs inline handling, concurrency model, etc.)

## Protocole Implementation

### Commands (current state)

| Command | Success response | Error responses | Events pushed |
|---|---|---|---|
| *(on connect)* | `OK hello proto=1` | none | none |
| any command before `CONNECT` | | `ERR 100 YOU_MUST_BE_CONNECTED` | |
| `CONNECT <name>` | `OK connected` | `ERR 10 YOU_HAVE_A_NAME`, `ERR 09 THIS_IS_NOT_A_NAME`, `ERR 201 NAME_IN_USE` | none |
| `QUIT` | `OK bye` | | none |
| `LOOK` | `OK look={"room": {"Id":"loc.square","Name":"Village Square","Description":"A bustling square.","Exits":{"east":"loc.market","north":"loc.forest"},"Players":["arbre"],"Items":["item.herbs","item.herbs"],"NPCs":["npc.village_guard"]}}` | | |
| `MOVE <direction>` | `OK room=<room-id>` | `ERR 301 NO_EXIT` | none |
| `TAKE <item id or name>` | `OK taken=<item-id>` | `ERR 07 ITEM_NOT_OBTAINABLE`, `ERR 404 ITEM_NOT_FOUND`, `ERR 17 FULL_INVENTORY` | `"EVT PLAYER TOOK <item_id>` |
| `DROP <item id or name>` | `OK dropped=<item-id>` | `ERR 404 ITEM_NOT_FOUND` | `"EVT PLAYER DROPPED <item_id>` |
| `INVENTORY` | `OK inventory=[item.a item.b]` | | |
| `TALK <npc id or name>` | `OK talk=<npc_id>: <dialogue line>` | `ERR 404 NPC_NOT_FOUND` | |
| `STATUS` | `OK status={"hp": <client_hp>, "max_hp": <client_maxp_hp>, "status": "<client_status>"}` | | | none |
| `ATTACK <npc id or name>` | `OK attack={"player_hp": 36, "player_max_hp": 100, "NPC_name": Robber, "NPC_hp": 1, "NPC_max_hp": 100, "damage": 8, "damage_type": normal, "status": combat}`, `OK attack=<ncp_id>=dead, "status"="<client_status>"` | `ERR 404 NPC_NOT_FOUND`, `ERR 405 NPC_NOT_HOSTILE`, `ERR 14 NPC_IN_FIGHT` | `EVT NPC <npc_id> DEAD`, `EVT NPC FLEE, "status"="<player_status>"`, `EVT NPC ATTACK`, `EVT NPC REGENERATION : +5 hp`, `EVT POISON : -5 npc hp during <nb_turns> turns`, `EVT STRENGH : +5 player strengh during <nb_turns> turns`, `EVT HEALTH : +5 player hp during <nb_turns> turns`, `EVT POISON : -5 player hp during <nb_turns> turns` |
| `FLEE` | `OK flee=succed, "status"="<client_status>"`, `OK flee=fail, "status"="combat"` | `ERR 13 NOT_IN_FIGHT` | `EVT NPC ATTACK`, `EVT NPC REGENERATION : +5 hp`, `EVT POISON : -5 npc hp during <nb_turns> turns`, `EVT STRENGH : +5 player strengh during <nb_turns> turns`, `EVT HEALTH : +5 player hp during <nb_turns> turns`, `EVT POISON : -5 player hp during <nb_turns> turns` |
| `DEFEND` | `OK defend, "status"="combat"`, `OK attack={"player_hp": 36, "player_max_hp": 100, "NPC_name": Robber, "NPC_hp": 1, "NPC_max_hp": 100, "damage": 8, "damage_type": normal, "status": combat}` | `ERR 13 NOT_IN_FIGHT` | `EVT NPC FLEE, "status"="<client_status>"`, `EVT NPC ATTACK`, `EVT NPC REGENERATION : +5 hp`, `EVT POISON : -5 npc hp during <nb_turns> turns`, `EVT STRENGH : +5 player strengh during <nb_turns> turns`, `EVT HEALTH : +5 player hp during <nb_turns> turns`, `EVT POISON : -5 player hp during <nb_turns> turns` |
| `USE_ITEM <item>` | `OK use=<item_id>` | `ERR 12 IN_FIGHT`, `ERR 404 ITEM_NOT_FOUND`, `ERR 13 NOT_IN_FIGHT` | `EVT BEER : -5 player hp during <nb_turns> turns`, `EVT SWORD : +5 player strengh during <nb_turns> turns`, `EVT POISON POTION : -5 npc hp during <nb_turns> turns`, `EVT STRENGH POTION : +5 player strengh during <nb_turns> turns`, `EVT HEALTH POTION : +5 player hp during <nb_turns> turns` |
| `WHO` | `OK players=<n>` | | |
| `CHAT GLOBAL <msg>` | `OK` (sender) | `ERR 06 NO_CHAT_SPECIFIED`, `ERR 08 MISSING_MESSAGE` | `EVT GLOBAL CHAT <username> <message>` to others |
| `CHAT ROOM <msg>` | `OK` (sender) | `ERR 06 NO_CHAT_SPECIFIED`, `ERR 08 MISSING_MESSAGE` | `EVT ROOM CHAT <username> <message>` to others in room |
| `CHAT GROUP <msg>` | `OK` (sender) | `ERR 401 NOT_IN_GROUP`, `ERR 06 NO_CHAT_SPECIFIED`, `ERR 08 MISSING_MESSAGE` | `EVT GROUP CHAT <username> <message>` to group members |
| `GROUP CREATE <id>` | `OK group_created=<id>` | `ERR 402 ALREADY_IN_GROUP` | none | done |
| `GROUP INVITE <user>` | `OK group_invited=<username>` | `ERR 03 CANT_FIND_USER`, `ERR 402 ALREADY_IN_GROUP` | `EVT GROUP INVITE <leadername>` to the invited user |
| `GROUP JOIN <leader>` | `OK group_joined=<id>` | `ERR 402 ALREADY_IN_GROUP`, `ERR 03 CANT_FIND_LEADER`, `ERR 04 THIS_USER_IS_NOT_A_LEADER`, `ERR 05 REQUIRE_INVITATION` | `EVT GROUP JOIN <username>` to members |
| `GROUP LEAVE` | `OK group_leaved=<group_id>` | `ERR 401 NOT_IN_GROUP` | `EVT GROUP LEAVE <username>` to remaining members |
| `INFO_GROUP` | `OK group_infos={"group": {"Id":"<group_ID>", "Leader":"group_Leader", "Members":"[<username> <username> <username>]", "Invited":"[<username>]"}` | `ERR 401 NOT_IN_GROUP` | none |
| `QUEST <npc name or npc id>` | `OK quest={"quest_id": "princess_necklace", "description": "Bring me a beautiful necklace please.", "reward": "Bravery", "status": "available"}` | `ERR 406 NO_QUEST_AVAILABLE`, `ERR 404 NPC_NOT_FOUND` | `EVT QUEST took=<quest_id>` |
| `QUESTS` | `OK quests=[]`, `OK quests=[{"quest_id": "princess_necklace", "status": "active", "progress": "0/1"}, {"quest_id: "example", "status": "completed"}]` | none | none |
| `COMPLETE_QUEST <quest_id>` | `OK complete_quest=quest_id: completed`, `OK taken=item.potion` | `ERR 406 NO_QUEST_AVAILABLE`, `ERR 11 INCOMPLETE_QUEST_REQUIREMENTS`, `ERR 404 NPC_NOT_FOUND`, `ERR 17 FULL_INVENTORY` | `EVT QUEST completed=<quest_id>` |
| `ABANDON_QUEST <quest_id>` | `OK abandoned=<quest_id>` | `406 NO_QUEST_AVAILABLE` | `EVT QUEST abandoned=<quest_id>` |
| `FIND_QUESTS` | `OK find_quests={"quests":[{"npc":"Bartender","quests":["kill_pirate","kill_captain"]}]}` | `ERR 15 NO_QUEST_AVAILABLE` |
| UNKNOWN COMMAND | | `ERR 01 UNKNOWN_COMMAND`|
| error CLI or GUI connection | | `900 CONNECTION_FAILED` |
| error CLI or GUI send | | `900 SEND_FAILED` |
| when someone CONNECT himself or QUIT | everyone receive `EVT STATS players=nb_players_in_game` and `EVT STATS players_in_room=nb_players_in_his_room` | |


### Error codes

| Code | Name | Source |
|------|------|--------|
| 201 | NAME_IN_USE | RFC |
| 301 | NO_EXIT | RFC |
| 401 | NOT_IN_GROUP | RFC |
| 402 | ALREADY_IN_GROUP | RFC |
| 404 | ITEM_NOT_FOUND / NPC_NOT_FOUND / ITEM_NOT_IN_INVENTORY | RFC |
| 405 | NPC_NOT_HOSTILE | RFC |
| 900 | CONNECTION_FAILED | RFC |
| 901 | SEND_FAILED | RFC |
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
| 11 | INCOMPLETE_QUEST_REQUIREMENTS | **custom (not in RFC)** |
| 12 | IN_FIGHT | **custom (not in RFC)** |
| 13 | NOT_IN_FIGHT | **custom (not in RFC)** |
| 14 | NPC_IN_FIGHT | **custom (not in RFC)** |
| 15 | NO_QUEST_AVAILABLE | **custom (not in RFC)** |
| 16 | LINE_TOO_LONG | **custom (not in RFC)** |
| 17 | FULL_INVENTORY | **custom (not in RFC)** |
| 500 | INTERNAL | **custom (not in RFC)** |

### Events

| Event | Format | Source |
|-------|--------|--------|
| Group invitation | `EVT GROUP INVITE <leader>` | RFC |
| Group join | `EVT GROUP JOIN <username>` | RFC |
| Group leave | `EVT GROUP LEAVE <username>` | RFC |
| Chat | `EVT <SCOPE> CHAT <username> <message>` | RFC |
| Room presence | `EVT ROOM PRESENCE ENTER/LEAVE <username>` | RFC |
| Stats | `EVT STATS players= <nb_player>` | RFC |
| NPC Death | `EVT NPC DEAD <npc_id>` | **custom (not in RFC)** |
| NPC Flee | `EVT NPC FLEE, "status"="<npc_id>"` | **custom (not in RFC)** |
| NPC Attack | `EVT NPC ATTACK` | **custom (not in RFC)** |
| NPC Regeneration | `EVT NPC REGENERATION : +5 hp` | **custom (not in RFC)** |
| I am dead | `EVT DEATH` | **custom (not in RFC)** |
| A player is dead | `EVT dead=<player_name>` | **custom (not in RFC)** |
| Use_item beer | `EVT BEER : -5 player hp during <nb_turns> turns` | **custom (not in RFC)** |
| Use_item sword | `EVT SWORD : +5 player strengh during <nb_turns> turns` | **custom (not in RFC)** |
| Use_item poison potion | `EVT POISON POTION : -5 npc hp during <nb_turns> turns` | **custom (not in RFC)** |
| Use_item strengh potion | `EVT STRENGH POTION : +5 player strengh during <nb_turns> turns` | **custom (not in RFC)** |
| Use_item health potion | `EVT HEALTH POTION : +5 player hp during <nb_turns> turns` | **custom (not in RFC)** |
| Poison npc turns | `EVT POISON : -5 npc hp during <nb_turns> turns` | **custom (not in RFC)** |
| Strengh turns | `EVT STRENGH : +5 player strengh during <nb_turns> turns` | **custom (not in RFC)** |
| Health turns | `EVT HEALTH : +5 player hp during <nb_turns> turns` | **custom (not in RFC)** |
| Poison player turns | `EVT POISON : -5 player hp during <nb_turns> turns` | **custom (not in RFC)** |
| A player abandon a quest | `EVT QUEST abandoned=<quest_id>` | **custom (not in RFC)** |
| A player complete a quest | `EVT QUEST completed=<quest_id>` | **custom (not in RFC)** |
| A player take a quest | `EVT QUEST took=<quest_id>` | **custom (not in RFC)** |

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
