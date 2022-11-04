NOTE: hackweek level code.

# Suma Downloade

Extract the /rhn/download endpoint to an external service.
This endpoint is stateless. 

It extracts autentication from header (if authentication is enabled) and checks the JWT token.
After this, connects to database to retrive the path in the file system for the package to be downloaded.


# Run
`go build`
Copy the artifact to sumas server
Run it `./go-suma`

port 8088 will expose the API

## Chance apache httd config

TODO

## repository

One can use the same suma repository endpoints but with port `8088` instead.
example: 
`http://localhost:8088/rhn/manager/download/sle-module-basesystem15-sp3-updates-x86_64`

# TODO
- [ ] Handle debian packages
- [ ] repo access verification. Now works as if check-tokens where set to false (property `java.salt_check_download_tokens`)
- [ ] Download for media files
- [ ] Parameterize the folder location
- [ ] Automatic tests