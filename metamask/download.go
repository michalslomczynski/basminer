package metamask

import (
	"archive/zip"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const (
	releasesUrl = "https://api.github.com/repos/metamask/metamask-extension/releases"
	browser  = "chrome"
	FileName = "metamask-app"
)

type ApiRelease struct {
	Assets []Asset `json:"assets,omitempty"`
}

type Asset struct {
	Name string `json:"name,omitempty"`
	Url string `json:"browser_download_url,omitempty"`
}

type ApiReleaseResp struct {
	releases []ApiRelease
}

func DownloadMetamask() error {
	if metamaskDownloaded() {
		fmt.Println("metamask is downloaded")
		return nil
	}
	fmt.Println("downloading metamask...")

	resp, err := http.Get(releasesUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	buff, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	var releases []ApiRelease
	err = json.Unmarshal(buff, &releases)
	if err != nil {
		log.Fatal(err)
	}

	url := findDownloadUrl(releases)
	if url == "" {
		return errors.New("could not fetch metamask download url")
	}

	err = downloadFile(fmt.Sprintf("%s.zip", FileName), url)
	if err != nil {
		return errors.New("could not download metamask")
	}

	files, err := unzip(fmt.Sprintf("%s.zip", FileName), FileName)
	if err != nil {
		log.Fatal("failed to unzip")
		return err
	}

	fmt.Println("Unzipped:\n" + strings.Join(files, "\n"))

	return nil
}

func metamaskDownloaded() bool {
	path, _ := filepath.Abs(FileName)
	if _, err := os.Stat(path); err != nil {
		return false
	}
	return true
}

func findDownloadUrl(releases []ApiRelease) string {
	for _, a := range releases[0].Assets {
		if strings.Contains(a.Name, browser) {
			return a.Url
		}
	}
	return ""
}

func downloadFile(filepath string, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	return err
}

func unzip(src string, dest string) ([]string, error) {
	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		fpath := filepath.Join(dest, f.Name)

		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {
			os.MkdirAll(fpath, os.ModePerm)
			continue
		}

		if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
			return filenames, err
		}

		outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
		if err != nil {
			return filenames, err
		}

		rc, err := f.Open()
		if err != nil {
			return filenames, err
		}

		_, err = io.Copy(outFile, rc)

		outFile.Close()
		rc.Close()

		if err != nil {
			return filenames, err
		}
	}
	return filenames, nil
}