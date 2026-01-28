package analyze

import (
	"encoding/json"
	//"fmt"
	"io"
	"net/http"
	"time"
	tea "github.com/charmbracelet/bubbletea"
)

const baseUrl = "https://api.ssllabs.com/api/v2/analyze"

type StatusMsg int
type SuccessMsg string



type errMsg struct {
	err error
}

// estructura de la resuesta esperada de la API de SSLLabs
type SSLLabsResponse struct {
	Status        string     `json:"status"`
	StatusMessage string     `json:"statusMessage"`
	Endpoints     []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Grade string `json:"grade"`
}

func CheckSomeUrl(maxAge int, url,  publish, stratnew, allopc string) tea.Cmd {
	fullUrl := baseUrl + "?host=" + url + "&publish=" + publish + "&startNew=" + stratnew + "&all=" + allopc + "&ignoreMismatch=on"
	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(fullUrl)
		if err != nil {
			return errMsg{err}
		}

		return StatusMsg(res.StatusCode)
	}
}

func KeepChecking(url string) tea.Cmd {
	urlCheck := baseUrl + "?host=" + url + "&publish=off&all=done&ignoreMismatch=on"

	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		for {
			time.Sleep(8 * time.Second)
			res, err := c.Get(urlCheck)
			if err != nil {
				return errMsg{err}
			}

			body, err := io.ReadAll(res.Body)
			res.Body.Close()
			if err != nil {
				return errMsg{err}
			}
			var sslResp SSLLabsResponse
			err = json.Unmarshal(body, &sslResp)
			if err != nil {
				return errMsg{err}
			}
			if sslResp.Status == "READY" {
				return SuccessMsg("Analysis for " + url + " is complete. Grade: " + sslResp.Endpoints[0].Grade)
			}
		}
	}
}

// checks for cache results if thre is no a cache result starts a new analysis
func Checkfromcache(url string) tea.Cmd {
	urlCheck := baseUrl + "?host=" + url + "&fromCache=on&publish=off&all=done&ignoreMismatch=on"

	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}

		res, err := c.Get(urlCheck)
		if err != nil {
			return errMsg{err}
		}

		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return errMsg{err}
		}
		var sslResp SSLLabsResponse
		err = json.Unmarshal(body, &sslResp)
		if err != nil {
			return errMsg{err}
		}

		switch sslResp.Status {
		case "ERROR":
			return errMsg{err}
		case "IN_PROGRESS":
			return StatusMsg(202)
		}
		if sslResp.Status == "READY" {
			return SuccessMsg("Cached analysis for " + url + " found. Grade: " + sslResp.Endpoints[0].Grade)
		} else {
			// start a new analysis
			fullUrl := baseUrl + "?host=" + url + "&publish=off&startNew=on&all=done&ignoreMismatch=on"
			_, err := c.Get(fullUrl)
			if err != nil {
				return errMsg{err}
			}
			return SuccessMsg("there are no cached results. Started new analysis for " + url)
		}
	}
}