//go:build ignore

package main

import (
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

func bucket_key(r *http.Request) (string, string) {
	slashIdx := strings.IndexByte(r.URL.Path[1:], '/')
	if slashIdx == -1 {
		return r.URL.Path[1:], ""
	} else {
		return r.URL.Path[1 : 1+slashIdx], r.URL.Path[1+slashIdx:]
	}

}
func handler(w http.ResponseWriter, r *http.Request) {
	reqDump, err := httputil.DumpRequest(r, true)
	if err != nil {
		log.Fatal(err)
	}
	bucket, key := bucket_key(r)

	fmt.Printf("REQUEST:\n%s", string(reqDump))
	fmt.Printf("Bucket: %s req %s b %t", bucket, r.Method, r.Method == http.MethodPut)
	bodyBytes, err := io.ReadAll(r.Body)
	fmt.Printf("Body: %s", string(bodyBytes))
	if err != nil {
		panic(err)
	}
	if r.Method == http.MethodPut {
		w.Header().Set("Location", fmt.Sprintf("/%s", bucket))
		err := os.MkdirAll(fmt.Sprintf("./tmp/%s", bucket), fs.ModePerm)
		if err != nil {
			panic(err)
		}
		if key != "" {
			err = os.WriteFile(fmt.Sprintf("./tmp/%s/%s", bucket, key), bodyBytes, fs.ModePerm)
			if err != nil {
				panic(err)
			}
		}
		/*
		   r.ParseMultipartForm(32 << 20)
		   var buf bytes.Buffer
		   file, header, err := r.FormFile("file")
		   if err != nil {
		       panic(err)
		   }
		   defer file.Close()
		   name := strings.Split(header.Filename, ".")
		   fmt.Printf("File name %s\n", name[0])
		   io.Copy(&buf, file)
		   contents := buf.String()
		   fmt.Println(contents)
		   buf.Reset()
		*/
	} else if r.Method == http.MethodGet {
		fileBytes, err := os.ReadFile(fmt.Sprintf("./tmp/%s/%s", bucket, key))
		if err != nil {
			panic(err)
		}
		w.Write(fileBytes)
	}
}

func main() {
	err := os.Mkdir("./tmp", fs.ModePerm)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":80", nil))
}
