### Cli Overview
1) Name: `taskcli`
2) Server Command: `taskcli start`
    - `--host 127.0.0.1 --port 3500` (Optional)
3) Client Commands:
    - Add task `taskcli add 'Work on blog post'`
    - Update task `taskcli update <taskid> -d 'Update a feature in x project' -s 'P/I/C'`
        - `-d` | `-decription` to change the task description.
        - `-s` | `-status` to change the task status. values: `P` for pending, `I` for in-progress, `C` for completed
    - Delete task `taskcli delete <taskid>`
    - Delete all tasks `taskcli delete -all`