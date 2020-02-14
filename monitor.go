package main

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/turnon/elastalarm/notifiers"
	"github.com/turnon/elastalarm/response"
)

type monitor struct {
	httpClient     *http.Client
	url            string
	done           chan bool
	timeoutRetried int
	*config
}

const (
	timeOut  = 3 * time.Second
	delayRun = 5 * time.Second
)

func initMonitors(host string, files []string) {
	for i, file := range files {
		if i != 0 {
			time.Sleep(delayRun)
		}

		cfg := loadConfig(file)
		if cfg.Skip {
			continue
		}

		m := newMonitor(host, cfg)
		m.run()
	}
}

func newMonitor(host string, cfg *config) *monitor {
	index := url.PathEscape(cfg.Index)
	url := strings.Join([]string{host, index, "_search?rest_total_hits_as_int=true"}, "/")
	return &monitor{
		config:         cfg,
		url:            url,
		httpClient:     &http.Client{Timeout: timeOut},
		done:           make(chan bool),
		timeoutRetried: 0,
	}
}

func (mon *monitor) run() {
	go func() {
		mon.check()

		for {
			select {
			case <-mon.ticker():
				mon.check()
			case <-mon.done:
				return
			}
		}
	}()
}

func (mon *monitor) check() {
	err := mon._check()
	if err != nil {
		mon.stopTicker()
		close(mon.done)
		log.Printf("%+v", err)
	}
}

func (mon *monitor) _check() error {
	req, err := http.NewRequest("GET", mon.url, mon.ReqBody())
	if err != nil {
		return errors.WithStack(err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := mon.httpClient.Do(req)
	if errHandler := mon.handleReqErr(err); errHandler != nil {
		return errHandler()
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.WithStack(err)
	}

	if resp.StatusCode != 200 {
		return errors.New(resp.Status + "|" + string(body))
	}

	respObj := &response.Response{}
	if err := respObj.Unmarshal(body); err != nil {
		return errors.WithStack(err)
	}

	mon.handleResp(respObj)

	return nil
}

func (mon *monitor) handleResp(respObj *response.Response) {
	var (
		found  bool
		result *response.Result
	)

	if mon.OnAggs() {
		found, result = mon.FoundOnAggs(respObj)
	} else {
		found, result = mon.Found(respObj)
	}

	if found {
		mon.notify(result)
	}
}

func (mon *monitor) notify(result *response.Result) {
	msg := notifiers.Msg{mon.Title, result}

	for _, notifier := range mon.notifiers {
		if err := notifier.Send(&msg); err != nil {
			log.Printf("%+v", err)
		}
	}
}

func (mon *monitor) handleReqErr(err error) func() error {
	if err == nil {
		if mon.timeoutRetried > 0 {
			if tickerErr := mon.backToDefaultTicker(); tickerErr != nil {
				return func() error {
					return errors.WithStack(tickerErr)
				}
			}
		}
		return nil
	}

	if e, ok := err.(net.Error); ok && e.Timeout() {
		mon.notify(&response.Result{
			Abstract: "retried (" + strconv.Itoa(mon.timeoutRetried) + "/" + strconv.Itoa(mon.TimeoutRetry) + ") " + err.Error(),
		})

		if mon.timeoutRetried >= mon.TimeoutRetry {
			return func() error {
				return errors.WithStack(err)
			}
		}

		return func() error {
			if err = mon.makeRetryTicker(); err != nil {
				return errors.WithStack(err)
			}
			return nil
		}
	}

	return func() error {
		mon.notify(&response.Result{Abstract: err.Error()})
		return errors.WithStack(err)
	}
}

func (mon *monitor) makeRetryTicker() error {
	mon.stopTicker()
	interval := strconv.Itoa(1<<uint(mon.timeoutRetried)) + "m"
	mon.timeoutRetried = mon.timeoutRetried + 1
	return mon.makeTicker(interval)
}

func (mon *monitor) backToDefaultTicker() error {
	mon.stopTicker()
	mon.timeoutRetried = mon.timeoutRetried - 1
	return mon.makeDefaultTicker()
}
