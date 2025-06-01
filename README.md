# taskcli

`taskcli` is a command-line tool designed to help you manage your tasks efficiently and effectively. With a simple and intuitive interface, `taskcli` allows you to create, view, update, and delete tasks directly from your terminal.

## Features

- Add new tasks with ease.
- View a list of all your tasks.
- Mark tasks as completed.
- Delete tasks you no longer need.
- Lightweight and fast.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/deevanshu-k/taskcli
   ```
2. Navigate to the project directory:
   ```bash
   cd taskcli
   ```
3. Build the package:
   ```bash
   make package
   ```
4. Install the deb package:
   ```bash
   sudo dpkg -i taskcli_<VERSION>_amd64.deb
   ```
5. Update user in `/etc/systemd/system/taskcli.service`: root -> current_user
6. Reload systemd daemon:
   ```bash
   sudo systemctl daemon-reload
   ```
7. Enable and start the service:
   ```bash
   sudo systemctl enable taskcli
   sudo systemctl start taskcli
   ```
8. Verify the service status:
   ```bash
   sudo systemctl status taskcli
   ```

## Info

- Systemd starts the taskcli server in background
- `taskcli` can then send all its commands to the server
- Both use the same config file to get the host and server

## TODO

- Dynamic version passing to go build flags
  `go build -ldflags "-X main.version=$(VERSION)" -o $(APP_NAME)`
- Add git actions to auto publish releases to GitHub
- Add flag for start command --port and --host
