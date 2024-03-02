package server

import "os"

const defaultDatabaseFile = "despensa.db"
const defaultServerPort = "8080"
const defaultServerAddress = ""
const defaultSessionKey = "057fafb6f84419aae83010c3855c27b1fee4865ceeae34fdc191f60c02c5d96"

func GetSessionKey() []byte {
	sessionKey := os.Getenv("SESSION_KEY")
	if sessionKey == "" {
		sessionKey = defaultSessionKey
	}
	return []byte(sessionKey)
}

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
