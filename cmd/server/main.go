package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/joaocezar1/integracaoBalancaPdv/configs"
	"github.com/tarm/serial"
	"go.bug.st/serial/enumerator"
)

var cfg = configs.NovoConfig()

func main() {
	http.HandleFunc("/balanca/config", configHandler)
	http.HandleFunc("/balanca/consulta", consultaHandler)
	log.Fatal(http.ListenAndServe(":8090", nil))
}

func configHandler(w http.ResponseWriter, r *http.Request) {
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

	cfg.Balanca = b

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configurações da balança atualizadas com sucesso"))
}

func consultaHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
		return
	}
	if cfg.Balanca.Modelo == "" ||
		cfg.Balanca.Protocolo == "" ||
		cfg.Balanca.Paridade == 0 ||
		cfg.Balanca.Baud == 0 {
		http.Error(w, "Balança não configurada corretamente", http.StatusBadRequest)
		return
	}

	peso := obterPesoDaBalanca()
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader((http.StatusOK))
	w.Write([]byte(peso))
}

func obterPesoDaBalanca() string {
	var peso string
	ports, err := enumerator.GetDetailedPortsList()
	if err != nil {
		log.Fatal(err)
	}

	if len(ports) == 0 {
		fmt.Println("Nenhuma porta serial encontrada.")
		return ""
	}
	var deviceName string
	for _, port := range ports {
		if port.Product != "" &&
			(port.Product == "Toledo Virtual COM Port" || port.VID == "10954") {
			deviceName = port.Name
			break
		}
	}
	if deviceName == "" {
		fmt.Println("Balança Toledo não detectada.")
		return ""
	}
	var parity serial.Parity
	if cfg.Balanca.Paridade == 1 {
		parity = serial.ParityNone
	} else {
		parity = serial.ParityEven
	}
	c := &serial.Config{
		Name:        deviceName,       // ou "/dev/ttyUSB0"
		Baud:        cfg.Balanca.Baud, // verifique no manual da balança; geralmente 9600 bps
		ReadTimeout: time.Second * 2,
		Parity:      parity,
	}
	s, err := serial.OpenPort(c)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	// Envia ENQ
	enq := []byte{0x05} // 05H
	_, err = port.Write(enq)
	if err != nil {
		log.Fatal(err)
	}

	// Espera a resposta da balança
	time.Sleep(100 * time.Millisecond) // ajuste se necessário

	// Lê a resposta
	buf := make([]byte, 128)
	n, err := port.Read(buf)
	if err != nil {
		log.Fatal(err)
	}

}
