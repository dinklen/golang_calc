package database

import (
	"database/sql"
	"log"
	"math"

	"golang_calc/internal/calc_libs/errors"
	"golang_calc/internal/calc_libs/expressions"

	_ "github.com/lib/pq"
)

var DataBase *Database

type Database struct {
	DB *sql.DB
}

func NewDatabase() (*Database, error) {
	db, err := sql.Open("postgres", "user=postgres password=db78903 dbname=subexpressions sslmode=disable")
	if err != nil {
		log.Fatal("[FATAL] failed to create database object: ", err)
		return nil, err
	}

	return &Database{DB: db}, nil
}

func (db *Database) Insert(expr *expressions.Expression, enabled1, enabled2 bool) error {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return err
	}

	query := `INSERT INTO calc_exprs (id, num1, operator, num2, enb1, enb2, used) VALUES ($1, $2, $3, $4, $5, $6, $7)`

	_, err := db.DB.Exec(
		query,
		expr.Id,
		expr.Arg1,
		expr.Operation,
		expr.Arg2,
		enabled1,
		enabled2,
		false,
	)

	if err != nil {
		log.Printf("[ERROR] failed to insert data to database: %v", err)
		return err
	}

	log.Printf("[INFO] inserting success")
	return nil
}

func (db *Database) UnloadTasks() (*expressions.Expressions, error) {
	tasks := []*expressions.Expression{}

	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return nil, err
	}

	rows, err := db.DB.Query(`
		SELECT id, num1, operator, num2
		FROM calc_exprs
		WHERE enb1 = TRUE AND enb2 = TRUE AND used = FALSE;
	`)

	if err != nil {
		log.Printf("[ERROR] failed to select tasks from database: %v", err)
		return nil, err
	}

	var (
		id       uint32
		num1     string
		num2     string
		operator string
	)

	for rows.Next() {
		err := rows.Scan(&id, &num1, &operator, &num2)
		if err != nil {
			log.Printf("[ERROR] failed to scan data from database: %v", err)
			return nil, err
		}

		tasks = append(tasks, expressions.NewExpression(
			id,
			num1,
			num2,
			operator,
		))
	}

	return &expressions.Expressions{Exprs: tasks}, nil
}

func (db *Database) UpdateValues() error {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return err
	}

	rows, err := db.DB.Query(`SELECT id, num1, num2 FROM calc_exprs WHERE (enb1 = FALSE OR enb2 = FALSE) AND used = FALSE`)

	if err != nil {
		log.Fatal("[FATAL] failed to update database: ", err)
		return err
	}

	args := make([]string, 2)

	var id uint32

	for rows.Next() {
		log.Printf("scaning...")
		rows.Scan(&id, &args[0], &args[1])

		for index := 0; index < 2; index++ {
			if rune(args[index][0]) == '{' {
				args[index] = args[index][1 : len(args[index])-1]
			} else {
				args[index] = ""
			}
		}

		if args[0] != "" {
			_, err = db.DB.Exec(
				`
				UPDATE calc_exprs
				SET num1 = CASE
				WHEN (SELECT result FROM calc_exprs WHERE id = $1) IS NOT NULL THEN
				(SELECT result::TEXT FROM calc_exprs WHERE id = $1)
				ELSE ('{' || $1 || '}')
				END,
				enb1 = (SELECT result FROM calc_exprs WHERE id = $1) IS NOT NULL
				WHERE id = $2
				`,

				args[0],
				id,
			)

			if err != nil {
				log.Fatal("[FATAL] failed to update database: ", err)
				return err
			}
		}

		if args[1] != "" {
			_, err = db.DB.Exec(
				`
				UPDATE calc_exprs
				SET num2 = CASE
				WHEN (SELECT result FROM calc_exprs WHERE id = $1) IS NOT NULL THEN
				(SELECT result::TEXT FROM calc_exprs WHERE id = $1)
				ELSE ('{' || $1 || '}')
				END,
				enb2 = (SELECT result FROM calc_exprs WHERE id = $1) IS NOT NULL
				WHERE id = $2
				`,

				args[1],
				id,
			)

			if err != nil {
				log.Fatal("[FATAL] failed to update database: ", err)
				return err
			}
		}
	}

	log.Printf("[INFO] update success %v", id)
	return nil
}

func (db *Database) UnloadAllTasks() ([]*expressions.ExpressionInfo, error) {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return nil, err
	}

	retExprs := []*expressions.ExpressionInfo{}

	rows, err := db.DB.Query("SELECT id, result, used, enb1, enb2 FROM calc_exprs")
	if err != nil {
		log.Printf("[ERROR] failed to unload all tasks from database: %v", err)
		return nil, err
	}

	var (
		id  uint32
		res float64

		enb1       bool
		enb2       bool
		usedStatus bool

		totalStatus string
	)

	for rows.Next() {
		rows.Scan(&id, &res, &usedStatus, &enb1, &enb2)

		if res == math.MaxFloat64 {
			totalStatus = "error"
		} else if !enb1 || !enb2 {
			totalStatus = "waiting"
		} else if !usedStatus {
			totalStatus = "calculating"
		} else {
			totalStatus = "calculated"
		}

		retExprs = append(
			retExprs,
			&expressions.ExpressionInfo{
				Id:     id,
				Status: totalStatus,
				Result: res,
			},
		)
	}

	log.Printf("[INFO] success to unload all tasks")
	return retExprs, nil
}

func (db *Database) InsertError(id uint32) error {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return err
	}

	_, err := db.DB.Exec(
		`INSERT INTO calc_exprs (id, result) VALUES ($1, $2)`,
		id,
		math.MaxFloat64,
	)

	if err != nil {
		log.Printf("[ERROR] failed to insert error expr status")
		return err
	}

	return nil
}

func (db *Database) UpdateUsedStatus(ids map[float64]uint32) error {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return err
	}

	for key, value := range ids {
		_, err := db.DB.Exec(
			`UPDATE calc_exprs SET result = $1, used = TRUE WHERE id = $2`,
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

func (db *Database) UnloadCurrentTask(id uint32) (*expressions.ExpressionInfo, error) {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return nil, err
	}

	row := db.DB.QueryRow(`SELECT id, result, used, enb1, enb2 FROM calc_exprs WHERE id = $1`, id)

	var (
		id_s uint32
		res  float64

		enb1       bool
		enb2       bool
		usedStatus bool

		totalStatus string
	)

	row.Scan(&id_s, &res, &usedStatus, &enb1, &enb2)

	if id_s == 0 {
		return nil, errors.ErrIncorrectQuery
	}

	if res == math.MaxFloat64 {
		totalStatus = "error"
	} else if !enb1 || !enb2 {
		totalStatus = "waiting"
	} else if !usedStatus {
		totalStatus = "calculating"
	} else {
		totalStatus = "calculated"
	}

	log.Printf("[INFO] success to unload current task")
	return &expressions.ExpressionInfo{
		Id:     id_s,
		Status: totalStatus,
		Result: res,
	}, nil
}

func (db *Database) Clean() error {
	if err := db.DB.Ping(); err != nil {
		log.Fatal("[FATAL] failed to ping database: ", err)
		return err
	}

	_, err := db.DB.Exec(`TRUNCATE TABLE calc_exprs`)
	if err != nil {
		log.Printf("[ERROR] failed to clean table: %v", err)
		return err
	}

	return nil
}
