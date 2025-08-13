## ✨ New Features

- Only reload config if validation runs successfully
- Whitelist for terminal commands:

  ```yaml
  # Terminal configuration
  terminal:
  # When set to true, any command will be allowed,
  # and the allowed_commands list will be ignored.
  allow_all_commands: false
  # Be careful with this, it will allow a command to be executed.
  # Do not allow any commands that can compromise the system!
  allowed_commands:
    - command: cat
      args:
        - /config/config.yaml
    - command: docker
      args:
        - ps
        - version
    - command: restic
      args:
        - version
    - command: rclone
      args:
        - version
  ```
