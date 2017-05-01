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
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/golangplus/fmt"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/vincent-petithory/dataurl"
)

type gpx struct {
	Creator     string       `xml:"creator,attr"`
	Title       string       `xml:"trk>name"`
	TrackPoints []trackPoint `xml:"trk>trkseg>trkpt"`
}

type trackPoint struct {
	Lat float64 `xml:"lat,attr"`
	Lon float64 `xml:"lon,attr"`
}

type gpxFile struct {
	FileName string
	Contents []byte
}

type resultModel struct {
	Success      bool
	ErrorMessage string
	URL          template.URL
	FileName     string
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
		numPoints, _ := strconv.Atoi(r.FormValue("points_per_file"))
		contents := new(bytes.Buffer)

		//First split files
		io.Copy(contents, file)
		gpxFileBytes := contents.Bytes()
		files, err := splitGpxFile(gpxFileBytes, numFiles)
		if err != nil {
			fmt.Println(err)
			message := fmt.Sprintf("Error reading GPX file, are you sure it's in the correct format?")
			result := &resultModel{Success: false, ErrorMessage: message}
			t, _ := template.ParseFiles("result.gtpl")
			t.Execute(w, result)
			return
		}

		files, err = simplifyGpx(files, numPoints)
		if err != nil {
			fmt.Println(err)
			message := fmt.Sprintf("Error simplifying GPX file")
			result := &resultModel{Success: false, ErrorMessage: message}
			t, _ := template.ParseFiles("result.gtpl")
			t.Execute(w, result)
			return
		}

		zipFileBytes, err := createZipFile(files)
		saveFile(zipFileBytes, "/tmp/output.zip")
		if err != nil {
			fmt.Println(err)
			message := fmt.Sprintf("Internal error packaging results: %s", err)
			result := &resultModel{Success: false, ErrorMessage: message}
			t, _ := template.ParseFiles("result.gtpl")
			t.Execute(w, result)
			return
		}
		dataEncodedURL := dataurl.EncodeBytes(zipFileBytes)

		//Then present back to the user
		targetFileName, _ := uuid.NewV4()
		result := &resultModel{Success: true, URL: template.URL(dataEncodedURL), FileName: targetFileName.String()}
		t, _ := template.ParseFiles("result.gtpl")
		t.Execute(w, result)

	}
}

func simplifyGpx(files []gpxFile, numPoints int) ([]gpxFile, error) {
	for x := range files {
		//Note: need to use pointer here as otherwise doesn't change the value in the slice
		file := &files[x]
		cmd := exec.Command("gpsbabel", "-i", "gpx", "-f", "-", "-x", "simplify,count="+strconv.Itoa(numPoints), "-o", "gpx", "-F", "-")
		cmd.Stdin = strings.NewReader(string(file.Contents))
		var outBuff bytes.Buffer
		var errBuff bytes.Buffer
		cmd.Stdout = &outBuff
		cmd.Stderr = &errBuff
		err := cmd.Run()
		if err != nil {
			return nil, err
		}
		file.Contents = outBuff.Bytes()
	}
	return files, nil
}

func saveFile(contents []byte, fileName string) {
	ioutil.WriteFile(fileName, contents, 0666)
}

func splitGpxFile(contents []byte, numFiles int) ([]gpxFile, error) {
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

	var gpxFiles = []gpxFile{}
	for i := 0; i < numFiles; i++ {
		//Create a new GPX for each
		buffer := new(bytes.Buffer)
		var fileName string
		if g.Title == "" {
			fileName = fmt.Sprintf("filepart_%d", (i + 1))
		} else {
			fileName = fmt.Sprintf("%s_%d", g.Title, (i + 1))
		}
		gpx := &gpx{Creator: "gpx-simplifier", Title: fileName}
		start := i * trackpointsPerFile
		end := start + trackpointsPerFile
		fmtp.Printfln("Start %d end %d", start, end)
		gpx.TrackPoints = g.TrackPoints[start:end]
		encoder := xml.NewEncoder(buffer)
		encoder.Encode(gpx)
		gpxFiles = append(gpxFiles, gpxFile{FileName: fileName, Contents: buffer.Bytes()})
	}

	return gpxFiles, nil

}

func createZipFile(files []gpxFile) ([]byte, error) {
	zipfile := new(bytes.Buffer)
	writer := zip.NewWriter(zipfile)

	for _, file := range files {
		fileName := fmt.Sprintf("%s.gpx", file.FileName)
		fmtp.Printfln("File: %s", fileName)
		zipFile, err := writer.Create(fileName)
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
