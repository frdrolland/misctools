package elk

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/google/uuid"
	//	accesslog "github.com/mash/go-accesslog"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
)

var (
	httpAddr   string
	wg         sync.WaitGroup
	running    bool = false
	outputFile *os.File
)

/*
type logger struct {
}

func (l logger) Log(record accesslog.LogRecord) {
	log.Println(record.Method + " " + record.Uri)
}
*/

/**
 * All requests.
 */
func allHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Server", "elk/1.0.0")
	w.Header().Set("Content-Type", "application/json")

	// Print request body
	bodyBuffer, _ := ioutil.ReadAll(r.Body)
	if nil != bodyBuffer && "" != string(bodyBuffer) {
		fmt.Printf("%s\n", bodyBuffer)
	}

	w.Write(prettyprint([]byte("{\"name\":\"tJHIKoQ\",\"cluster_name\":\"go-cluster\",\"cluster_uuid\":\"wRI0fPD1R1iBwKsLN4J5Ww\",\"version\":{\"number\":\"5.6.3\",\"build_hash\":\"1a2f265\",\"build_date\":\"2017-10-06T20:33:39.012Z\",\"build_snapshot\":false,\"lucene_version\":\"6.6.1\"},\"tagline\":\"You Know, for Search\"}\n")))
}

/**
 * Healthcheck implementation.
 */
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

/**
 * Healthcheck implementation.
 */
func bulkHandler(w http.ResponseWriter, req *http.Request) {
	// Print request body
	defer req.Body.Close()
	start := time.Now()

	var counter uint64 = 0

	reader := bufio.NewReader(req.Body)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)
	items := bytes.NewBufferString("")

	for scanner.Scan() {

		// First we have header
		var header *BulkHeader
		var message *BulkMessage
		err := json.Unmarshal(scanner.Bytes(), &header)
		if nil != err {
			log.Println("Error while decoding JSON from bulk header", err)
		} else {
			if scanner.Scan() {
				bytes := scanner.Bytes()
				outputFile.Write(bytes)
				outputFile.WriteString("\n")
				err := json.Unmarshal(bytes, &message)
				if nil != err {
					log.Println("Error while decoding JSON from bulk message", err)
					log.Printf("JSON was: %s\n", scanner.Text())

					// try to print message (but can be partial)
					if 1 < counter {
						items.WriteString(",")
					}
					items.WriteString(fmt.Sprintf("{\"index\":{\"_index\":\"%s\",\"_type\":\"%s\",\"_id\":\"%s\",\"_version\":1,\"result\":\"created\",\"_shards\":{\"total\":2,\"successful\":0,\"failed\":1},\"created\":true,\"status\":201}}", header.Index.Index, header.Index.Type, uuid.New().String()))
					log.Printf("[index:%s|type:%s] timestamp:%s msgtype:%s\n", header.Index.Index, header.Index.Type, message.Timestamp, message.C.Msg)
				} else {
					counter++
					if 1 < counter {
						items.WriteString(",")
					}
					items.WriteString(fmt.Sprintf("{\"index\":{\"_index\":\"%s\",\"_type\":\"%s\",\"_id\":\"%s\",\"_version\":1,\"result\":\"created\",\"_shards\":{\"total\":2,\"successful\":1,\"failed\":0},\"created\":true,\"status\":201}}", header.Index.Index, header.Index.Type, uuid.New().String()))
					log.Printf("[index:%s|type:%s] timestamp:%s msgtype:%s\n", header.Index.Index, header.Index.Type, message.Timestamp, message.C.Msg)
				}
			}
		}
	}
	w.Header().Set("Content-Type", "application/json")

	elapsed := time.Since(start)

	// Write response body
	w.Write([]byte(fmt.Sprintf("{\"took\":%d,\"errors\":false,\"items\":[", int64(elapsed/time.Millisecond))))
	w.Write(items.Bytes())
	w.Write([]byte("]}"))

	log.Printf("[elk] bulk_received_count: %d\n", counter)

	outputFile.Sync()
}

/**
 * Healthcheck implementation.
 */
func shutdownHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("shutting down http server"))
	Shutdown()
}

// formatRequest generates ascii representation of a request
func formatRequest(r *http.Request) string {
	//	// Create return string
	//	var request []string
	//	// Add the request string
	//	url := fmt.Sprintf("%v %v %v", r.Method, r.URL, r.Proto)
	//	request = append(request, url)
	//	// Add the host
	//	request = append(request, fmt.Sprintf("Host: %v", r.Host))
	//	// Loop through headers
	//	for name, headers := range r.Header {
	//		name = strings.ToLower(name)
	//		for _, h := range headers {
	//			request = append(request, fmt.Sprintf("%v: %v", name, h))
	//		}
	//	}
	//	return request
	// Buffer the body
	// Put it back before you call client.Do()
	//req.Body = myReader{bytes.NewBuffer(buf)}
	return ""
}

func Shutdown() {
	if running {
		log.Println("elasticsearch server is shutting down...")
		wg.Done()
		running = false
		os.Exit(0)
	} else {
		log.Println("elasticsearch server was already stopped...")
	}
}

//dont do this, see above edit
func prettyprint(b []byte) []byte {
	var out bytes.Buffer
	json.Indent(&out, b, "", "  ")
	return out.Bytes()
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

/**
 * Startup "elasticsearch" server to listen on http port
 */
func Startup() {

	log.Println("[elk] config file : ", viper.ConfigFileUsed())

	httpAddr = viper.GetString("elasticsearch.listen")
	outputFileName := viper.GetString("elasticsearch.output")

	// Catch signals to cleanup before exiting
	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM|syscall.SIGKILL)
	go func() {
		<-c
		Shutdown()
	}()

	running = true

	mux := http.NewServeMux()
	mux.HandleFunc("/healthcheck", healthCheckHandler)
	mux.HandleFunc("/_bulk", bulkHandler)
	mux.HandleFunc("/_shutdown", shutdownHandler)
	mux.HandleFunc("/", allHandler)
	logger := negroni.NewLogger()
	logger.SetFormat("{{.StartTime}} | {{.Status}} | \t {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}}")
	n := negroni.New(negroni.NewRecovery(), logger, negroni.NewStatic(http.Dir("public")))
	n.UseHandler(mux)

	if 0 < len(outputFileName) {
		var err error
		log.Println("[elk] creating file : ", outputFileName)
		outputFile, err = os.Create(outputFileName)
		check(err)
	}

	log.Println("[elk] starting HTTP server on address : ", httpAddr)
	// Start HTTP server in another thread
	//http.Handle("/static/", http.FileServer(http.Dir("./static")))
	//go log.Fatal(http.ListenAndServe(httpAddr, accesslog.NewLoggingHandler(handler, logger)))
	go func() {
		log.Fatal(http.ListenAndServe(httpAddr, n))
	}()
	log.Println("[elk] HTTP server started")

	// Wait for termination
	wg.Wait()
}
