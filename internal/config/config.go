package config

import (
	"log"
	"os"
	"strconv"
)

var Conf *Config

type Config struct {
	AppPort   string
	AgentPort string

	PlusTime     int
	MinusTime    int
	MultipTime   int
	DivisionTime int

	ComputingPower int
}

func getEnv(arg, value string) int {
	envVar := os.Getenv(arg)
	if envVar == "" {
		envVar = value
	}

	retVar, err := strconv.Atoi(envVar)
	if err != nil {
		log.Fatal("failed to convert env variable to int")
		panic(err)
	}

	return retVar
}

func ConfigFill() *Config {
	config := new(Config)

	// env
	config.AppPort = os.Getenv("APP_PORT")
	config.AgentPort = os.Getenv("AGENT_PORT")

	config.PlusTime = getEnv("TIME_ADDITION_MS", "100")
	config.MinusTime = getEnv("TIME_SUBSTRACTION_MS", "100")
	config.MultipTime = getEnv("TIME_MULTIPLICATIONS_MS", "100")
	config.DivisionTime = getEnv("TIME_DIVISIONS_MS", "100")

	config.ComputingPower = getEnv("COMPUTING_POWER", "5")
	if config.AppPort == "" {
		config.AppPort = "8080"
	}

	if config.AgentPort == "" && config.AppPort != "8081" {
		config.AgentPort = "8081"
	} else {
		config.AgentPort = "8082"
	}

	return config
}
