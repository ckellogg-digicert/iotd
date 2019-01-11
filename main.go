package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"time"
)

type bingImage struct {
	URL       string `json:"url"`
	URLBase   string `json:"urlbase"`
	StartDate string `json:"startdate"`
	Title     string `json:"title"`
	Copyright string `json:"copyright"`
}

type bingImages struct {
	Images []bingImage `json:"images"`
}

func main() {
	usr, _ := user.Current()
	homeDir := usr.HomeDir
	iotdDir := homeDir + "/Pictures/Iotd"
	iotdPrefix := time.Now().Format("20060102")

	iotdFilename, err := getIotdFile(iotdDir, iotdPrefix)
	if err != nil {
		iotdURL := getIotdURL(iotdPrefix)
		iotdFilename = createIotdImage(iotdURL, iotdDir, iotdPrefix)
	}

	applScript := `/usr/bin/osascript<<END
tell application "System Events"
	set picture of every desktop to POSIX file "%s"
end tell
END`

	if _, err := exec.Command("sh", "-c", fmt.Sprintf(applScript, iotdFilename)).Output(); err != nil {
		log.Fatalln("Command failed: ", err)
	}

	log.Printf("Set desktop background: %s\n", iotdFilename)
}

func getIotdFile(iotdDir, iotdPrefix string) (iotdFilename string, err error) {
	files, err := ioutil.ReadDir(iotdDir)
	if err != nil {
		return "", err
	}

	for _, file := range files {
		if strings.HasPrefix(file.Name(), iotdPrefix) {
			iotdFilename = iotdDir + "/" + file.Name()
			break
		}
	}

	if iotdFilename == "" {
		return "", fmt.Errorf("File not found")
	}

	return iotdFilename, nil
}

func getIotdURL(startDate string) string {
	iotd, _ := http.Get("https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=8")
	var j bingImages

	defer iotd.Body.Close()
	b, _ := ioutil.ReadAll(iotd.Body)
	if err := json.Unmarshal(b, &j); err != nil {
		log.Fatalln(err)
	}

	// Only change the date if we have a new image
	if j.Images[0].StartDate != startDate {
		log.Printf("Start date mismatch: expected: %s actual: %s\n", startDate, j.Images[0].StartDate)
		os.Exit(0)
	}

	return "https://www.bing.com" + j.Images[0].URL
}

func createIotdImage(iotdURL, iotdDir, iotdPrefix string) (iotdFilename string) {
	iotdFilename = fmt.Sprintf("%s/%s-%s.jpg", iotdDir, iotdPrefix, base64.StdEncoding.EncodeToString([]byte(iotdURL)))

	if _, err := os.Stat(iotdDir); os.IsNotExist(err) {
		log.Printf("Creating directory %s\n", iotdDir)
		os.MkdirAll(iotdDir, os.ModePerm)
	}

	log.Printf("Creating file %s\n", iotdFilename)

	resp, err := http.Get(iotdURL)
	if err != nil {
		log.Fatalf("Could not retrieve %s. Error %s\n", iotdURL, err)
	}

	defer resp.Body.Close()

	iotdFile, err := os.Create(iotdFilename)
	if err != nil {
		log.Fatalf("Could not create %s. Error %s\n", iotdFilename, err)
	}

	defer iotdFile.Close()
	io.Copy(iotdFile, resp.Body)

	return iotdFilename
}
