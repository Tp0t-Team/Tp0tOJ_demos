// Copyright 2022 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"compress/gzip"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type staticFiles map[string]string

func (m *staticFiles) String() string {
	return "List of files to serve"
}

func (m *staticFiles) Set(value string) error {
	parts := strings.Split(value, "=")
	if len(parts) == 2 {
		(*m)[parts[0]] = parts[1]
	} else {
		(*m)[path.Base(value)] = value
	}
	return nil
}

var (
	portFlag    = flag.Int("port", 1337, "Server HTTP port number")
	captureFlag = flag.String("capture", "traces.json.gz", "Capture file to load")
	staticsFlag = make(staticFiles)
)

type Trace struct {
	Pt                []byte    `json:"pt"`
	Ct                []byte    `json:"ct"`
	PowerMeasurements []float64 `json:"pm"`
}

type TraceMetadata struct {
	Id         int    `json:"Id"`
	Pt         string `json:"PT"`
	Ct         string `json:"CT"`
	NumSamples int    `json:"NumSamples"`
}

type Capture []Trace

func loadCaptureIo(src io.Reader) (Capture, error) {
	var capture Capture
	zipper, err := gzip.NewReader(src)
	if err != nil {
		return nil, fmt.Errorf("gzip NewReader failed %v", err)
	}
	decoder := json.NewDecoder(zipper)
	if err = decoder.Decode(&capture); err != nil {
		return nil, fmt.Errorf("JSON decoder failed %v", err)
	}
	return capture, nil
}

func loadCapture(filename string) (Capture, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("Error opening capture file: %v", err)
	}
	defer f.Close()
	return loadCaptureIo(f)
}

// ./elmo firmware.bin -Ntrace 50
func collection(quantity int) error {
	cmd := exec.Command("./elmo", "firmware.bin", "-Ntrace", strconv.Itoa(quantity))
	log.Println(cmd.String())
	err := cmd.Run()
	if err != nil {
		log.Println(err)
		return err
	}
	return nil
}

//python3 collect.py \
//--plaintext output/randdata.txt \
//--ciphertext output/printdata.txt \
//--tracedir output/traces/trace%05d.trc \
//--ntraces 50 \
//traces_raw.json.gz
func pyCollect(quantity int) error {
	commander := "python3 traces/collect.py  --plaintext output/randdata.txt --ciphertext output/printdata.txt --tracedir output/traces/trace%05d.trc --ntraces " + strconv.Itoa(quantity) + " traces_raw.json.gz"
	cmd := exec.Command("/bin/bash", "-c", commander)
	log.Println(cmd.String())
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}

	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		return err
	}

	//err := cmd.Run()
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	return nil
}

//python3 downsample.py \
//--input traces_raw.json.gz \
//--factor 5 \
//traces.json.gz
func pyDownSample(factor int) error {
	commander := "python3 traces/downsample.py --input traces_raw.json.gz --factor " + strconv.Itoa(factor) + " traces.json.gz"
	cmd := exec.Command("/bin/bash", "-c", commander)
	log.Println(cmd.String())
	stdout, err := cmd.StdoutPipe()
	cmd.Stderr = cmd.Stdout
	if err != nil {
		return err
	}
	if err = cmd.Start(); err != nil {
		return err
	}

	for {
		tmp := make([]byte, 1024)
		_, err := stdout.Read(tmp)
		fmt.Print(string(tmp))
		if err != nil {
			break
		}
	}
	if err = cmd.Wait(); err != nil {
		return err
	}

	//err := cmd.Run()
	//if err != nil {
	//	log.Println(err)
	//	return err
	//}
	return nil
}

func main() {
	flag.Var(&staticsFlag, "static", "List of static files to serve.")
	flag.Parse()
	e := echo.New()

	//capture, err := loadCapture(*captureFlag)
	//if err != nil {
	//	e.Logger.Fatal(err)
	//	return
	//}
	var capture Capture

	// Static files.
	e.File("/", "index.html")
	e.File("/viewer.js", "viewer.js")
	e.File("/viewer.css", "viewer.css")
	for k, v := range staticsFlag {
		e.File("/"+k, v)
	}

	//Returns list of static files.
	e.GET("/files", func(c echo.Context) error {
		var files []string
		for k := range staticsFlag {
			files = append(files, k)
		}
		return c.JSON(http.StatusOK, files)
	})

	// Returns trace data from a single capture file.
	e.GET("/data", func(c echo.Context) error {
		var metadata []TraceMetadata
		for i, t := range capture {
			metadata = append(metadata, TraceMetadata{i,
				hex.EncodeToString(t.Pt),
				hex.EncodeToString(t.Ct),
				len(t.PowerMeasurements)})
		}
		return c.JSON(http.StatusOK, metadata)
	})
	e.GET("/data/:trace", func(c echo.Context) error {
		trace, err := strconv.Atoi(c.Param("trace"))
		if err != nil || trace < 0 || trace >= len(capture) {
			return c.String(http.StatusInternalServerError, "Invalid trace")
		}
		return c.JSON(http.StatusOK, capture[trace].PowerMeasurements)
	})
	e.GET("/collection/:quantity", func(c echo.Context) error {
		quantity, err := strconv.Atoi(c.Param("quantity"))
		log.Println(quantity)
		if err != nil || (quantity != 50 && quantity != 100 && quantity != 200) {
			return c.String(http.StatusInternalServerError, "Invalid trace")
		}
		_, err = os.Stat("./elmo")
		if err != nil {
			return c.String(http.StatusInternalServerError, "No elmo")
		}
		err = collection(quantity)
		if err != nil {
			return c.String(http.StatusInternalServerError, "elmo exec fault")
		}
		err = pyCollect(quantity)
		if err != nil {
			return c.String(http.StatusInternalServerError, "pyCollect exec fault")
		}
		err = pyDownSample(5)
		if err != nil {
			return c.String(http.StatusInternalServerError, "pyDownSample exec fault")
		}
		capture, err = loadCapture(*captureFlag)
		if err != nil {
			log.Println(err)
			return c.String(http.StatusInternalServerError, "capture reload fault")
		}
		return c.String(http.StatusOK, "")

	})

	e.Logger.Fatal(e.Start(fmt.Sprintf("0.0.0.0:%d", *portFlag)))
}
