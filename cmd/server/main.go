package main

import (
	"log"
	"net/http"

	"github.com/joaocezar1/integracaoBalancaPdv/configs"
	"github.com/joaocezar1/integracaoBalancaPdv/internal/handlers"
)

var cfg = configs.NovoConfig()

func main() {
	http.HandleFunc("/balanca/config", handlers.ConfigHandler(cfg))
	http.HandleFunc("/balanca/consultaPeso", handlers.ConsultaPesoHandler(cfg))
	http.HandleFunc("/balanca/consultaConfig", handlers.ConsultaConfigHandler(cfg))
	log.Fatal(http.ListenAndServe(":8090", nil))
}
