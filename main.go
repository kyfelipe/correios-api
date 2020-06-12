package main

import (
	"correios/services"
	"net/http"
)

// Baseado: https://medium.com/eaciit-engineering/soap-wsdl-request-in-go-language-3861cfb5949e
// https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl
func main() {
	http.HandleFunc("/consultaCEP", services.ConsultaCEP)
	http.ListenAndServe(":4000", nil)
}
