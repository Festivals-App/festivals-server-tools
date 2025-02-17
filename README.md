<p align="center">
   <a href="https://github.com/festivals-app/festivals-server-tools/commits/" title="Last Commit"><img src="https://img.shields.io/github/last-commit/festivals-app/festivals-server-tools?style=flat"></a>
   <a href="https://github.com/festivals-app/festivals-server-tools/issues" title="Open Issues"><img src="https://img.shields.io/github/issues/festivals-app/festivals-server-tools?style=flat"></a>
   <a href="./LICENSE" title="License"><img src="https://img.shields.io/github/license/festivals-app/festivals-server-tools.svg"></a>
</p>

<h1 align="center">
  <br/><br/>
    FestivalsApp Server Tools
  <br/><br/>
</h1>

The festivals server tools repository contains server functions shared between server components of the FestivalsApp.

![Figure 1: Architecture Overview Highlighted](https://github.com/Festivals-App/festivals-documentation/blob/main/images/architecture/export/architecture_overview.svg "Figure 1: Architecture Overview")

<hr/>
<p align="center">
  <a href="#development">Development</a> •
  <a href="#deployment">Deployment</a> •
  <a href="#engage">Engage</a>
</p>
<hr/>

## Development

This repository serves as a toolkit for shared functionality among the server components of the `FestivalsApp`. It consolidates common tasks like sending a heartbeat to the [festivals-gateway](https://github.com/Festivals-App/festivals-gateway), updating server binaries from GitHub, responding to client requests, and more. By centralizing these functions, the repository reduces code duplication, simplifies maintenance, and makes refactoring easier.

The development of festivals-server-tools is closely tied to the specific functionality it supports. For example, the heartbeat sender implementation relies heavily on the [festivals-gateway](https://github.com/Festivals-App/festivals-gateway) implementation, whereas the response tools are designed to operate independently of other project components. Each task is designed to minimize dependencies, ensuring the tools remain lightweight and efficient. **However, direct imports from any other `FestivalsApp` components are strictly prohibited.**

### Requirements

- [Golang](https://go.dev/) Version 1.23.5+
- [Visual Studio Code](https://code.visualstudio.com/download) 1.96.0+
  - Plugin recommendations are managed via [workspace recommendations](https://code.visualstudio.com/docs/editor/extension-marketplace#_recommended-extensions).
- [Bash script](https://en.wikipedia.org/wiki/Bash_(Unix_shell)) friendly environment

## Deployment

Add the festivals-server-tools to your go project by running this command

`go get github.com/festivals-app/festivals-server-tools`

## Engage

I welcome every contribution, whether it is a pull request or a fixed typo. The best place to discuss questions and suggestions regarding the festivals-server-tools is the [issues](https://github.com/festivals-app/festivals-server-tools/issues/) section. More general information and a good starting point if you want to get involved is the [festival-documentation](https://github.com/Festivals-App/festivals-documentation) repository.

The following channels are available for discussions, feedback, and support requests:

| Type                     | Channel                                                |
| ------------------------ | ------------------------------------------------------ |
| **General Discussion**   | <a href="https://github.com/festivals-app/festivals-documentation/issues/new/choose" title="General Discussion"><img src="https://img.shields.io/github/issues/festivals-app/festivals-documentation/question.svg?style=flat-square"></a> </a>   |
| **Other Requests**    | <a href="mailto:simon@festivalsapp.org" title="Email me"><img src="https://img.shields.io/badge/email-Simon-green?logo=mail.ru&style=flat-square&logoColor=white"></a>   |

## Licensing

Copyright (c) 2023-2025 Simon Gaus. Licensed under the [**GNU Lesser General Public License v3.0**](./LICENSE)