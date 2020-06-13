package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/gorilla/mux"
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
			Return  Cep `xml:"return" json:"return"`
		} `xml:"consultaCEPResponse" json:"consultaCEPResponse"`
	}
}

type Cep struct {
	Bairro       string `xml:"bairro" json:"bairro"`
	Cep          string `xml:"cep" json:"cep"`
	Cidade       string `xml:"cidade" json:"cidade"`
	Complemento2 string `xml:"complemento2" json:"complemento2"`
	End          string `xml:"end" json:"end"`
	Uf           string `xml:"uf" json:"uf"`
}

// ConsultaCEP godoc
// @Summary Consulta Cep
// @Description Consulta Cep
// @Tags cep
// @Accept json
// @Produce json
// @Param cep path string true "CEP"
// @Success 200 {object} Cep
// @Router /consultaCEP/{cep} [get]
func ConsultaCEP(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	cep := params["cep"]

	if len(cep) < 1 {
		log.Println("Url Param 'cep' is missing")
		return
	}

	// wsdl service url
	url := "https://apps.correios.com.br/SigepMasterJPA/AtendeClienteService/AtendeCliente"

	// payload
	payload := []byte(strings.TrimSpace(fmt.Sprintf(`
      <soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" xmlns:ser="http://cliente.bean.master.sigep.bsb.correios.com.br/">
         <soapenv:Body>
            <ser:consultaCEP>
               <cep>%s</cep>
            </ser:consultaCEP>
         </soapenv:Body>
      </soapenv:Envelope>`, cep),
	))

	httpMethod := "POST"
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
	w.Header().Set("Content-Type", "application/json")
	w.Write(j)

	log.Println("-> Everything is good, printing users data")
}
