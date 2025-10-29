package handlers

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/joaocezar1/integracaoBalancaPdv/configs"
)

func ConfigHandler(cfg *configs.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
			return
		}
		//A requisição http pode vir com charset, usar o HasPrefix deixa mais robusto
		if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
			http.Error(w, "Content-Type deve ser application/json", http.StatusUnsupportedMediaType)
			return
		}

		var b configs.Balanca
		if err := json.NewDecoder(r.Body).Decode(&b); err != nil {
			http.Error(w, "Erro ao decodificar JSON", http.StatusBadRequest)
			return
		}

		cfg.Balanca = &b

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Configurações da balança atualizadas com sucesso"))
	}
}
