# addressbook
### Config
`config.json` file has 3 parameters that need to be configured:
- JWT signing secret
- PostgreSQL URL
- HTTP server port (if you change this and use Docker Compose then remember to change the exposed port in docker-compose.yml as well)

### Running
addressbook can be run either manually or using Docker Compose. 
To run manually, bring up your PostgreSQL and run `go build` to build the project, then run the resulting executable `addressbook`.
To run with Docker Compose run `sudo docker-compose up`.

### API
#### User
There are two user endpoints, one for registering and one for obtaining JWT tokens.
You can register by POSTing to `/api/user` with JSON payload containing `username`, `email`, and `password` fields. Password can't be shorter than 6 characters.
You can create and obtain a new JWT token by POSTing to `/api/user/token` with `username` and `password` JSON fields. All subsequent API endpoints expect you to send this token in header as `Authorization: Bearer [token]`.
