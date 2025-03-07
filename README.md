# golang_calc

---

Можно пока 0 не ставить? Мне серьёзно очень немного осталось, скоро сюда запушу всё.

**UPD**: прогресс есть, я думаю, этой ночью всё запушу. Огромное спасибо, если дал шанс)

**UPD2**: ...чуть времени ночью не хватило, скоро запушу

tg: @dinklen08

---

## О проекте
Это сервис, принимающий арифметические выражения и посылающий их на микросервис (агент). Он состоит из оркестратора, агента и базы данных.

---

## Схема и описание работы

---

## Запуск

1. Для начала необходимо клонировать этот репозиторий к себе в директорию:
```bash
git clone https://github.com/dinklen/golang_calc.git

cd golang_calc
```

2. Следующим этапом идёт объявление переменных окружения (не обязательно, по умолчанию задержка у всех будет равна *100 ms*, порт у оркестратора и агента - 8080 и 8081 соответственно, а максимальное кол-во горутин - *5*).
```bash
export APP_PORT=8080
export AGENT_PORT=8081
export TIME_ADDITION_MS=100
export TIME_SUBTRACTION_MS=100
export TIME_MULTIPLICATIONS_MS=100
export TIME_DIVISIONS_MS=100
```

| **Название** | **Описание** | **Значение по умолчанию** |
| :---: | :---: | :---: |
| `APP_PORT` | Порт, на котором запускается основной сервер (оркестратор) | 8080 |
| `AGENT_PORT` | Порт, на котором запускается микросервис, на котором происходят вычисления подвыражений (агент) | 8081 |
| `TIME_ADDITION_MS` | Задержка операции сложения в миллисекундах | 100 |
| `TIME_SUBTRACTION_MS` | Зажержка операции вычитания в миллисекундах | 100 |
| `TIME_MULTIPLICATIONS_MS` | Зажержка операции умножения в миллисекундах | 100 |
| `TIME_DIVISIONS_MS` | Зажержка операции деления в миллисекундах | 100 |
| `COMPUTING_POWER` | Максимальное количество горутин, которых агент может использовать одновременно | 5 |

> В примере приведены значения по умолчанию. Менять их необязательно.


3. Затем можно производить запуск (первоначально нужно запустить оркестратор, затем уже агент (как микросервис)).

На ***Linux***:
```bash
# запуск оркестратора в первом терминале
go run cmd/orchestrator/main.go

# запуск агента во втором терминале
go run cmd/agent/main.go
```

На ***Windows***:

Те же команды, только заменить слеши на обратные.

---

## Использование
В зависимости от запроса будут определённые статус-коды и сообщения.

**Статус-код** | **Сообщение** | **Описание**
:---: | :---: | :---:
`404` | `Page not found` | Введён некорректный адрес отправки запроса
`405` | `Access denied` | Использован некорректный метод
`422` | `Invalid input` | Отправлено некорректное выражение
`500` | `Internal server error` | Внутренняя ошибка сервера
`201` | `Received` | Данные успешно получены
`200` | `OK` | Результат выполнения задачи - ОК

### Позитивные сценарии
Для взаимодействия с калькулятором необходимо открыть третий терминал. Далее, для отправки выражения ввести команду:
```bash
curl --location 'localhost/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": <ваше выражение>
}'
```

Результат будет в формате **JSON**:
```json
{
    "id":<id выражения>
}
```

Это - id выражения. Для получения результата необходимо сделать ещё один запрос:
```bash
curl --location 'localhost/api/v1/expressions/:<id>'
```

Результат также будет в формате **JSON**:
```json
{
    "id":<id>
    "status":<"calculating"/"calculated">
    "result":<результат вычислений>
}
```

Но при этом можно получить список всех подвыражений:
```bash
curl --location 'localhost/api/v1/expressions/'
```

Тогда ответ будет таким:
```json
{
    "expressions":[
        {
            "id":<id>
            "status":<"calculating"/"calculated">
            "result":<результат вычислений/"null">
        },
        {
            "id":<id>
            "status":<"calculating"/"calculated">
            "result":<результат вычислений/"null">
        }
    ]
}
```

### Негативные сценарии
1. Введя некорректное выражение и затем запросив его результат (после подсчётов) мы получим следующее:
```json
{
    "id":<id>
    "status":"error"
    "result":"null"
}
```

В логах можно будет увидеть причину (логируется всё в терминал). Код ответа, в таком случае, - `422`.

2. При вводе некорректного адреса для отправки метода сразу же получим вывод с статус кодом `404`: `page not found`

3. Попытавшись использовать неразрешённый метод, получим результат с кодом `405`:`Access denied`

4. Если что-то пойдёт не так во время обработки вычислений/отправки сообщений/..., то мы получим результат с кодом `500`:`Internal server error`

Все будет получено в формате **JSON**.

---

**Telegram** | [@dinklen08](https://t.me/@dinklen08)
