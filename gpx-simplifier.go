package main

import (
	"archive/zip"
	"bytes"
	"crypto/md5"
	"encoding/xml"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golangplus/fmt"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/vincent-petithory/dataurl"
)

type gpx struct {
	Creator     string       `xml:"creator,attr"`
	Title       string       `xml:"trk>name"`
	TrackPoints []TrackPoint `xml:"trk>trkseg>trkpt"`
}

type TrackPoint struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
}

type GpxFile struct {
	FileName string
	Contents []byte
}

func handleGpxUpload(w http.ResponseWriter, r *http.Request) {
	if r.Method == "GET" {
		crutime := time.Now().Unix()
		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("upload.gtpl")
		t.Execute(w, token)
	} else {
		r.ParseMultipartForm(32 << 20)
		file, _, err := r.FormFile("uploadfile")

		if err != nil {
			fmt.Println(err)
			return
		}
		numFiles, _ := strconv.Atoi(r.FormValue("num_files"))
		//num_points, _ := strconv.Atoi(r.FormValue("points_per_file"))
		contents := new(bytes.Buffer)

		//First split files
		io.Copy(contents, file)
		gpxFileBytes := contents.Bytes()
		saveFile(gpxFileBytes, "/tmp/test.txt")
		files, _ := splitGpxFile(gpxFileBytes, numFiles)
		//TODO: Then use gpsbabel to limit the points in each file

		//Now zip
		zipFileBytes, _ := createZipFile(files)
		dataEncodedUrl := dataurl.EncodeBytes(zipFileBytes)

		//Then present back to the user
		targetFileName, _ := uuid.NewV4()
		fmt.Fprintf(w, "<a download=\"%s\" href=\"%s\">Click me</a>", targetFileName.String(), dataEncodedUrl)

	}
}

func saveFile(contents []byte, fileName string) {
	ioutil.WriteFile(fileName, contents, 0666)
}

func splitGpxFile(contents []byte, numFiles int) ([]GpxFile, error) {
	buf := bytes.NewBuffer(contents)
	dec := xml.NewDecoder(buf)

	var g gpx
	err := dec.Decode(&g)
	if err != nil {
		fmtp.Printfln("Can't read file as a GPX :(")
		return nil, err
	}

	numTrackpoints := len(g.TrackPoints)
	trackpointsPerFile := numTrackpoints / numFiles
	//Due to rounding this may drop couple of points but shouldn't be a big problem on my tracks
	fmtp.Printfln("Num Trackpoints: %d, per file %d", numTrackpoints, trackpointsPerFile)

	var gpxFiles = []GpxFile{}
	for i := 0; i < numFiles; i++ {
		//Create a new GPX for each
		buffer := new(bytes.Buffer)
		fileName := fmt.Sprintf("%s_%d", g.Title, (i + 1))
		gpx := &gpx{Creator: "gpx-simplifier", Title: fileName}
		start := i * trackpointsPerFile
		end := start + trackpointsPerFile
		fmtp.Printfln("Start %d end %d", start, end)
		gpx.TrackPoints = g.TrackPoints[start:end]
		encoder := xml.NewEncoder(buffer)
		encoder.Encode(gpx)
		gpxFiles = append(gpxFiles, GpxFile{FileName: fileName, Contents: buffer.Bytes()})
	}

	return gpxFiles, nil

}

func createZipFile(files []GpxFile) ([]byte, error) {
	zipfile := new(bytes.Buffer)
	writer := zip.NewWriter(zipfile)

	for _, file := range files {
		file_name := fmt.Sprintf("%s.gpx", file.FileName)
		fmtp.Printfln("File: %s", file_name)
		zipFile, err := writer.Create(file_name)
		if err != nil {
			fmtp.Printfln("Error adding file %s to zip: %s", file.FileName, err)
			return nil, err
		}
		_, err = zipFile.Write([]byte(file.Contents))
		if err != nil {
			fmtp.Printfln("Error creating zip: %s", err)
			return nil, err
		}
	}
	//Can't defer as otherwise will be unwritten bytes
	err := writer.Close()
	if err != nil {
		fmtp.Printfln("Error closing zipfile: %s", err)
		return nil, err
	}

	zipFileAsByteArr := zipfile.Bytes()
	return zipFileAsByteArr, nil
}

func main() {
	fmt.Println("Starting up on :/8081 ...")
	http.HandleFunc("/upload", handleGpxUpload)

	log.Fatal(http.ListenAndServe(":8081", nil))
	fmt.Println("Done!")
}
