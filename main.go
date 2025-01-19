package main

import (
	"crypto/sha1"
	_ "embed"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"time"

	"github.com/siva1danil/SteamIconFix/appinfoparser"
	webview "github.com/webview/webview_go"
)

//go:embed "main.html"
var html []byte
var urlIcon string = "https://cdn.cloudflare.steamstatic.com/steamcommunity/public/images/apps/%d/%s.ico"
var wv webview.WebView
var dlDone bool
var dlResult string

func main() {
	wv = webview.New(true)
	defer func() {
		wv.Destroy()
		wv = nil
		cleanup()
	}()

	wv.SetTitle("Steam Icon Fixer")
	wv.SetSize(500, 400, webview.HintFixed)

	wv.Navigate("data:text/html;base64," + base64.StdEncoding.EncodeToString(html))
	wv.Bind("findSteamPath", findSteamPath)
	wv.Bind("findInstalledGames", findInstalledGames)
	wv.Bind("findAppInfo", findAppInfo)
	wv.Bind("findIcons", findIcons)
	wv.Bind("downloadIcon", downloadIcon)
	wv.Bind("downloadIconResult", downloadIconResult)

	wv.Run()
}

func cleanup() {
	switch runtime.GOOS {
	case "windows":
		for {
			dir := os.Getenv("APPDATA")
			exe, err := os.Executable()
			if err != nil || dir == "" {
				return
			}
			dir += "/" + filepath.Base(exe)
			err = os.RemoveAll(dir)
			if err != nil && err != os.ErrPermission {
				time.Sleep(100 * time.Millisecond)
				fmt.Println(err)
			} else {
				return
			}
		}
	}
}

func findSteamPath() (string, error) {
	switch runtime.GOOS {
	case "windows":
		path := `C:\Program Files (x86)\Steam`
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
		return "", fmt.Errorf("steam path not found for Windows")
	case "linux":
		if path := os.Getenv("HOME"); path != "" {
			steamPath := path + "/.steam/steam"
			if _, err := os.Stat(steamPath); err == nil {
				return steamPath, nil
			}
			steamPath = path + "/.local/share/Steam"
			if _, err := os.Stat(steamPath); err == nil {
				return steamPath, nil
			}
		}
		return "", fmt.Errorf("steam path not found for Linux")
	case "darwin":
		if path := os.Getenv("HOME"); path != "" {
			steamPath := path + "/Library/Application Support/Steam"
			if _, err := os.Stat(steamPath); err == nil {
				return steamPath, nil
			}
		}
		return "", fmt.Errorf("steam path not found for macOS")
	default:
		return "", fmt.Errorf("unsupported platform: %s", runtime.GOOS)
	}
}

func findInstalledGames(steam string) ([]int, error) {
	vdfPath := steam + "/steamapps/libraryfolders.vdf"
	vdfData, err := os.ReadFile(vdfPath)
	if err != nil {
		return nil, err
	}

	re, err := regexp.Compile(`\t\t\t"(\d+)"\t\t"(\d+)"`)
	if err != nil {
		return nil, err
	}

	matches := re.FindAllStringSubmatch(string(vdfData), -1)
	games := make([]int, 0, len(matches))
	for _, match := range matches {
		idString := match[1]
		idInt, err := strconv.Atoi(idString)
		if err != nil {
			return nil, err
		}
		games = append(games, idInt)
	}

	return games, nil
}

func findAppInfo(steam string) (*appinfoparser.AppInfo, error) {
	vdfPath := steam + "/appcache/appinfo.vdf"
	vdfReader, err := os.Open(vdfPath)
	if err != nil {
		return nil, err
	}
	info, err := appinfoparser.AppInfoFromReader(vdfReader)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func findIcons(steam string) ([]string, error) {
	icons := make([]string, 0)
	hashes := make(map[string]string)
	result := make([]string, 0)

	matches, err := filepath.Glob(fmt.Sprintf("%s/steam/games/*.ico", steam))
	if err != nil {
		return nil, err
	}
	for _, match := range matches {
		fileName := filepath.Base(match)
		iconName := fileName[:len(fileName)-4]

		file, err := os.Open(match)
		if err != nil {
			return nil, err
		}
		defer file.Close()

		hasher := sha1.New()
		if _, err := io.Copy(hasher, file); err != nil {
			return nil, err
		}
		icons = append(icons, iconName)
		hashes[iconName] = fmt.Sprintf("%x", hasher.Sum(nil))
	}

	for i := range icons {
		if hashes[icons[i]] == icons[i] {
			result = append(result, icons[i])
		}
	}

	return result, nil
}

func downloadIcon(steam string, id int, icon string) error {
	dlDone = false
	dlResult = ""
	go func() {
		url := fmt.Sprintf(urlIcon, id, icon)
		resp, err := http.Get(url)
		if err != nil {
			dlDone = true
			dlResult = err.Error()
			return
		}
		if resp.StatusCode != 200 {
			dlDone = true
			dlResult = fmt.Sprintf("status code: %d", resp.StatusCode)
			return
		}
		defer resp.Body.Close()

		file, err := os.Create(fmt.Sprintf("%s/steam/games/%s.ico", steam, icon))
		if err != nil {
			dlDone = true
			dlResult = err.Error()
			return
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		if err != nil {
			dlDone = true
			dlResult = err.Error()
			return
		}

		dlDone = true
		dlResult = ""
	}()
	return nil
}

func downloadIconResult() *string {
	if !dlDone {
		return nil
	}
	return &dlResult
}
