# addressbook
### Config
`config.json` file has 3 parameters that need to be configured:
- JWT signing secret
- PostgreSQL URL
- HTTP server port (if you change this and use Docker Compose then remember to change the exposed port in docker-compose.yml as well)

### Schema
Running addressbook will make it automatically try to run the `schema.sql` on the database. API will be served regardless of whether schema initialization fails or succeeds.

### Running
addressbook can be run either manually or using Docker Compose.
To run manually, bring up your PostgreSQL, run `go build` to build the project, then run the resulting executable `addressbook`.
To run with Docker Compose run `sudo docker-compose up`.

### API
#### User
/api/user POST -> Create user
/api/user/token POST -> Create token

There are two user endpoints, one for registering and one for obtaining JWT tokens.
You can register by POSTing to `/api/user` with JSON payload containing "username", "email", and "password" fields. Password can't be shorter than 6 characters and email has to be in valid format.
You can create and obtain a new JWT token by POSTing to `/api/user/token` with "username" (or "email") and "password" JSON fields. All subsequent API endpoints expect you to send this token in header as `Authorization: Bearer [token]`.

All operations on contacts and contact-lists can only be done by the user that has created them.

#### Contact
/api/contact POST -> Create contact
/api/contact/{id} DELETE -> Delete contact
/api/contact GET -> Get contacts
/api/contact/{id} GET -> Get contact

When creating a contact you should pass in a JSON payload with fields "name", "surname", and "email".

#### Contact-list
/api/contact-list POST -> Create contact-list
/api/contact-list/{id} DELETE -> Delete contact-list
/api/contact-list GET -> Get contact-lists
/api/contact-list/{id} GET -> Get contact-list
/api/contact-list/search POST -> Search contact-lists by name
/api/contact-list/{id}/contact GET -> List contacts of contact-list
/api/contact-list/{id}/contact POST -> Add contact to contact-list
/api/contact-list/{id}/contact DELETE -> Delete a contact from contact-list

When creating a contact-list you should pass in a JSON payload with field "name".
When searching for contact-lists by name you should pass in a JSON payload with field "term", referring to search term.
When adding/deleting a contact to/from contact-list you should pass in a JSON payload with field "id", referring to contact ID. Also note that "id" should be of JSON Number type.
