package prusalink

import (
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
	GetPrinterType() (string, error)
	GetJob() (Job, error)
	GetPrinter() (PrinterJSON, error)
	GetFiles() (Files, error)
	GetVersion() (Version, error)
	GetStatus() (Status, error)
	GetInfo() (Info, error)
	GetSettings() (Settings, error)
	GetCameras() (Cameras, error)

	GetMetricLabels(job Job, labelValues ...string) []string
}

func NewPrinter(cfg config.Printers) (Printer, error) {
	if cfg.Type == "" {
		printerType, err := GetPrinterType(cfg)
		if err != nil {
			log.Error().Msg("Error while probing printer at " + cfg.Address + " - " + err.Error())
			return nil, err
		}
		cfg.Type = printerType
	}

	switch printerBoards[cfg.Type] {
	case PrinterBoardTypeBuddy:
		return buddyPrinter{basePrinter{cfg}}, nil
	case PrinterBoardTypeEinsy:
		return einsyPrinter{basePrinter{cfg}}, nil
	case PrinterBoardTypeSL:
		return slPrinter{basePrinter{cfg}}, nil
	default:
		return nil, ErrUnknownPrinterType
	}
}

type basePrinter struct {
	cfg config.Printers
}

func (p basePrinter) GetPrinterType() (string, error) {
	return GetPrinterType(p.cfg)
}

func (p basePrinter) GetJob() (Job, error) {
	return GetJob(p.cfg)
}

func (p basePrinter) GetPrinter() (PrinterJSON, error) {
	return GetPrinter(p.cfg)
}

func (p basePrinter) GetFiles() (Files, error) {
	return GetFiles(p.cfg)
}

func (p basePrinter) GetVersion() (Version, error) {
	return GetVersion(p.cfg)
}

func (p basePrinter) GetStatus() (Status, error) {
	return GetStatus(p.cfg)
}

func (p basePrinter) GetInfo() (Info, error) {
	return GetInfo(p.cfg)
}

func (p basePrinter) GetSettings() (Settings, error) {
	return GetSettings(p.cfg)
}

func (p basePrinter) GetCameras() (Cameras, error) {
	return GetCameras(p.cfg)
}

func (p basePrinter) GetMetricLabels(job Job, labelValues ...string) []string {
	return GetLabels(p.cfg, job, labelValues...)
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
