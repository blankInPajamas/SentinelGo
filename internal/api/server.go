package api

import (
	"github.com/blankInPajamas/SentinelGo/internal/storage/memory"
	"github.com/gin-gonic/gin"
)

func GetEventHandler(store *memory.InMemoryStorage) gin.HandlerFunc {
	return func() {

	}
}

func GetAlertHandler(store *memory.InMemoryStorage) gin.HandlerFunc {
	return func() {

	}
}