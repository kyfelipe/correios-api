package services

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/kyfelipe/correios-api/utils"
	"log"
	"net/http"
	"strings"
)

type CorreiosResponse struct {
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
// @Param cep query string true "CEP"
// @Success 200 {object} Cep
// @Failure 400 {object} utils.HTTPError
// @Failure 500 {object} utils.HTTPError
// @Router /consultaCEP [get]
func ConsultaCEP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	cep := r.URL.Query().Get("cep")

	if len(cep) < 8 {
		w.WriteHeader(http.StatusBadRequest)
		setErrorMessage(&w, http.StatusBadRequest, "Url Param 'cep' is invalid")
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
	log.Println("Preparing the request")

	// prepare the request
	req, err := http.NewRequest(httpMethod, url, bytes.NewBuffer(payload))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		setErrorMessage(&w, http.StatusInternalServerError, fmt.Sprintf("Error on creating request object. %s", err.Error()))
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
	log.Println("Dispatching the request")

	// dispatch the request
	res, err := client.Do(req)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		setErrorMessage(&w, http.StatusInternalServerError, fmt.Sprintf("Error on dispatching request. %s", err.Error()))
		return
	}
	defer res.Body.Close()

	log.Println("Retrieving and parsing the response")

	// read and parse the response body
	result := &CorreiosResponse{}
	err = xml.NewDecoder(res.Body).Decode(result)
	if err != nil {
		fmt.Println(err)
	}

	data := &CorreiosResponse{}
	b, _ := xml.Marshal(result)
	_ = xml.Unmarshal(b, data)

	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(data.Body.ConsultaCEPResponse.Return)

	log.Println("Everything is good, response cep data")
}

func setErrorMessage(w *http.ResponseWriter, statusCode int, message string) {
	_ = json.NewEncoder(*w).Encode(utils.HTTPError{
		Code:    statusCode,
		Message: message,
	})
	log.Println(message)
}
