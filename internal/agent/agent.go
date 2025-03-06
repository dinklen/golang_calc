package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"

	"golang_calc/internal/application"
	"golang_calc/internal/calc_libs/errors"
	"golang_calc/internal/calc_libs/expressions"
)

type safetyMap struct {
	Map map[string]int `json:"map"`
	mu  sync.RWMutex
}

func NewSafetyMap() *safetyMap {
	return &safetyMap{
		Map: make(map[string]int),
	}
}

func (sm *safetyMap) Add(expr string, id int) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.Map[expr] = id
}

func (sm *safetyMap) Get(expr string) int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	return sm.Map[expr]
}

func AgentHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

	exprs := expressions.NewExpressions()
	retExprs := NewSafetyMap()
	//make([]endExprs, len=exprs) и в них добавлять в форе

	data, err := io.ReadAll(r.Body)
	if err != nil {
		// error_output
	}

	err = json.Unmarshal(data, &exprs)
	if err != nil {
		// error output
	}

	cp, err := strconv.Atoi(application.App.Configuration.ComputingPower)         // error_output
	minusTime, err := strconv.Atoi(application.App.Configuration.MinusTime)       // error_output
	plusTime, err := strconv.Atoi(application.App.Configuration.PlusTime)         // error_output
	multipTime, err := strconv.Atoi(application.App.Configuration.MultipTime)     // error_output
	divisionTime, err := strconv.Atoi(application.App.Configuration.DivisionTime) // error_output

	var wg sync.WaitGroup
	sem := make(chan struct{}, cp)

	for exprsCounter := 0; exprsCounter < len(exprs.Exprs); exprsCounter++ {
		wg.Add(1)

		go func(index int, wg *sync.WaitGroup, sem chan struct{}) {
			var result float64

			defer wg.Done()

			sem <- struct{}{}

			switch exprs.Exprs[index].Operator {
			case "*":
				time.Sleep(time.Duration(multipTime) * time.Millisecond)
				result = exprs.Exprs[index].Arg1 * exprs.Exprs[index].Arg2
			case "/":
				if exprs.Exprs[index].Arg2 == 0 {
					log.Printf("[ERROR] ", errors.ErrDivisionByZero)
					return
				}
				time.Sleep(time.Duration(divisionTime) * time.Millisecond)
				result = exprs.Expressions[index].Arg1 / exprs.Expressions[index].Arg2
			case "+":
				time.Sleep(time.Duration(plusTime) * time.Millisecond)
				result = exprs.Expressions[index].Arg1 + exprs.Expressions[index].Arg2
			case "-":
				time.Sleep(time.Duration(minusTime) * time.Millisecond)
				result = exprs.Expressions[index].Arg1 - exprs.Expressions[index].Arg2
			}

			retExprs.Add(fmt.Sprintf("%s", result), exprs.Expressions[index].Id)

			log.Printf("[INFO] worker %d: success", index)

			<-sem
		}(exprsCounter, &wg, sem)
	}

	wg.Wait()

	err = json.NewEncoder(w).Encode(retExprs.Map)

	if err != nil {
		log.Printf("[ERROR] failed to send message to orchestrator: %v", err)
	}

	log.Printf("[INFO] subcalc success")

	w.WriteHeader(http.StatusOK)
}
