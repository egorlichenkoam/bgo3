package qrcodegenerator

import (
	"bufio"
	"context"
	"errors"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Service struct {
	Timeout int64
}

const name = "qrCode service"
const externalServiceURL = "http://api.qrserver.com/v1/create-qr-code/"

func NewServive(timeout int64) *Service {
	log.Printf("%s : %s", name, "create")
	return &Service{
		Timeout: timeout,
	}
}

func (s *Service) Encode(stringtoencode string, filename string) (string, error) {
	log.Printf("%s : %s", name, "Encode - start")
	//TODO: контекст с таймаутом
	lng := len(stringtoencode)
	if (lng < 1) || (lng > 900) {
		return "", errors.New("string to encode has not proper size. should be from 1 to 900 symbols")
	}
	body, err := s.get(stringtoencode)
	if err != nil {
		return "", err
	}
	fp, err := s.export(body, filename)
	if err != nil {
		return "", err
	}
	log.Println(fp)
	log.Printf("%s : %s", name, "Encode - end")
	return fp, nil
}

func (s *Service) get(stringtoencode string) ([]byte, error) {
	log.Printf("%s : %s", name, "get - start")
	contentTypeName := "Content-Type"
	contentTypeValue := "application/x-www-form-urlencoded"
	values := make(url.Values)
	values.Set("data", stringtoencode)
	values.Set("format", "png")
	client := &http.Client{}
	ctxWithTimeout, _ := context.WithTimeout(context.Background(), time.Duration(s.Timeout)*time.Millisecond)
	req, err := http.NewRequestWithContext(ctxWithTimeout, http.MethodPost, externalServiceURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, err
	}
	req.Header.Set(contentTypeName, contentTypeValue)
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	defer func(rc io.ReadCloser) {
		if cerr := rc.Close(); cerr != nil {
			log.Println(cerr)
			err = cerr
		}
	}(resp.Body)
	if err != nil {
		return nil, err
	}
	log.Printf("%s : %s", name, "get - end")
	return body, nil
}

func (s *Service) export(qrcodebytes []byte, filename string) (string, error) {
	log.Printf("%s : %s", name, "export - start")
	filePath := ""
	file, err := os.Create(filename)
	filepath.Dir("/")
	if err != nil {
		return "", err
	}
	defer func(c io.Closer) {
		if cerr := c.Close(); cerr != nil {
			err = cerr
		}
	}(file)
	w := bufio.NewWriter(file)
	_, err = w.Write(qrcodebytes)
	w.Flush()
	if err != nil {
		return "", err
	}
	//TODO:сформирвоать путь до файла
	log.Printf("%s : %s", name, "export - end")
	return filePath, nil
}
