package server

import "os"

const defaultDatabaseFile = "despensa.db"
const defaultServerPort = "8080"
const defaultServerAddress = ""

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultServerPort
	}
	return port
}

func GetDatabaseFile() string {
	databaseFile := os.Getenv("DATABASE_FILE")
	if databaseFile == "" {
		databaseFile = defaultDatabaseFile
	}
	return databaseFile
}

func GetAddress() string {
	serverAddress := os.Getenv("ADDRESS")
	if serverAddress == "" {
		serverAddress = defaultServerAddress
	}
	return serverAddress
}
