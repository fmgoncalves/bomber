package main

import (
	"github.com/fmgoncalves/bomber/encoding/mail/sample"
	"fmt"
	"io"
	"io/ioutil"
	"net/smtp"
	"flag"
	"time"
	"math/rand"
	"strings"
	"runtime"
	"os"
)

const (
	default_host = "localhost"
	default_dir = "samples"
	default_port = 25
)

func die(err error) {
	if err != nil {
		panic(fmt.Sprintf("%v", err))
	}
	return
}

func send_default_header(key string, value string, sample_keys map[string] struct{}, wc io.WriteCloser){
	if _, contains := sample_keys[strings.ToLower(key)] ; !contains {
		wc.Write([]byte(key + ":" + value + "\r\n"))
	}
}

func hasCategory(s sample.MailSample, category string) bool {
	c := strings.ToLower(category)
	for _, s_category := range s.Type {
		if c == strings.ToLower(s_category) {
			return true
		}
	}
	return false
}

func send_sample(host string, s sample.MailSample) error {
	var err error
	// Connect to the remote SMTP server.
	var port int
	if s.Port > 0 {
		port = s.Port
	} else {
		port = default_port
	}
	c, err := smtp.Dial(fmt.Sprintf("%s:%d",host, port))
	if err != nil { return err }
	//die(err)
	defer c.Quit()
	// MAIL_FROM, TO
	c.Mail(s.From)
	c.Rcpt(s.To)
	// DATA
	wc, err := c.Data()
	if err != nil { return err }
	//die(err)
	defer wc.Close()
	// Default headers
	keyz := s.DefinedHeaders()
	send_default_header("From", s.From, keyz, wc)
	send_default_header("To", s.To, keyz, wc)
	send_default_header("Date", time.Now().UTC().Format(time.RFC822), keyz, wc)
	send_default_header("Message-Id", fmt.Sprintf("<%v@crap>",rand.Int()), keyz, wc)

	for _, h := range s.Headers {
		wc.Write([]byte(h + "\r\n"))
	}
	wc.Write([]byte("\r\n"))
	wc.Write([]byte(s.Body))
	return err
}

func main() {
	var (
		SAMPLES_DIRECTORY string
		HOST string
		CATEGORY string
		N_MESSAGES int
		THROTTLING float64
		LIST_ONLY bool
		VERBOSE bool
	)

	flag.StringVar(&HOST, "s", default_host, "SMTP server to send the message to")
	//flag.StringVar(&SAMPLES_DIRECTORY, "samples", default_dir, "Directory containing email message samples in JSON")
	samples_dir := os.Getenv("BOMBER_SAMPLES")
	if len(samples_dir) == 0 {
		samples_dir = "samples"
	}
	flag.StringVar(&SAMPLES_DIRECTORY, "samples", samples_dir, "Directory containing email message samples in JSON")
	flag.StringVar(&CATEGORY, "c", "all", "Category of messages to send")
	flag.IntVar(&N_MESSAGES, "n", 100, "Number of message to be sent")
	flag.Float64Var(&THROTTLING, "throttle", 0, "Throttle the message flow (msg/second).")
	flag.BoolVar(&LIST_ONLY, "l", false, "Only list available categories.")
	flag.BoolVar(&VERBOSE, "v", false, "Prints details of execution.")
	flag.Parse()

	l, err := ioutil.ReadDir(SAMPLES_DIRECTORY)
	die(err)
	sample_list := make([]sample.MailSample,len(l))
	categories := make(map[string]bool)
	idx := 0
	for _, val := range l {
		if ! val.IsDir() && strings.HasSuffix(val.Name(),".json") {
			s, err := sample.Unmarshal(SAMPLES_DIRECTORY+"/"+val.Name())
			if err != nil {
				if VERBOSE {
					fmt.Printf("Failed to load %s sample\n",val.Name())
				}
			}
			if VERBOSE {
				fmt.Printf("Loaded %s sample\n",val.Name())
			}
			if LIST_ONLY {
				for _, s_category := range s.Type {
					asdfg := strings.ToLower(s_category)
					categories[asdfg] = true
				}
			}
			if CATEGORY == "all" || hasCategory(s,CATEGORY) {
				sample_list[idx] = s
				idx++
			}
		}
	}
	if LIST_ONLY {
		fmt.Printf("Categories found in samples:\n")
		for cat, _ := range categories {
			fmt.Printf("\t%v\n",cat)
		}
		return
	}

	if ( idx > 0 ){
		sample_list = sample_list[0:idx]
	} else {
		panic("No samples to use.")
	}

	pn := runtime.NumCPU()
	//runtime.GOMAXPROCS(pn)
	if THROTTLING > 0 {
		pn = int(THROTTLING)
	}
	if VERBOSE {
		fmt.Printf("Sending %v messages in batches of %v\n", N_MESSAGES, pn);
	}

	c := make(chan int, pn)
	for i:= 0; i < pn; i++ {
		c <- 1
	}
	for i:= 0; i < N_MESSAGES; i++ {
		<-c
		go func (j int) {
			if VERBOSE {
				fmt.Printf("Sending message #%v\n",j+1)
			}
			err := send_sample(HOST, sample_list[rand.Int() % len(sample_list)])
			if err != nil {
					fmt.Printf("Failed to send message #%v: %v\n", j+1, err)
			}
			if THROTTLING > 0 {
				// TODO throttling should consider the running time of send_sample
				time.Sleep(time.Duration( (float64(pn) /THROTTLING) * 1000.0 * float64(time.Millisecond)))
			}
			c <- 1
		}(i)
	}
	for i:= 0; i < pn; i++ {
		<- c
	}

	return
}
