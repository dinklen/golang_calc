# golang_calc
The repository includes a function for calculating mathematical expressions

---

<h2>Functional</h2>
The repository contains a function 

`CalcHandler` - she processes the questions submitted to her in **JSON format**. When launched, a local server is started: `localhost:8080/api/v1/calculate`

<h2>How to use it?</h2>
To use this app, you should follow these steps:

1. Cloning the repository

```
git clone https://github.com/dinklen08/golang_calc/
```

2. Start the app ;)

On **Unix**:
```
go run golang_calc/cmd/main.go
```

On **Windows**:
```
go run golang_calc\cmd\main.go
```

And then with the help of cURL submit requests to him. For example (send it in other terminal/command prompt):

```
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```

...result (in **JSON format**):

```
{"result":"6"}
```

<h2>Program behavior</h2>
But this server will not always be able to answer you with any number. At least, when submitting incorrect data. Below are the status codes that characterize the server's response.

---

| **Status Codes** | **Output** |
| :---: | :---: |
| 200 | Status - OK. The result will be returned as a number. `{"result":"*number*"}` |
| 422 | Status - expression error. An invalid expression was supplied to the input. `{"error":"Expression is not valid"}` |
| 500 | Status - server error. An error occurred while processing your request. `{"error":"Internal server error"}` |
