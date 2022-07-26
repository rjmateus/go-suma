NOTE: hackweek level code. No test, taking several shortcuts, didn't spend much time on code organization.

# Suma Update Manager

Allow the management of suma server, monitory status and apply updates

It also exposes suma repositories using only the File system and the postgresql database.

# Run
`go build`
Copy the artifact to sumas server
Run it `./suma-updater`

port 8088 will expose the API

## Available methods

| name          | http Method | description                        | example                                                                                                                   |
|---------------|-------------|------------------------------------|---------------------------------------------------------------------------------------------------------------------------|
| status        | GET         | Return server status               | curl localhost:8088/api/status                                                                                            |
| serviceStatus | GET         | Return server service status       | curl localhost:8088/api/serviceStatus                                                                                     |
| updates       | GET         | List of available updates          | curl localhost:8088/api/updates                                                                                           |
| patches       | GET         | List of available patches          | curl localhost:8088/api/patches                                                                                           |
| refresh       | POST        | cal `zypper ref -f`                | curl -X POST localhost:8088/api/refresh                                                                                   |
| patch         | POST        | options: withOptional, withUpdates | curl -X POST localhost:8088/api/patch -H "Content-Type: application/json" -d '{"withUpdate": true, "withOptional": true}' |
| updatePackage | POST        | options: packages                  | curl -X POST localhost:8088/api/updatePackage -H "Content-Type: application/json" -d '{"packages": ["packageName"]}'      |

## repository

One can use the same suma repositories serving endpoints but with port `8088` instead.
example: 
`http://localhost:8088/rhn/manager/download/sle-module-basesystem15-sp3-updates-x86_64`

# TODO
- [ ] repo access verification. Now works as if check-tokens where set to false (property `java.salt_check_download_tokens`)
- [ ] Run updates in a go routine, to be non-blocking. Add mechanisms to monitory the update progress.
- [ ] Allow to install of a single patch or update