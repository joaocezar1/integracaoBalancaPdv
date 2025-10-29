package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/joaocezar1/integracaoBalancaPdv/configs"
	"github.com/joaocezar1/integracaoBalancaPdv/internal/balanca"
)

func ConsultaConfigHandler(cfg *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(cfg); err != nil {
			http.Error(w, "Erro ao converter JSON", http.StatusInternalServerError)
		}
	}
}

func ConsultaPesoHandler(cfg *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var peso *string
		if r.Method != http.MethodGet {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		if cfg.Balanca == nil ||
			cfg.Balanca.Modelo == "" ||
			cfg.Balanca.Protocolo == "" ||
			cfg.Balanca.Paridade == 0 ||
			cfg.Balanca.Baud == 0 {
			http.Error(w, "Balança não configurada corretamente", http.StatusBadRequest)
			return
		}

		peso, err := balanca.ObterPesoDaBalanca(cfg)
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao consultar balança: %v", err), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader((http.StatusOK))
		w.Write([]byte(*peso))
	}
}
