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
	Job() (JobJSON, error)
	Printer() (PrinterJSON, error)
	Files() (FilesJSON, error)
	Version() (VersionJSON, error)
	Status() (StatusJSON, error)
	Info() (InfoJSON, error)
	Settings() (SettingsJSON, error)
	Cameras() (CamerasJSON, error)

	GetBaseLabels() []string
	GetMetricLabels(job JobJSON, labelValues ...string) []string
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
		return buddyPrinter{bp}, nil
	case PrinterBoardTypeEinsy:
		return einsyPrinter{bp}, nil
	case PrinterBoardTypeSL:
		return slPrinter{bp}, nil
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

func (p basePrinter) Job() (JobJSON, error) {
	var job JobJSON
	response, err := accessPrinterEndpoint("job", p.cfg)

	if err != nil {
		return job, err
	}

	err = json.Unmarshal(response, &job)

	return job, err
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

func (p basePrinter) GetMetricLabels(job JobJSON, labelValues ...string) []string {
	if job == (JobJSON{}) {
		return append([]string{p.cfg.Address, p.cfg.Type, p.cfg.Name, "", ""}, labelValues...)
	}
	return append([]string{p.cfg.Address, p.cfg.Type, p.cfg.Name, job.Job.File.Name, job.Job.File.Path}, labelValues...)
}

type buddyPrinter struct {
	basePrinter
}

type einsyPrinter struct {
	basePrinter
}

type slPrinter struct {
	basePrinter
}
