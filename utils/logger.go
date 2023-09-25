package utils

import "log"

var logger log.Logger

func InitLogger() {

}
func GetLogger() *log.Logger {
	return &logger
}
func LogInfo(message string) {
	logger.Println(message)
}
func LogError(message string) {
	logger.Println("Error, ", message)
}
