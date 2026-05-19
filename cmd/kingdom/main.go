package kingdom

import (
	"context"
	"log"
	"os/signal"
	"syscall"

	"github.com/Rioverde/mine/internal/config"
	"github.com/Rioverde/mine/internal/repo"
)

func main() {

	cfg, err := config.ReadConfig("config.yaml")
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	conn, err := repo.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	_ = cfg

}
