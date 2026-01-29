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
	Host            string     `json:"host"`
	Port            int        `json:"port"`
	Protocol        string     `json:"protocol"`
	IsPublic        bool       `json:"isPublic"`
	Status          string     `json:"status"`
	StartTime       int64      `json:"startTime"`
	TestTime        int64      `json:"testTime"`
	EngineVersion   string     `json:"engineVersion"`
	CriteriaVersion string     `json:"criteriaVersion"`
	Endpoints       []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	IpAddress         string `json:"ipAddress"`
	StatusMessage     string `json:"statusMessage"`
	Grade             string `json:"grade"`
	GradeTrustIgnored string `json:"gradeTrustIgnored"`
	HasWarnings       bool   `json:"hasWarnings"`
	IsExceptional     bool   `json:"isExceptional"`
	Progress          int    `json:"progress"`
	Duration          int    `json:"duration"`
	Eta               int    `json:"eta"`
	Delegation        int    `json:"delegation"`
}

type AResponse struct {
	Typeofres string
	Report    SSLLabsResponse
}

func CheckSomeUrl(maxAge int, url, publish, stratnew, allopc string) tea.Cmd {
	maxAgeStr := fmt.Sprintf("%d", maxAge)
	fullUrl := baseUrl + "?host=" + url + "&publish=" + publish + "&startNew=" + stratnew + "&maxAge=" + maxAgeStr + "&all=" + allopc + "&ignoreMismatch=on"
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
					return AResponse{
						Typeofres: "fromnewanalysis",
						Report:    sslResp,
					}
				}
				// El bucle continúa si Status != "READY"
			}
		}

	}
}
