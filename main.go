package main

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type ConsultaCEPResponse struct {
	XMLName xml.Name
	Body    struct {
		XMLName             xml.Name
		ConsultaCEPResponse struct {
			XMLName xml.Name
			Return  struct {
				Bairro       string `xml:"bairro" json:"bairro"`
				Cep          string `xml:"cep" json:"cep"`
				Cidade       string `xml:"cidade" json:"cidade"`
				Complemento2 string `xml:"complemento2" json:"complemento2"`
				End          string `xml:"end" json:"end"`
				Uf           string `xml:"uf" json:"uf"`
			} `xml:"return" json:"return"`
		} `xml:"consultaCEPResponse" json:"consultaCEPResponse"`
	}
}

// Baseado: https://medium.com/eaciit-engineering/soap-wsdl-request-in-go-language-3861cfb5949e
// https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente?wsdl
func main() {
	// wsdl service url
	url := "https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente"

	// payload
	payload := []byte(strings.TrimSpace(`
      <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://cliente.bean.master.sigep.bsb.correios.com.br/">
         <soapenv:Body>
            <ser:consultaCEP>
               <cep>24456422</cep>
            </ser:consultaCEP>
         </soapenv:Body>
      </soapenv:Envelope>`,
	))

	httpMethod := "POST"

	// soap action
	//soapAction := "urn:consultaCEP"

	log.Println("-> Preparing the request")

	// prepare the request
	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return
	}

	// set the content type header, as well as the other required headers
	req.Header.Set("Content-Type", "text/xml; charset=utf-8")
	req.Header.Set("Accept", "text/xml, multipart/related")
	//req.Header.Set("SOAPAction", soapAction)

	// prepare the client request
	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	log.Println("-> Dispatching the request")

	// dispatch the request
	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return
	}
	defer res.Body.Close()

	log.Println("-> Retrieving and parsing the response")

	// read and parse the response body
	result := &ConsultaCEPResponse{}
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		fmt.Println(err)
	}

	data := &ConsultaCEPResponse{}
	b, _ := xml.Marshal(result)
	_ = xml.Unmarshal(b, data)
	j, _ := json.Marshal(data.Body.ConsultaCEPResponse.Return)
	fmt.Println(string(j))

	log.Println("-> Everything is good, printing users data")
}
