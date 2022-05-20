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
	repoFlag   = flag.String("repo", "", "Repository like elsgaard/tractor")
	patFlag    = flag.String("p", "", "Personal access token (no default)")
	nameFlag   = flag.String("a", "", "Artifact name (no default)")
	pathFlag   = flag.String("o", "", "path/filename (no default)")
	sourceFlag = flag.String("s", "", "Source instead of release artifact: 'zip' or 'tar'(no default)")
)

func main() {

	flag.Parse()

	if *nameFlag != "" && *sourceFlag != "" {
		fmt.Println("Only artifact or source is allowed")
		return
	}

	if *pathFlag == "" {
		fmt.Println("path/Filename is mandatory")
		return
	}

	if *repoFlag == "" {
		fmt.Println("Github repository is mandatory")
		return
	}

	s := getRelease()

	if *sourceFlag != "" {
		switch *sourceFlag {
		case "zip":
			err := downloadFile(*pathFlag, gjson.Get(s, "zipball_url").String())
			if err != nil {
				return
			}
		case "tar":
			err := downloadFile(*pathFlag, gjson.Get(s, "tarball_url").String())
			if err != nil {
				return
			}
		default:
			fmt.Println("Invalid source, must be zip or tar")
		}
		// We are done, exit
		return
	}

	// We are looking for an asset
	assets := gjson.Get(s, "assets")
	assets.ForEach(func(key, value gjson.Result) bool {
		assetID := gjson.Get(value.String(), "id")
		assetName := gjson.Get(value.String(), "name")

		matched, _ := regexp.MatchString(*nameFlag, assetName.String())
		url := createDownloadUrl(assetID.String())

		if matched {
			err := downloadFile(*pathFlag, url)
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

	client := &http.Client{}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("Authorization", "Bearer "+*patFlag+"")
	req.Header.Set("Accept", "*/*")
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

// getRelease is downloading a list of the GitHub releases i JSON format
func getRelease() string {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", "https://api.github.com/repos/"+*repoFlag+"/releases/latest", nil)
	req.Header.Set("Authorization", "Bearer "+*patFlag+"")
	req.Header.Set("Accept", "application/vnd.github.v3+json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	return string(bodyText)

}
