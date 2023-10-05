package logging

import (
	log "github.com/sirupsen/logrus"
)

// creazione di un'interfaccia standard per gli errori
type Event struct{
	id int
	message string
}

type StandardLogger struct {
	*log.Logger
}

// Inizializzo il nuovo logger
func NewLogger() *StandardLogger {
	var baseLogger = log.New()

	var standardLogger = &StandardLogger{baseLogger}

	standardLogger.Formatter = &log.JSONFormatter{}

	return standardLogger
}

// Variabili per la traduzione del numero errore in messaggio
var (
	messageReceived = 	Event{ 1, "Received input data: %s"}
)

func (l * StandardLogger) MessageReceived(data string){
	l.Infof(messageReceived.message, data)
}