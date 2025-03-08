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

func getEnv(arg string) int {
	envVar := os.Getenv(arg)
	if envVar == "" {
		envVar = "100"
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

	config.PlusTime = getEnv("TIME_ADDITION_MS")
	config.MinusTime = getEnv("TIME_SUBSTRACTION_MS")
	config.MultipTime = getEnv("TIME_MULTIPLICATION_MS")
	config.DivisionTime = getEnv("TIME_DIVISION_MS")

	config.ComputingPower = getEnv("COMPUTING_POWER")
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
