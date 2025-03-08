package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"golang_calc/internal/calc_libs/errors"
	"golang_calc/internal/calc_libs/expressions"
	"golang_calc/internal/config"
)

type safetyMap struct {
	Map map[string]uint32 `json:"map"`
	mu  sync.RWMutex
}

func NewSafetyMap() *safetyMap {
	return &safetyMap{
		Map: make(map[string]uint32),
	}
}

func (sm *safetyMap) Add(expr string, id uint32) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	sm.Map[expr] = id
}

func (sm *safetyMap) Get(expr string) uint32 {
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

	var wg sync.WaitGroup
	sem := make(chan struct{}, config.Conf.ComputingPower)

	for exprsCounter := 0; exprsCounter < len(exprs.Exprs); exprsCounter++ {
		wg.Add(1)

		go func(index int, wg *sync.WaitGroup, sem chan struct{}) {
			var result float64

			defer wg.Done()

			sem <- struct{}{}

			switch exprs.Exprs[index].Operation {
			case "*":
				time.Sleep(time.Duration(config.Conf.MultipTime) * time.Millisecond)
				result = exprs.Exprs[index].Arg1 * exprs.Exprs[index].Arg2
			case "/":
				if exprs.Exprs[index].Arg2 == 0 {
					log.Printf("[ERROR] %v", errors.ErrDivisionByZero)
					return
				}
				time.Sleep(time.Duration(config.Conf.DivisionTime) * time.Millisecond)
				result = exprs.Exprs[index].Arg1 / exprs.Exprs[index].Arg2
			case "+":
				time.Sleep(time.Duration(config.Conf.PlusTime) * time.Millisecond)
				result = exprs.Exprs[index].Arg1 + exprs.Exprs[index].Arg2
			case "-":
				time.Sleep(time.Duration(config.Conf.MinusTime) * time.Millisecond)
				result = exprs.Exprs[index].Arg1 - exprs.Exprs[index].Arg2
			}

			retExprs.Add(fmt.Sprintf("%s", result), exprs.Exprs[index].Id)

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
