package main

import (
	"crypto/tls"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/yankeguo/rg"
	"gopkg.in/yaml.v3"
)

type ValidationOptions struct {
	Debug       bool `yaml:"debug"`
	DaysAdvance int  `yaml:"days_advance"`
}

type Options struct {
	Domains    []string          `yaml:"domains"`
	Validation ValidationOptions `yaml:"validation"`
	Notify     struct {
		WxworkBot string `yaml:"wxwork_bot"`
	} `yaml:"notify"`
}

func main() {
	var err error
	defer func() {
		if err == nil {
			return
		}
		log.Println("exited with error:", err.Error())
		os.Exit(1)
	}()
	defer rg.Guard(&err)

	if len(os.Args) < 2 {
		err = fmt.Errorf("usage: %s <config.yaml>", os.Args[0])
		return
	}

	var opts Options
	rg.Must0(yaml.Unmarshal(rg.Must(os.ReadFile(os.Args[1])), &opts))

	if opts.Validation.DaysAdvance <= 0 {
		opts.Validation.DaysAdvance = 7
	}

	log.Printf("options: %+v", opts)

	var results []string

	for _, domain := range opts.Domains {
		result := rg.Must(validateDomain(domain, opts.Validation))

		if result != "" {
			results = append(results, result)
		}
	}

	if len(results) == 0 {
		return
	}

	client := resty.New()

	res := rg.Must(client.R().SetBody(map[string]any{
		"msgtype": "text",
		"text": map[string]any{
			"content": strings.Join(results, "\n"),
		},
	}).Post(opts.Notify.WxworkBot))

	if res.IsError() {
		err = errors.New(res.String())
		return
	}
}

func validateDomain(domain string, opts ValidationOptions) (result string, err error) {
	defer rg.Guard(&err)

	var (
		host string
		port string
		name string
	)

	splits := strings.Split(domain, ":")

	if len(splits) == 1 {
		host = splits[0]
		port = "443"
		name = splits[0]
	} else if len(splits) == 2 {
		host = splits[0]
		port = splits[1]
		name = splits[0]
	} else if len(splits) == 3 {
		host = splits[0]
		port = splits[1]
		name = splits[2]
	} else {
		err = fmt.Errorf("invalid domain: %s", domain)
		return
	}

	var conn *tls.Conn

	if conn, err = tls.Dial("tcp", host+":"+port, &tls.Config{
		ServerName: name,
	}); err != nil {
		result = fmt.Sprintf("failed to connect to %s: %s", domain, err.Error())
		err = nil
		return
	}
	defer conn.Close()

	cert := conn.ConnectionState().PeerCertificates[0]

	if cert.NotAfter.AddDate(0, 0, -opts.DaysAdvance).Before(time.Now()) || opts.Debug {
		result = fmt.Sprintf("certificate of %s will expire at %s", domain, cert.NotAfter.Format("2006-01-02"))
	}

	return
}
