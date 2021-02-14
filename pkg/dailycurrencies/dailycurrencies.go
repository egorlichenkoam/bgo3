package dailycurrencies

import (
	"encoding/json"
	"encoding/xml"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type valuteJSON struct {
	Code  string  `json:"code"`
	Name  string  `json:"name"`
	Value float64 `json:"value"`
}

type valuteDTO struct {
	XMLName  string  `xml:"Valute"`
	Id       string  `xml:"ID,attr"`
	NumCode  int     `xml:"NumCode"`
	CharCode string  `xml:"CharCode"`
	Nominal  int64   `xml:Nominal`
	Name     string  `xml:Name`
	Value    float64 `xml:"Value"`
}

type valcursDTO struct {
	XMLName string      `xml:"ValCurs"`
	Date    string      `xml:"Date,attr"`
	Name    string      `xml:"name,attr"`
	Valutes []valuteDTO `xml:"Valute"`
}

type Service struct {
}

const name = "DailyCurrencies service"

func NewService() *Service {
	log.Printf("%s : %s", name, "create")
	return &Service{}
}

func (s *Service) Extract() error {
	log.Printf("%s : %s", name, "extract - start")
	body, err := get()
	if err != nil {
		log.Printf("%s : %s", name, "some trouble")
		log.Println(err)
		return err
	}
	vcs, err := unmarshalResponseBody(body)
	if err != nil {
		log.Printf("%s : %s", name, "some trouble")
		log.Println(err)
		return err
	}
	vcjs, err := fillValuteJsons(vcs)
	if err != nil {
		log.Printf("%s : %s", name, "some trouble")
		log.Println(err)
		return err
	}
	err = export(vcjs)
	if err != nil {
		log.Printf("%s : %s", name, "some trouble")
		log.Println(err)
		return err
	}
	log.Printf("%s : %s", name, "extract - end")
	return err
}

func get() ([]byte, error) {
	log.Printf("%s : %s", name, "get - start")
	//запрашиваем
	resp, err := http.Get("https://raw.githubusercontent.com/netology-code/bgo-homeworks/master/10_client/assets/daily.xml")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	//считываем тело ответа
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer func(rc io.ReadCloser) {
		if cerr := rc.Close(); cerr != nil {
			log.Println(cerr)
			err = cerr
		}
	}(resp.Body)
	//на все случаи жизни
	if err != nil {
		log.Println(err)
		return nil, err
	}
	log.Printf("%s : %s", name, "get - end")
	return body, nil
}

func unmarshalResponseBody(body []byte) (valcursDTO, error) {
	log.Printf("%s : %s", name, "unmarshalResponseBody - start")
	vcs := valcursDTO{}
	err := xml.Unmarshal(body, &vcs)
	if err != nil {
		log.Println(err)
		return valcursDTO{}, err
	}
	log.Printf("%s : %s", name, "unmarshalResponseBody - end")
	return vcs, nil
}

func fillValuteJsons(vcs valcursDTO) ([]valuteJSON, error) {
	log.Printf("%s : %s", name, "fillValuteJsons - start")
	vcjs := make([]valuteJSON, 0)
	i := 0
	for _, vc := range vcs.Valutes {
		vcj := valuteJSON{
			Code:  vc.CharCode,
			Name:  vc.Name,
			Value: vc.Value / float64(vc.Nominal),
		}
		vcjs = append(vcjs, vcj)
		i++
	}
	if len(vcjs) < 1 {
		log.Printf("Zero valuteJSON filled, but %d valuteDTOs were passed", len(vcjs))
		return nil, errors.New("Zero valuteJSON filled")
	}
	log.Printf("%s : %s (count : %d)", name, "fillValuteJsons - end", i)
	return vcjs, nil
}

func export(vcjs []valuteJSON) error {
	log.Printf("%s : %s", name, "export - start")
	file, err := os.Create("currencies.json")
	if err != nil {
		return err
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			err = cerr
		}
	}(file)
	encoder := json.NewEncoder(file)
	err = encoder.Encode(vcjs)
	if err != nil {
		log.Println(err)
		return err
	}
	log.Printf("%s : %s", name, "export - end")
	return nil
}
