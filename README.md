# Building & Deploying
## Taskfile
The project uses [Task](https://taskfile.dev/) to generate
[templ](https://github.com/a-h/templ) templates and to build the executable.

Use `task templ` to (re-)generate the templates and `task build` to build the
executable.

## Docker
The ideal way to run this project using Docker, is by using the included
`docker-compose.yml`, by running `docker-compose up` in the project directory.
