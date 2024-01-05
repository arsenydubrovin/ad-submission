package main

import (
	"fmt"

	"github.com/arsenydubrovin/ad-submission/internal/config"
)

func main() {
	cfg := config.Load()

	fmt.Println(cfg)
}
