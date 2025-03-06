package database

import (
	"database/sql"
	"log"

	"golang_calc/internal/calc_libs/expressions"

	_ "github.com/lib/pq"
)

type Database struct {
	DB interface{}
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "user=postgres dbname=subexpressions sslmode=disable")
	if err != nil {
		log.Fatal("failed to create database object: ", err)
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (db *Database) Insert(expr *expressions.Expression, enabled1, enabled2 bool) error {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return err
	}

	_, err = db.DB.Exec(
		`INSERT INTO calc_exprs (id, num1, operator, num2, used, enb1, enb2)
		VALUES (\$1, \$2, \$3, \$4, \$5, \$6, \$7)`,

		expr.Id,
		expr.Arg11,
		expr.Operation,
		expr.Arg2,
		false,
		enabled1,
		enabled2,
	)

	if err != nil {
		log.Printf("[ERROR] failed to insert data to database: ", err)
		return err
	}

	log.Printf("[INFO] inserting success")
	return nil
}

/*
func (db *Database) GetResult(id string) (float64, error) {
	var result float64

	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return err
	}

	row := db.DB.QueryRow(`SELECT result FROM calc_exprs WHERE id = \$1`, id)
	err = row.Scan(&result)

	if err != nil {
		log.Printf("[ERROR] failed get data from database: ", err)
		return err
	}

	log.Printf("[INFO] result got")
	return result, nil
}
*/

func (db *Database) UnloadTasks() (*expressions.Expressions, error) {
	tasks := []*expressions.Expression{}

	if err := db.DB.Ping(); err != nil {
		log.Fatal("failed to ping database: ", err)
		return err
	}

	rows, err := db.DB.Query(`
		SELECT (id, num1, operator, num2)
		FROM calc_exprs
		WHERE enb1 = TRUE AND enb2 = TRUE AND used = FALSE;
	`)

	if err != nil {
		log.Printf("[ERROR] failed to select tasks from database: ", err)
		return []*expressions.Expression{}, err
	}

	for rows.Next() {
		tasks = append(tasks, expressions.NewExpression(
			id,
			num1,
			num2,
			operator,
		))
	}

	return tasks, nil
}

func (db *Database) UpdateValues() error {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return err
	}

	// select items where unk_nums >= 1 and write logic
	rows, err := db.DB.Query(`
		SELECT (id, num1, num2)
		FROM calc_exprs
		WHERE (enb1 = FALSE OR enb2 = FALSE) AND used = FALSE;
	`)

	if err != nil {
		log.Fatal("failed to update database: ", err)
		return err
	}

	args := make([]string, 2)
	arg := ""

	var id int

	for rows.Next() {
		rows.Scan(&id, &args[0], &args[1]) // init id

		// parse expression id (for example: "{432}" => "432")
		for index := 0; index < 2; index++ {
			if []rune(args[index])[0] == '{' {
				args[index][len(args[index])-1] = ""
				args[index][0] = ""
			}

			args[index] = arg
			arg = ""
		}

		if args[0] != "" {
			_, err = db.DB.Exec(
				`
				UPDATE calc_exprs
				SET num1 = (
					SELECT result FROM calc_exprs WHERE id = \$1
				) WHERE id = \$2
				`,
				args[0],
				id,
			)
		}

		if args[1] != "" {
			_, err = db.DB.Exec(
				`
				UPDATE calc_exprs
				SET num2 = (
					SELECT result FROM calc_exprs WHERE id = \$1
				) WHERE id = \$2
				`,
				args[1],
				id,
			)
		}

		if err != nil {
			log.Fatal("[FATAL] failed to update database: ", err)
			return err
		}
	}

	log.Printf("[INFO] update success")
	return nil
} // update all values status {enable} from database

func (db *Database) UnloadAllTasks() ([]*expressions.ExpressionInfo, error) {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return nil, err
	}

	retExprs := []*expressions.ExpressionInfo{}

	rows, err := db.DB.Query("SELECT (id, result, used, enb1, enb2) FROM calc_exprs")
	if err != nil {
		log.Printf("[ERROR] failed to unload all tasks from database: %v", err)
		return nil, err
	}

	var (
		id  int
		res float64

		enb1       bool
		enb2       bool
		usedStatus bool

		totalStatus string
	)

	for rows.Next() {
		rows.Scan(&id, &res, &usedStatus, &enb1, &enb2)

		if !enb1 || !enb2 {
			totalStatus = "waiting"
		} else if !usedStatus {
			totalStatus = "calculating"
		} else {
			totalStatus = "calculated"
		}

		retExprs = append(
			retExprs,
			&expressions.ExressionInfo{
				Id:     id,
				Status: totalStatus,
				Result: res, // ?? а если его ещё нет?
			},
		)
	}

	log.Printf("[INFO] success to unload all tasks")
	return retExprs, nil
}

func (db *Database) UpdateUsedStatus(ids map[string]int) error {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return err
	}

	for key, value := range ids {
		_, err := db.DB.Exec(
			`
			UPDATE calc_exprs
			SET result = \$1,
				used = TRUE
			WHERE id = \$2;
			`,
			key,
			value,
		)

		if err != nil {
			log.Printf("[ERROR] failed to update used status and result in database: %v", err)
			return err
		}
	}

	log.Printf("[INFO] success to update status and result value")
	return nil
}

func (db *Database) UnloadCurrentTask(id int) (*expressions.Expression, error) {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return nil, err
	}

	row, err := db.DB.QueryRow("SELECT (id, result, used, enb1, enb2) FROM calc_exprs")

	if err != nil {
		log.Printf("[ERROR] failed to unload all tasks from database: %v", err)
		return nil, err
	}

	var (
		id  int
		res float64

		enb1       bool
		enb2       bool
		usedStatus bool

		totalStatus string
	)

	row.Scan(&id, &res, &usedStatus, &enb1, &enb2)

	if !enb1 || !enb2 {
		totalStatus = "waiting"
	} else if !usedStatus {
		totalStatus = "calculating"
	} else {
		totalStatus = "calculated"
	}

	log.Printf("[INFO] success to unload current task")
	return &expressions.ExpressionInfo{
		Id:     id,
		Status: totalStatus,
		Result: res,
	}, nil
}

/*
func connect() error {
	connectData := "user=postgres dbname=subexpressions sslmode=disable"
	db, err := sql.Open("postgres", connectData)
	if err != nil {
		log.Fatal("connection to database failed: ", err)
		return err
	}

	defer db.Close()

	err = db.Ping()
	if err != nil {
		log.Fatal("database ping failed: ", err)
		return err
	}

	log.Info("success to connect to database")
}

func Run(command string) error { // so dangerous!!!
	err := connect()
	if err != nil {
		return err
	}
}

connectionData := "user=postgres dbname=subexpressions sslmode=disable"
	db, err := sql.Open("postgres", connectionData)
	if err != nil {
		log.Fatal("failed to connect to database: ", err)
		return err
	}

	defer db.Close()
*/

// load database from file?
