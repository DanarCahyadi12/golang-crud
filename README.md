
# Golang CRUD

Golang CRUD operations with JWT authentication



## Tech Stack

- Golang

- MySQL



## Framework & Library

- Fiber (Http framework)

- GORM (ORM)

- Viper (Configuration)

- Golang Migrate (Database Migration)

- Go Playground Validator (Validation) 
## Configuration

All config is in `config.json` file.
## Run migrations

```bash
 migrate -database "mysql://<your_username>:<your_password>@tcp(<your_host>:<your_port>)/<your_database>?charset=utf8mb4&parseTime=true&loc=Local" -path database/migrations up
```
    

## Run application

```bash
go run cmd/web/main.go
```
## API Reference

#### Sign up

```http
  POST /signup
```

| Body field | Type     | Description                |
| :-------- | :------- | :------------------------- |
| `name` | `string` | Required  |
| `email`| `string` | Requried, max 255 character |
| `password` | `string` | Required, minimum 8 character |

#### Sign in

```http
  POST /auth/signin
```

| Body field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `email`      | `string` | Required |
| `password` | `string` | required |


#### Sign out

```http
  GET /signout
```

#### Get access token

```http
  POST /auth/token
```

#### Create product

```http
  POST /products
```

| Headers | Description | Value |
| :--------- | :------- | :----------|
| `Authorization` | `Type` :`Bearer token` | `Bearer <YOUR_ACCESS_TOKEN` |

| Body field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`      | `string` | Required |
| `price` | `number` | required |
| `stock` | `number` | required 

#### Update product

```http
  PUT /products/:id
```

| Headers | Description | Value |
| :--------- | :------- | :----------|
| `Authorization` | `Type` :`Bearer token` | `Bearer <YOUR_ACCESS_TOKEN` |

| Body field | Type     | Description                       |
| :-------- | :------- | :-------------------------------- |
| `name`      | `string` | Required |
| `price` | `number` | required |
| `stock` | `number` | required 

#### Delete product

```http
  DELETE /products/:id
```

| Headers | Description | Value |
| :--------- | :------- | :----------|
| `Authorization` | `Type` :`Bearer token` | `Bearer <YOUR_ACCESS_TOKEN` |

#### Get products

```http
  GET /products
```

| Headers | Description | Value |
| :--------- | :------- | :----------|
| `Authorization` | `Type` :`Bearer token` | `Bearer <YOUR_ACCESS_TOKEN` |


Query params
| Key | Description | Type |
| :--------- | :------- | :----------|
| `page` | `default value` : `1` | `number` |
| `limit` | `default value` : `50` | `number`|

#### Get detail products

```http
  GET /product/:id
```

| Headers | Description | Value |
| :--------- | :------- | :----------|
| `Authorization` | `Type` :`Bearer token` | `Bearer <YOUR_ACCESS_TOKEN` |


