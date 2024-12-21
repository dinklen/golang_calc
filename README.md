# golang_calc
The repository includes a function for calculating mathematical expressions.

---

<h2>Functional</h2>
This app receives a mathematical expression from the user in the

**JSON format**
. After processing it, it returns the answer to the user (also in the **JSON format**). After starting the application, the server will be launched along the path
`localhost:8080/api/v1/calculate`

> Disadvantage - does not accept negative numbers. The result itself can be negative, but for example, the expression `-1+6` or `4*(7-9)` will return an error.

<h2>How to use it?</h2>
Before using it, you should make sure that Golang and Git are installed on your computer.

___

To use it, you should follow these steps:

1. Cloning the repository

    ```
    git clone https://github.com/dinklen08/golang_calc/
    ```

2. Start the app

    On **\*Unix**:
    ```
    go run golang_calc/cmd/main.go
    ```

    On **Windows**:
    ```
    go run golang_calc\cmd\main.go
    ```

And then with the help of cURL submit requests to him. For example (send it in other terminal/command prompt):

**\*Unix**

```
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "*your expression*"
}'
```

**Windows**

```
 curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\"expression\": \"*expression*\"}"
```

...result (in **JSON format**):

```
{"result":"6"}
```

<h2>Program behavior</h2>
But this server will not always be able to answer you with any number. At least, when submitting incorrect data. Below are the status codes that characterize the server's response.

___

| **Status Code** | **Output** |
| :---: | :---: |
| 200 | Status - OK. The result will be returned as a number. `{"result":"*number*"}` |
| 405 | Status - invalid method. The only method allowed is POST. `{"error":"Access denied"}` |
| 422 | Status - expression error. An invalid expression was supplied to the input. `{"error":"Expression is not valid: *description*"}` |
| 500 | Status - server error. An error occurred while processing your request. `{"error":"Internal server error: *description*"}` |

___
