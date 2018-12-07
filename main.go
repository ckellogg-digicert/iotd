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
	"time"
)

type bingImage struct {
	URL     string `json:"url"`
	URLBase string `json:"urlbase"`
}

type bingImages struct {
	Images []bingImage `json:"images"`
}

func main() {
	iotdURL := getIotdURL()
	usr, _ := user.Current()
	homeDir := usr.HomeDir

	iotdDir := homeDir + "/Pictures/Iotd"
	iotdFilename := fmt.Sprintf("%s/%s-%s.jpg", iotdDir, time.Now().Format("2006-01-02"), base64.StdEncoding.EncodeToString([]byte(iotdURL)))
	if _, err := os.Stat(iotdFilename); os.IsNotExist(err) {
		createIotdImage(iotdURL, iotdDir, iotdFilename)
	}

	applScript := `/usr/bin/osascript<<END
tell application "System Events"
	set picture of every desktop to POSIX file "%s"
end tell
END`

	if _, err := exec.Command("sh", "-c", fmt.Sprintf(applScript, iotdFilename)).Output(); err != nil {
		log.Fatalln("command failed: ", err)
	}

	log.Printf("Set desktop wallpaper: %s\n", iotdFilename)
}

func getIotdURL() string {
	iotd, _ := http.Get("https://www.bing.com/HPImageArchive.aspx?format=js&idx=0&n=8")
	var j bingImages

	defer iotd.Body.Close()
	b, _ := ioutil.ReadAll(iotd.Body)
	if err := json.Unmarshal(b, &j); err != nil {
		log.Fatalln(err)
	}

	return "https://www.bing.com" + j.Images[0].URL
}

func createIotdImage(iotdURL, iotdDir, iotdFilename string) {
	if _, err := os.Stat(iotdDir); os.IsNotExist(err) {
		log.Printf("Creating directory %s\n", iotdDir)
		os.MkdirAll(iotdDir, os.ModePerm)
	}

	log.Printf("Creating file %s\n", iotdFilename)

	resp, _ := http.Get(iotdURL)
	defer resp.Body.Close()

	iotdFile, err := os.Create(iotdFilename)
	if err != nil {
		log.Fatalln(err)
	}

	defer iotdFile.Close()
	io.Copy(iotdFile, resp.Body)
}
