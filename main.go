package main

import (
	"flag"
	"fmt"
	"github.com/tidwall/gjson"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
)

var (
	repoFlag = flag.String("repo", "", "Repository like elsgaard/tractor")
	relFlag  = flag.String("r", "latest", "A release version (default: 'latest')")
	patFlag  = flag.String("p", "", "Personal access token (no default)")
	nameFlag = flag.String("a", "", "Artifact name (no default)")
	pathFlag = flag.String("o", "", "path/filename (no default)")
)

func main() {

	flag.Parse()

	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+*repoFlag+"/releases/"+*relFlag+"", nil)
	req.Header.Set("Authorization", "Bearer "+*patFlag+"")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	s := string(bodyText)
	result := gjson.Get(s, "assets")

	result.ForEach(func(key, value gjson.Result) bool {
		assetID := gjson.Get(value.String(), "id")
		assetName := gjson.Get(value.String(), "name")

		matched, _ := regexp.MatchString(*nameFlag, assetName.String())
		url := createDownloadUrl(assetID.String())

		if matched {
			err = downloadFile(*pathFlag, url)
			if err != nil {
				fmt.Println(err)
			}
		}

		return true // keep iterating
	})

	return

}

func downloadFile(filepath string, url string) (err error) {
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get the file
	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+*patFlag+"")
	req.Header.Set("Accept", "application/octet-stream")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	defer resp.Body.Close()

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	// Writer the body to file
	size, err := io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	fmt.Printf("Downloaded a file %s with size %d", filepath, size)
	return nil
}

func createDownloadUrl(id string) string {
	url := "https://api.github.com/repos/" + *repoFlag + "/releases/assets/" + id + ""
	return url

}
