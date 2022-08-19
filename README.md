
# Entrega do desafio de GoLang do [Imersão fullcycle](https://imersao.fullcycle.com.br/desafio/Imersao30/desafio1-golang)

## Informações do desafio

Neste desafio, você deverá criar uma aplicação Golang que realiza transações bancárias.

A aplicação será uma API Rest contendo 2 endpoints:


- POST /bank-accounts - Criar contas bancárias


No corpo da requisição deverá ser enviado:

```http
{
"account_number": "1111-11"
}
```

A resposta HTTP deverá ser 201, contendo o ID da conta criada e o "account_number"


- POST /bank-accounts/transfer - Transferência entre contas bancárias


No corpo da requisição deverá ser enviado:

```http
{
"from": "1111-11"
"to": "1111-11"
"amount": 100
}
```

A aplicação deverá persistir os dados no banco de dados SQLite.

A resposta HTTP deverá conter o saldo da conta from e o saldo da conta to.

Disponibilize esta aplicação Golang com Docker Compose na porta 8000.

Ao rodar docker compose up todo ambiente deverá já estar disponível.

# Bibliotecas usadas

* [Gin](github.com/gin-gonic/gin): para fazer a API
* [go-sqlite3](github.com/mattn/go-sqlite3): para lidar com o SQLite

# Como rodar a aplicação:

```sh
docker compose up
```

# Como testar:

## Criar conta

```sh
curl -X POST -H "Content-Type: application/json" -d '{"account_number":"1111-11"}' http://localhost:8000/bank-accounts
```

## Transferir saldo

```sh
curl -X POST -H "Content-Type: application/json" -d '{"from":"2222-22", "to": "1111-11", "amount": 100}' http://localhost:8080/bank-accounts/transfer
```