package analyze

import (
	"fmt"
	"time"
)


func Resumegrades(report SSLLabsResponse) string { 

	grades := ""
	for _, endpoint := range report.Endpoints {
		grades += endpoint.Grade + " "
	}
	return grades
}

func fromunixtoutc(unixTime int64) string {
    if unixTime == 0 {
        return "N/A"
    }
    tm := time.Unix(unixTime/1000, 0).UTC()
    return tm.Format("2006-01-02 15:04:05 UTC")
}

func Viewfullreport(report SSLLabsResponse) string {
    s := "==========================================\n"
    s += "           SSL LABS FULL REPORT           \n"
    s += "==========================================\n"
    s += fmt.Sprintf("Host:            %s\n", report.Host)
    s += fmt.Sprintf("Port/Protocol:   %d/%s\n", report.Port, report.Protocol)
    s += fmt.Sprintf("Status:          %s\n", report.Status)
    s += fmt.Sprintf("Engine Version:  %s (Criteria: %s)\n", report.EngineVersion, report.CriteriaVersion)
    s += fmt.Sprintf("Start Time:      %s\n", fromunixtoutc(report.StartTime))
    s += fmt.Sprintf("Test Time:       %s\n", fromunixtoutc(report.TestTime))
    s += "------------------------------------------\n"
    s += "Endpoints:\n"

    for i, ep := range report.Endpoints {
        s += fmt.Sprintf(" [%d] IP Address: %s\n", i+1, ep.IpAddress)
        s += fmt.Sprintf("     Grade:      %s (Trust Ignored: %s)\n", ep.Grade, ep.GradeTrustIgnored)
        s += fmt.Sprintf("     Message:    %s\n", ep.StatusMessage)
        s += fmt.Sprintf("     Progress:   %d%% (Duration: %dms)\n", ep.Progress, ep.Duration)
        s += fmt.Sprintf("     Warnings:   %t | Exceptional: %t\n", ep.HasWarnings, ep.IsExceptional)
    }
    s += "==========================================\n"
    return s
}