### Overall Architecture

#### Problem
1) Need cli tool to manage local todos or tasks.
2) Data should be persist locally.
3) On system restart, their should be no inconsistency.

#### Solution
1) Application will use client-server architecture.
2) Single cli application is used to start server and client.
3) eg
    `taskcli start` will start the tcp server
    `taskcli add <task>` send the add task command to tcp server
4) Everytime tcp server starts
    - It first checks for config file.
        - If present, then load the config and start server.
        - Else create config and storage file.
    - Starts listening on a port for commands.
    - It saves data on storage file.
5) On cli commands
    - It also checks for config file.
        - If present, then get host and port.
        - Else throw error that cli demon server is not running.
    - It then send the required command to tcp server with config host and port or with passed host and port.