package analyze

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

const baseUrl = "https://api.ssllabs.com/api/v2/analyze"

type StatusMsg int
type SuccessMsg string

type ErrMsg struct {
	Err error
}

// estructura de la resuesta esperada de la API de SSLLabs
type SSLLabsResponse struct {
	Host          string     `json:"host"`
	Status        string     `json:"status"`
	StatusMessage string     `json:"statusMessage"`
	Endpoints     []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Grade string `json:"grade"`
}

type AResponse struct {
	Typeofres string
	Report    SSLLabsResponse
}

func CheckSomeUrl(maxAge int, url, publish, stratnew, allopc string) tea.Cmd {
	maxAgeStr := fmt.Sprintf("%d", maxAge)
	fullUrl := baseUrl + "?host=" + url + "&publish=" + publish + "&startNew=" + stratnew +"&maxAge="+ maxAgeStr + "&all=" + allopc + "&ignoreMismatch=on"
	urlCheck := baseUrl + "?host=" + url + "&fromCache=on" + "&publish=" + publish + "&all=" + allopc + "&ignoreMismatch=on"
	return func() tea.Msg {
		c := &http.Client{Timeout: 10 * time.Second}
		res, err := c.Get(fullUrl)
		if err != nil {
			return ErrMsg{err}
		}
		//
		body, err := io.ReadAll(res.Body)
		res.Body.Close()
		if err != nil {
			return ErrMsg{err}
		}
		var sslResp SSLLabsResponse
		err = json.Unmarshal(body, &sslResp)
		if err != nil {
			return ErrMsg{err}
		}
		if sslResp.Status == "READY" {
			return AResponse{
				Typeofres: "fromcache",
				Report:    sslResp,
			}
		} else {
			for {
				time.Sleep(8 * time.Second)
				res, err := c.Get(urlCheck)
				if err != nil {
					return ErrMsg{err}
				}

				body, err := io.ReadAll(res.Body)
				res.Body.Close()
				if err != nil {
					return ErrMsg{err}
				}
				var sslResp SSLLabsResponse
				err = json.Unmarshal(body, &sslResp)
				if err != nil {
					return ErrMsg{err}
				}
				if sslResp.Status == "READY" {
					return SuccessMsg("the new Analysis for " + url + " is complete. Grade: " + sslResp.Endpoints[0].Grade)
				} else {
					return AResponse{
						Typeofres: "waiting for completion",
					}
				}

			}

		}

	}
}

