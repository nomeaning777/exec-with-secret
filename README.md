# exec-with-secret

Run commands with secrets from GCP Secret Manager.
This program resolves environment variables that starts with `secretmanager://`.

## Usage

```shell
$ NON_SECRET=not_secret SECRET=secretmanager://<PROJECT_NAME>/<SECRET_NAME>/<VERSION> ./exec-with-secret sh -c 'echo $NON_SECRET $SECRET'
not_secret very_important_secret
```

## License
MIT
