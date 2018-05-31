package influx

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
	"syscall"

	"github.com/google/uuid"
	"github.com/spf13/viper"
	"github.com/urfave/negroni"
)

var (
	httpAddr         string
	running          bool   = false
	influxdbVersion  string = "1.3.1"
	kapacitorVersion string = "1.3.1"
)

/*
type logger struct {
}

func (l logger) Log(record accesslog.LogRecord) {
	log.Println(record.Method + " " + record.Uri)
}
*/

/**
 * Write common headers, for any request handler
 */
func writeCommonHeaders(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Influxdb-Version", influxdbVersion)
	w.Header().Set("X-Kapacitor-Version", kapacitorVersion)
	w.Header().Set("Request-Id", uuid.New().String())
}

/**
 * All requests.
 */
func allHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Print request body
	bodyBuffer, _ := ioutil.ReadAll(r.Body)
	fmt.Printf("%s\n", bodyBuffer)

	w.Write(prettyprint([]byte("{\"name\":\"tJHIKoQ\",\"cluster_name\":\"go-cluster\",\"cluster_uuid\":\"wRI0fPD1R1iBwKsLN4J5Ww\",\"version\":{\"number\":\"5.6.3\",\"build_hash\":\"1a2f265\",\"build_date\":\"2017-10-06T20:33:39.012Z\",\"build_snapshot\":false,\"lucene_version\":\"6.6.1\"},\"tagline\":\"You Know, for Search\"}\n")))
}

/**
 * Healthcheck implementation.
 */
func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.Write([]byte("OK"))
}

/**
 * Healthcheck implementation.
 */
func pingHandler(w http.ResponseWriter, r *http.Request) {
	writeCommonHeaders(w)
	w.WriteHeader(http.StatusNoContent) //204
}

/**
 * Healthcheck implementation.
 */
func writeHandler(w http.ResponseWriter, req *http.Request) {
	// Print request body
	defer req.Body.Close()

	// Write common handlers
	writeCommonHeaders(w)

	var counter uint64 = 0

	reader := bufio.NewReader(req.Body)
	scanner := bufio.NewScanner(reader)

	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		//log.Printf("[index:%s|type:%s] timestamp:%s msgtype:%s\n", header.Index.Index, header.Index.Type, message.Timestamp, message.C.Msg)
		counter++
	}
	log.Printf("bulk_received_count:%d\n", counter)
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
		//wg.Done()
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

/**
 * Startup "elasticsearch" server to listen on http port
 */
func Startup() {

	httpAddr = viper.GetString("influxdb.listen")

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
	mux.HandleFunc("/ping", pingHandler)
	mux.HandleFunc("/write", writeHandler)
	mux.HandleFunc("/", allHandler)
	n := negroni.Classic() // Includes some default middlewares
	n.UseHandler(mux)

	log.Println("[influxdb] config file : ", viper.ConfigFileUsed())
	log.Println("[influxdb] starting HTTP server on address : ", httpAddr)
	// Start HTTP server in another thread
	//http.Handle("/static/", http.FileServer(http.Dir("./static")))
	//go log.Fatal(http.ListenAndServe(httpAddr, accesslog.NewLoggingHandler(handler, logger)))
	go func() {
		log.Fatal(http.ListenAndServe(httpAddr, n))
	}()
	log.Println("[elk] HTTP server started")

	// Wait for termination
}
