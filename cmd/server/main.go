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
	http.HandleFunc("/balanca/consultaPeso", consultaPesoHandler)
	http.HandleFunc("/balanca/consultaConfig", consultaConfigHandler)
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

	cfg.Balanca = &b

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Configurações da balança atualizadas com sucesso"))
}

func consultaConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(cfg); err != nil {
		http.Error(w, "Erro ao converter JSON", http.StatusInternalServerError)
	}
}

func consultaPesoHandler(w http.ResponseWriter, r *http.Request) {
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

	peso, err := obterPesoDaBalanca()
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao consultar balança: %v", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader((http.StatusOK))
	w.Write([]byte(*peso))
}

func obterPortaSerialBalanca() string {
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
	return deviceName
}

func configurarPortaSerial(deviceName string) *serial.Config {

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
	return c
}

func obterPesoDaBalanca() (*string, error) {
	deviceName := obterPortaSerialBalanca()
	c := configurarPortaSerial(deviceName)
	serialPort, err := serial.OpenPort(c)
	if err != nil {
		return nil, fmt.Errorf("erro ao abrir porta serial: %w", err)
	}
	defer serialPort.Close()
	var resp string
	// Envia ENQ
	for {
		enq := []byte{0x05}
		_, err = serialPort.Write(enq)
		if err != nil {
			return nil, fmt.Errorf("erro ao enviar ENQ: %w", err)
		}
		buf := make([]byte, 128)
		var n int
		var resp string

		n, err = serialPort.Read(buf)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler balança: %w", err)
		}
		if n > 0 {
			resp = string(buf[:n])
			byte2 := resp[2]
			if byte2 == 'I' || byte2 == 'C' {
				continue
			}
			break
		}
	}
	switch cfg.Balanca.Protocolo {
	case "P05A":
		return tratarRespostaP05AeP06(resp), nil
	case "P05B":
		return tratarRespostaP05BeP07(resp, 4), nil
	case "P06":
		return tratarRespostaP05AeP06(resp), nil
	case "P07":
		return tratarRespostaP05BeP07(resp, 5), nil
	default:
		return nil, fmt.Errorf("protocolo não suportado: %s", cfg.Balanca.Protocolo)
	}
}

func tratarRespostaP05AeP06(resp string) *string {
	var peso string
	if len(resp) > 2 {
		peso = resp[1 : len(resp)-1]
	}
	return &peso
}
func tratarRespostaP05BeP07(resp string, pontuacao int) *string {
	var peso string
	if len(resp) > 2 {
		peso = resp[1:pontuacao] + resp[pontuacao+1:len(peso)-1]
	}
	return &peso
}
