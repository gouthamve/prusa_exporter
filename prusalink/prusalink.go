package prusalink

import (
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/icholy/digest"
	"github.com/pstrobl96/prusa_exporter/config"
	"github.com/rs/zerolog/log"
)

var (
	printerBoards = map[string]string{
		"MINI":    PrinterBoardTypeBuddy,
		"MK35":    PrinterBoardTypeBuddy,
		"MK39":    PrinterBoardTypeBuddy,
		"MK4":     PrinterBoardTypeBuddy,
		"XL":      PrinterBoardTypeBuddy,
		"IX":      PrinterBoardTypeBuddy,
		"I3MK3S":  PrinterBoardTypeEinsy,
		"I3MK3":   PrinterBoardTypeEinsy,
		"I3MK25S": PrinterBoardTypeEinsy,
		"I3MK25":  PrinterBoardTypeEinsy,
		"SL1":     PrinterBoardTypeSL,
		"SL1S":    PrinterBoardTypeSL,
	}

	// used for autodetection - does not work with changed hostname :sad:
	printerTypes = map[string]string{
		"PrusaMINI":         "MINI",
		"PrusaMK4":          "MK4", // unfortunately MK3.5 is also detected as MK4
		"PrusaXL":           "XL",
		"PrusaLink I3MK3S":  "I3MK3S",
		"PrusaLink I3MK3":   "I3MK3",
		"PrusaLink I3MK25S": "I3MK25S",
		"PrusaLink I3MK25":  "I3MK25",
		"prusa-sl1":         "SL1",
		"prusa-sl1s":        "SL1S",
		"Prusa_iX":          "IX", // can be found in src/common/config.h in firmware source code
	}

	configuration config.Config
)

// BoolToFloat is used for basic parsing boolean to float64
// 0.0 for false, 1.0 for true
func BoolToFloat(boolean bool) float64 {
	if !boolean {
		return 0.0
	}

	return 1.0
}

// getStateFlag returns the state flag for the given printer.
// The state flag is a float64 value representing the current state of the printer.
// It is used for tracking the printer's status and progress.
func getStateFlag(printer PrinterJSON) float64 {
	if printer.State.Flags.Operational {
		return 1
	} else if printer.State.Flags.Prepared {
		return 2
	} else if printer.State.Flags.Paused {
		return 3
	} else if printer.State.Flags.Printing {
		return 4
	} else if printer.State.Flags.Cancelling {
		return 5
	} else if printer.State.Flags.Pausing {
		return 6
	} else if printer.State.Flags.Error {
		return 7
	} else if printer.State.Flags.SdReady {
		return 8
	} else if printer.State.Flags.ClosedOrError || printer.State.Flags.ClosedOnError {
		return 9
	} else if printer.State.Flags.Ready {
		return 10
	} else if printer.State.Flags.Busy {
		return 11
	} else if printer.State.Flags.Finished {
		return 12
	} else {
		return 0
	}
}

// accessPrinterEndpoint is used to access the printer's API endpoint
func accessPrinterEndpoint(path string, printer config.Printers) ([]byte, error) {
	url := string("http://" + printer.Address + "/api/" + path)
	var (
		res    *http.Response
		result []byte
		err    error
	)

	if printer.Apikey == "" {
		client := &http.Client{
			Transport: &digest.Transport{
				Username: printer.Username,
				Password: printer.Password,
			},
			Timeout: time.Duration(configuration.Exporter.ScrapeTimeout) * time.Millisecond,
		}
		res, err = client.Get(url)

		if err != nil {
			return result, err
		}
	} else {
		req, err := http.NewRequest("GET", url, nil)
		client := &http.Client{
			Timeout: time.Duration(configuration.Exporter.ScrapeTimeout) * time.Millisecond,
		}

		if err != nil {
			return result, err
		}

		req.Header.Add("X-Api-Key", printer.Apikey)
		res, err = client.Do(req)
		if err != nil {
			return result, err
		}
	}
	result, err = io.ReadAll(res.Body)
	res.Body.Close()

	if err != nil {
		log.Error().Msg(err.Error())
	}

	return result, nil
}

// GetJobV1 is used to get the printer's job v1 API endpoint
func GetJobV1(printer config.Printers) (JobV1JSON, error) {
	var job JobV1JSON
	response, err := accessPrinterEndpoint("v1/job", printer)

	if err != nil {
		return job, err
	}

	err = json.Unmarshal(response, &job)

	return job, err
}

// GetStorageV1 is used to get the printer's storage v1 API endpoint
func GetStorageV1(printer config.Printers) (StorageV1JSON, error) {
	var storage StorageV1JSON
	response, err := accessPrinterEndpoint("v1/storage", printer)

	if err != nil {
		return storage, err
	}

	err = json.Unmarshal(response, &storage)

	return storage, err
}

// GetPrinterProfiles is used to get the printer's printerprofiles API endpoint
func GetPrinterProfiles(printer config.Printers) (PrinterProfilesJSON, error) {
	var profiles PrinterProfilesJSON
	response, err := accessPrinterEndpoint("v1/printerprofiles", printer)

	if err != nil {
		return profiles, err
	}

	err = json.Unmarshal(response, &profiles)

	return profiles, err
}

// ProbePrinter is used to probe the printer - just testing the connection
func ProbePrinter(printer config.Printers) (bool, error) {
	req, _ := http.NewRequest("GET", "http://"+printer.Address+"/", nil)
	client := &http.Client{Timeout: time.Duration(configuration.Exporter.ScrapeTimeout) * time.Millisecond}
	r, e := client.Do(req)

	if e != nil {
		return false, e
	}

	if r.StatusCode == 401 {
		log.Debug().Msg("401 Unauthorized, trying to access with API key - " + printer.Address)
		req, _ := http.NewRequest("GET", "http://"+printer.Address+"/api/v1/status", nil)
		req.Header.Add("X-Api-Key", printer.Apikey)
		r, e = client.Do(req)
		if e != nil {
			return false, e
		}
	}

	return r.StatusCode == 200, nil
}
