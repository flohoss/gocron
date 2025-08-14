## ✨ New Features

- Improve the terminal command whitelist to be more efficient

  ```yaml
  # Terminal configuration
  terminal:
    # When set to true, any command will be allowed,
    # and the allowed_commands list will be ignored.
    allow_all_commands: false
    # Be careful with this, it will allow a command to be executed.
    # Do not allow any commands that can compromise the system!
    allowed_commands:
      cat:
        args:
          - /config/config.yaml
      docker:
        args:
          - ps
          - version
      restic:
        args:
          - version
      export:
        # Do not check for arguments, allow all arguments
        allow_all_args: true
  ```
