package prusalink

import (
	"encoding/json"
	"errors"

	"github.com/pstrobl96/prusa_exporter/config"
	"github.com/rs/zerolog/log"
)

const (
	PrinterBoardTypeBuddy = "buddy"
	PrinterBoardTypeEinsy = "einsy"
	PrinterBoardTypeSL    = "sl"
)

var (
	ErrUnknownPrinterType = errors.New("unknown printer type")
)

type Printer interface {
	Address() string

	PrinterType() (string, error)
	Job() (Job, error)
	Printer() (PrinterJSON, error)
	Files() (FilesJSON, error)
	Version() (VersionJSON, error)
	Status() (StatusJSON, error)
	Info() (InfoJSON, error)
	Settings() (SettingsJSON, error)
	Cameras() (CamerasJSON, error)

	GetBaseLabels() []string
	GetMetricLabels(job Job, labelValues ...string) []string
}

func NewPrinter(cfg config.Printers) (Printer, error) {
	bp := basePrinter{cfg}

	if cfg.Type == "" {
		printerType, err := bp.PrinterType()
		if err != nil {
			log.Error().Msg("Error while probing printer at " + cfg.Address + " - " + err.Error())
			return nil, err
		}
		bp.cfg.Type = printerType
	}

	switch printerBoards[cfg.Type] {
	case PrinterBoardTypeBuddy:
		return v1APIPrinter{bp}, nil
	case PrinterBoardTypeSL, PrinterBoardTypeEinsy:
		return bp, nil
	default:
		return nil, ErrUnknownPrinterType
	}
}

type basePrinter struct {
	cfg config.Printers
}

func (p basePrinter) Address() string {
	return p.cfg.Address
}

func (p basePrinter) PrinterType() (string, error) {
	if p.cfg.Type != "" {
		return p.cfg.Type, nil
	}

	version, err := p.Version()
	if err != nil {
		return "unknown", err
	}

	printerType := version.Hostname

	if version.Hostname == "" {
		if version.Original == "" {
			info, err := p.Info()
			if err != nil {
				return "unknown", err
			}
			printerType = info.Hostname
		} else {
			printerType = version.Original
		}
	} else if version.Original != "" {
		printerType = version.Original
	}

	if printerTypes[printerType] != "" {
		printerType = printerTypes[printerType]
	}

	if printerType == "" {
		printerType = "unknown"
	}

	log.Trace().Msg(printerType + " detected for " + p.cfg.Address + " (" + p.cfg.Name + ")")

	return printerType, nil
}

func (p basePrinter) Job() (Job, error) {
	var jobjson JobJSON
	response, err := accessPrinterEndpoint("job", p.cfg)

	if err != nil {
		return Job{}, err
	}

	err = json.Unmarshal(response, &jobjson)
	if err != nil {
		return Job{}, err
	}

	job := Job{}
	job.File.Name = jobjson.Job.File.Name
	job.File.Path = jobjson.Job.File.Path
	job.Progress.PercentDone = jobjson.Progress.Completion
	job.Progress.TimeLeft = jobjson.Progress.PrintTimeLeft
	job.Progress.TimeElapsed = jobjson.Progress.PrintTime

	return job, nil
}

func (p basePrinter) Printer() (PrinterJSON, error) {
	var printerData PrinterJSON
	response, err := accessPrinterEndpoint("printer", p.cfg)

	if err != nil {
		return printerData, err
	}

	err = json.Unmarshal(response, &printerData)

	return printerData, err
}

func (p basePrinter) Files() (FilesJSON, error) {
	var files FilesJSON
	response, err := accessPrinterEndpoint("files?recursive=true", p.cfg)

	if err != nil {
		return files, err
	}

	err = json.Unmarshal(response, &files)

	return files, err
}

func (p basePrinter) Version() (VersionJSON, error) {
	var version VersionJSON
	response, err := accessPrinterEndpoint("version", p.cfg)

	if err != nil {
		return version, err
	}

	err = json.Unmarshal(response, &version)

	return version, err
}

func (p basePrinter) Status() (StatusJSON, error) {
	var status StatusJSON
	response, err := accessPrinterEndpoint("v1/status", p.cfg)

	if err != nil {
		return status, err
	}

	err = json.Unmarshal(response, &status)

	return status, err
}

func (p basePrinter) Info() (InfoJSON, error) {
	var info InfoJSON
	response, err := accessPrinterEndpoint("v1/info", p.cfg)

	if err != nil {
		return info, err
	}

	err = json.Unmarshal(response, &info)

	return info, err
}

func (p basePrinter) Settings() (SettingsJSON, error) {
	var settings SettingsJSON
	response, err := accessPrinterEndpoint("settings", p.cfg)

	if err != nil {
		return settings, err
	}

	err = json.Unmarshal(response, &settings)

	return settings, err
}

func (p basePrinter) Cameras() (CamerasJSON, error) {
	var cameras CamerasJSON
	response, err := accessPrinterEndpoint("v1/cameras", p.cfg)

	if err != nil {
		return cameras, err
	}

	err = json.Unmarshal(response, &cameras)

	return cameras, err
}

func (p basePrinter) GetBaseLabels() []string {
	return []string{p.cfg.Address, p.cfg.Type, p.cfg.Name}
}

func (p basePrinter) GetMetricLabels(job Job, labelValues ...string) []string {
	if job == (Job{}) {
		return append([]string{p.cfg.Address, p.cfg.Type, p.cfg.Name, "", ""}, labelValues...)
	}
	return append([]string{p.cfg.Address, p.cfg.Type, p.cfg.Name, job.File.Name, job.File.Path}, labelValues...)
}

type v1APIPrinter struct {
	basePrinter
}

func (p v1APIPrinter) Job() (Job, error) {
	var jobV1 JobV1JSON
	response, err := accessPrinterEndpoint("v1/job", p.cfg)
	if err != nil {
		return Job{}, err
	}

	err = json.Unmarshal(response, &jobV1)
	if err != nil {
		return Job{}, err
	}

	job := Job{}
	job.ID = jobV1.ID
	job.Progress.PercentDone = jobV1.Progress
	job.Progress.TimeLeft = jobV1.TimeRemaining
	job.Progress.TimeElapsed = jobV1.TimePrinting

	job.File.Name = jobV1.File.DisplayName
	job.File.Path = jobV1.File.Path + "/" + jobV1.File.Name

	return job, nil
}

type Job struct {
	File struct {
		Name string
		Path string
	}
	Progress struct {
		PercentDone float64
		TimeLeft    float64
		TimeElapsed float64
	}
	ID int
}
