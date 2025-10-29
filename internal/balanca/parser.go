package balanca

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
		peso = resp[1:pontuacao] + resp[pontuacao+1:len(resp)-1]
	}
	return &peso
}
