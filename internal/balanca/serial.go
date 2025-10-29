package balanca

import (
	"fmt"
	"log"
	"time"

	"github.com/joaocezar1/integracaoBalancaPdv/configs"
	"github.com/tarm/serial"
	"go.bug.st/serial/enumerator"
)

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
		if port.Product == "Toledo Virtual COM Port" || port.VID == "10954" ||
			(port.VID == "1509" && port.PID == "2206") {
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

func configurarPortaSerial(deviceName string, cfg *configs.Config) *serial.Config {

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

func ObterPesoDaBalanca(cfg *configs.Config) (*string, error) {
	deviceName := obterPortaSerialBalanca()
	c := configurarPortaSerial(deviceName, cfg)
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
