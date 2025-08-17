## ✨ New Features

- Possibility to disable cron scheduling

  ```yaml
  jobs:
    - name: 'Disabled cron'
      envs:
        - key: SLEEP_TIME_LONG
          value: '90'
      # this will disable the scheduling for this job
      disable_cron: true
      commands:
        - sleep ${SLEEP_TIME_LONG}
        - echo "Job Done!"
  ```
