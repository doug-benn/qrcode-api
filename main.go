package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/yeqown/go-qrcode/v2"
	"github.com/yeqown/go-qrcode/writer/standard"
	//"github.com/yeqown/go-qrcode/writer/compressed"
)

func main() {
	mux := http.NewServeMux()
	// Registering our handler functions, and creating paths.
	mux.HandleFunc("POST /qrcode", qrcodeHandler)

	log.Println("Starting our simple http server.")
	log.Println("Started HTTP Listener")
	fmt.Println("To close connection CTRL+C")

	// Spinning up the server.
	err := http.ListenAndServe(":7878", mux)
	if err != nil {
		log.Fatal(err)
	}
}

type qrcReqBody struct {
	QrcData string `json:"qrcData"`
}

func qrcodeHandler(w http.ResponseWriter, r *http.Request) {
	//Parse Request Data
	var bodyData qrcReqBody

	if r.Body == nil {
		http.Error(w, "Body is required", http.StatusBadRequest)
		return
	}
	err := json.NewDecoder(r.Body).Decode(&bodyData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	log.Println(bodyData.QrcData)

	//Generate QR code
	qrc, err := qrcode.NewWith(bodyData.QrcData,
		qrcode.WithEncodingMode(qrcode.EncModeByte),
		qrcode.WithErrorCorrectionLevel(qrcode.ErrorCorrectionHighest),
	)
	if err != nil {
		panic(err)
	}

	buffer := new(bytes.Buffer)
	wr := nopCloser{Writer: buffer}
	w2 := standard.NewWithWriter(wr, standard.WithBuiltinImageEncoder(standard.PNG_FORMAT))
	if err = qrc.Save(w2); err != nil {
		panic(err)
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Disposition", "attachment; filename=qrcode.png") //"inline" == Display if possible else save - "attachment" will save it
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))
	if _, err := w.Write(buffer.Bytes()); err != nil {
		log.Println("unable to write image.")
	}

}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error { return nil }
