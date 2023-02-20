package usecase

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sync"

	"github.com/opentracing/opentracing-go"
)

func (u Usecase) SyncInscriptionEvents(rootSpan opentracing.Span) error {
	span, log := u.StartSpan("SyncInscriptionEvents", rootSpan)
	defer u.Tracer.FinishSpan(span, log)

	var err error

	errChan := make(chan error, 1)
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func(wg *sync.WaitGroup, errChan chan error) {
		defer wg.Done()
		err := u.syncInscriptionEvents(span)
		errChan <- err
	}(wg, errChan)

	wg.Wait()
	close(errChan)

	for e := range errChan {
		if e != nil {
			err = e
			log.Error("error when sync data", err.Error(), err)
		}
	}

	return err
}

func (u Usecase) syncInscriptionEvents(rootSpan opentracing.Span) error {
	url := os.Getenv("CUSTOM_ORD_SERVER")
	if url == "" {
		return errors.New("CUSTOM_ORD_SERVER is empty")
	}

	return nil
}

func getLatestOrdServiceBlockCount(rootSpan opentracing.Span, ordServer string) (int64, error) {
	url := fmt.Sprintf("%s//block-count", ordServer)
	fmt.Println("url", url)
	var result int64
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return result, err
	}

	client := &http.Client{}
	res, err := client.Do(req)

	if err != nil {
		return result, err
	}

	defer func(r *http.Response) {
		err := r.Body.Close()
		if err != nil {
			fmt.Println("Close body failed", err.Error())
		}
	}(res)

	fmt.Println("http.StatusOK", http.StatusOK, "res.Body", res.Body)

	if res.StatusCode != http.StatusOK {
		return result, errors.New("getLatestOrdServiceBlockCount Response status != 200")
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return result, errors.New("Read body failed")
	}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}
