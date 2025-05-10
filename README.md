# go-app

A simple Rest API made in GO with session based authentication.

It was made to learn the GO programming language.

## Endpoints

Base path is `/api/v1/auth`

- `POST /register` - Register a new user
- `POST /login` - Login and receive a session token
- `GET /session` - Retrieve the current session information
- `POST /logout` - Logout and invalidate the session token

## Run

- Clone the repository

```bash
git clone git@github.com/Zigl3ur/go-app.git
cd go-app
```

- Install dependencies

```bash
make install_deps
```

- Migrate and seed database

```bash
make migrate
make seed
```

- Run the project

```bash
make run
```
