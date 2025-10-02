// Copyright 2015 Nevio Vesic
// Please check out LICENSE file for more information about what you CAN and what you CANNOT do!
// Basically in short this is a free software for you to do whatever you want to do BUT copyright must be included!
// I didn't write all of this code so you could say it's yours.
// MIT License

package goesl

import "log"

var logger *log.Logger

func SetLogger(l *log.Logger) {
	logger = l
}

func Debug(message string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Printf(message, args...)
}

func Error(message string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Printf(message, args...)
}

func Notice(message string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Printf(message, args...)
}

func Info(message string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Printf(message, args...)
}

func Warning(message string, args ...interface{}) {
	if logger == nil {
		return
	}
	logger.Printf(message, args...)
}
