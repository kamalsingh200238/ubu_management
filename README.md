# Project ubu_management

One Paragraph of project description goes here

## MakeFile

run all make commands with clean tests

```bash
make all build
```

build the application

```bash
make build
```

run the application

```bash
make run
```

Create DB container

```bash
make docker-run
```

Shutdown DB container

```bash
make docker-down
```

live reload the application

```bash
make watch
```

run the test suite

```bash
make test
```

clean up binary from the last build

```bash
make clean
```

## Tasks

- [x] Create function to redirect
- [x] login page
- [x] login functionality for coordinator
- [ ] middleware for routes
- [ ] Home page redirects the user to either login page or their dashboard
- [ ] login page redirects the user if they are already logged in and if they have bad token or something then remove the cookies and then let them login
- [x] Add templ, Create a util to render templ easily
- [x] Create function to create and validate jwt easily
- [x] Create function to hash passwords and check if password matches
