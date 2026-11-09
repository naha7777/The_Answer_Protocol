_This project has been created as part of the 42 curriculum by anacharp, frrenaux._

# The_Answer_Protocol

## Description

## Instructions

## Architecture
section explaining your server design choices (dispatcher/router vs inline handling, concurrency model, etc.)

## Protocole Implementation
section documenting any deviations from RFC 42TAP and justifying your choices.

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

## Add Error Codes

During developpement, we needed to add more error codes :

| Code | Message | Description |
| ---- | ------- | ----------- |
| 100 | YOU_MUST_BE_CONNECTED | Request a connection |
| 107 | CANT_FIND_USER | Invalid username for a command argument |
| 108 | THIS_USER_IS_NOT_A_LEADER | Invalid leader for group join command |
| 104 | MISSING_ARGUMENT=LEADER_NAME | Missing argument for group join command |
| 102 | THIS_IS_NOT_A_NAME | Invalid username during connect command |
| 101 | YOU_HAVE_A_NAME | Try to connect twice without deconnection |
| 109 | REQUIRE_INVITATION | If a client try to join a group without invitation |

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

#### GUI
- [Package Fyne](https://docs.fyne.io/)

### AI Usage
