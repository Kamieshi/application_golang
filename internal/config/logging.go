package config

import log "github.com/sirupsen/logrus"

func InitLogger() {
	log.SetLevel(log.TraceLevel)
	log.SetFormatter(&log.TextFormatter{})
}
