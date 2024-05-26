package main

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"hammingCode/internal/pkg/hamming/delivery"
	"hammingCode/internal/pkg/hamming/usecase"
	"os"

	_ "hammingCode/docs"
)

// @title Вариант 11
// @description Отправка текста. Передаваемую информацию защитить передаваемую информацию [7,4]-кодом Хэмминга. Длина сегмента (X) 150 байт, период сборки сегментов (N) 2 секунд, вероятность ошибки (P) 8%, вероятность потери кадра (R) 1%.
// @contact.name Дмитрий Белозеров ИУ5-64Б
// @contact.url http://t.me/belozerovmsk
func main() {
	if err := run(); err != nil {
		os.Exit(1)
	}
}

func run() error {
	router := gin.Default()

	usecase := usecase.NewHammingUsecase()
	handler := handler.NewHammingHandler(usecase)

	{
		router.POST("/code", handler.DataLink)
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	}

	return router.Run("127.0.0.1:5000")
}
