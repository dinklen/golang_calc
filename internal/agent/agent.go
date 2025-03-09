package agent

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
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

	w.Header().Set("Content-Type", "application/json")

	exprs := expressions.NewExpressions()
	retExprs := NewSafetyMap()

	data, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] failed to read body")
		w.WriteHeader(500)
		return
	}

	err = json.Unmarshal(data, &exprs)
	if err != nil {
		log.Printf("[ERROR] failed to unmarshal data")
		w.WriteHeader(500)
		return
	}

	var ec bool = false
	var wg sync.WaitGroup
	sem := make(chan struct{}, config.Conf.ComputingPower)

	for exprsCounter := 0; exprsCounter < len(exprs.Exprs); exprsCounter++ {

		wg.Add(1)

		go func(index int, wg *sync.WaitGroup, sem chan struct{}) {
			var result float64

			defer wg.Done()

			sem <- struct{}{}

			arg1, err1 := strconv.ParseFloat(exprs.Exprs[index].Arg1, 64)
			arg2, err2 := strconv.ParseFloat(exprs.Exprs[index].Arg2, 64)

			if err1 != nil || err2 != nil {
				log.Printf("[ERROR] worker %d: parsing error (err1: %v; err2: %v)", index+1, err1, err2)
				ec = true
				retExprs.Add(fmt.Sprintf("%v", math.MaxFloat64), exprs.Exprs[index].Id)
				<-sem
				return
			}

			switch exprs.Exprs[index].Operation {
			case "*":
				time.Sleep(time.Duration(config.Conf.MultipTime) * time.Millisecond)
				result = arg1 * arg2
			case "/":
				if arg2 == 0 {
					log.Printf("[ERROR] worker %d: %v", index+1, errors.ErrDivisionByZero)
					ec = true
					retExprs.Add(fmt.Sprintf("%v", math.MaxFloat64), exprs.Exprs[index].Id)
					<-sem
					return
				}

				time.Sleep(time.Duration(config.Conf.DivisionTime) * time.Millisecond)
				result = arg1 / arg2
			case "+":
				time.Sleep(time.Duration(config.Conf.PlusTime) * time.Millisecond)
				result = arg1 + arg2
			case "-":
				time.Sleep(time.Duration(config.Conf.MinusTime) * time.Millisecond)
				result = arg1 - arg2
			}

			retExprs.Add(fmt.Sprintf("%v", result), exprs.Exprs[index].Id)

			log.Printf("[INFO] worker %d: success", index+1)

			<-sem
		}(exprsCounter, &wg, sem)
	}

	wg.Wait()

	if ec {
		w.WriteHeader(422)
		return
	}

	retData, err := json.Marshal(retExprs)
	if err != nil {
		log.Printf("[ERROR] failed to convert data to json: %v", err)
		w.WriteHeader(500)
		return
	}

	_, err = w.Write(retData)

	if err != nil {
		log.Printf("[ERROR] failed to send message to orchestrator: %v", err)
		w.WriteHeader(500)
		return
	}

	log.Printf("[INFO] subcalc success")

	w.WriteHeader(200)
}
