package tests

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/yonisaka/queue-worker/config"
	"github.com/yonisaka/queue-worker/utils"
	"log"
)

type Test struct {
	Ctx context.Context
	Cfg *config.Config
}

func Init() *Test {
	if err := godotenv.Load(fmt.Sprintf("%s/.env", utils.RootDir())); err != nil {
		log.Fatalf("no .env file provided.")
	}

	ctx := context.Background()
	cfg := config.New()

	return &Test{
		Ctx: ctx,
		Cfg: cfg,
	}
}
