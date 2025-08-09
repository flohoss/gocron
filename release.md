# 🚀 GoCron v0.7.0

## ‼️ BREAKING CHANGES

- The database and the configuration has been completely rewritten.
- Please delete your current db.sqlite file inside the config folder and migrate to the new config file

## ✨ New Features

- Everything is now configured in a config.yaml file
- The database only stores runs and logs for jobs
- Changing the config file will automatically trigger a change in the GUI
- Healthchecks are now supported to be send before and after jobs or in case of an error
- Pre- and post-commands can be defined in the defaults section
- Software is installed based on the needs of the user at docker start
- A new Terminal where any command can be executed inside the docker (Use carefully!)

## 🐛 Bug Fixes

- Buggy notifications have been removed
- Race-condition in GUI when loading the job detail page

---

Generated changelog:
{{ changelog }}
