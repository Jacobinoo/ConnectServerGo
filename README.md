![Connect Server Logo](jacobinoo.github.io/assets/ConnectServerLogo.png)
# Connect Server for Go (Work In Progress)

Connect Server (Minerva) is a REST backend that implements and exposes set of APIs needed to deliver functional and performant server infrastructure for Connect Messenger clients to communicate with. **Based on routing. Built using Go.**

## Building
![Connect Server Go Build](https://github.com/Jacobinoo/ConnectServerGo/actions/workflows/go.yml/badge.svg)

Go has to be installed in order to build the module. Navigate to the root of the project folder. Use `go build -o PATH/TO/BINARY/OUTPUT/FOLDER` command. You can modify the output flag to change the location of the binary.

```bash
go build -o ./bin
```
The command above will build and compile Connect Server to `bin` directory.

## Usage

Execute the `ConnectServer` executable file using the Terminal: supply the `.env` environment file path using the `APP_ENV_PATH` environment variable. Schema of the environment file is available inside `.env.schema` file in the source code.

```bash
APP_ENV_PATH=PATH/TO/ENV/FILE ./ConnectServer
```

## Contributions

At the moment, we do only accept **security contributions**. If you see a vulnerability of some kind, we kindly encourage you to share this information, and report it using Github Issues.

More information on security contributions is available inside [SECURITY.md]().

## License

[Not yet implemented - All rights reserved]()

Any form of redistribution is forbidden.

The source code is only for viewing purposes, excluding security contributors who can modify the code.

The source code may only be modified for security contribution reasons.

Copyright: Jakub Banasiewicz 2024
