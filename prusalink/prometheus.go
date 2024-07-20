package prusalink

import (
	"strings"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/pstrobl96/prusa_exporter/config"
	"github.com/rs/zerolog/log"
)

// Collector is a struct of all printer metrics
type Collector struct {
	printerBedTemp            *prometheus.Desc
	printerPrintSpeed         *prometheus.Desc
	printerFiles              *prometheus.Desc
	printerPrintTime          *prometheus.Desc
	printerPrintTimeRemaining *prometheus.Desc
	printerPrintProgress      *prometheus.Desc
	printerMaterial           *prometheus.Desc
	printerUp                 *prometheus.Desc
	printerNozzleSize         *prometheus.Desc
	printerStatus             *prometheus.Desc
	printerAxis               *prometheus.Desc
	printerFlow               *prometheus.Desc
	printerInfo               *prometheus.Desc
	printerMMU                *prometheus.Desc
	printerCover              *prometheus.Desc
	printerAmbientTemp        *prometheus.Desc
	printerCPUTemp            *prometheus.Desc
	pritnerUVTemp             *prometheus.Desc
	printerBedTempTarget      *prometheus.Desc
	printerBedTempOffset      *prometheus.Desc
	printerChamberTemp        *prometheus.Desc
	printerChamberTempTarget  *prometheus.Desc
	printerChamberTempOffset  *prometheus.Desc
	printerToolTemp           *prometheus.Desc
	printerToolTempTarget     *prometheus.Desc
	printerToolTempOffset     *prometheus.Desc
	printerHeatedChamber      *prometheus.Desc
	printerPrintSpeedRatio    *prometheus.Desc
	printerLogs               *prometheus.Desc
	printerLogsDate           *prometheus.Desc
	printerFarmMode           *prometheus.Desc
	printerCameras            *prometheus.Desc
	printerFanSpeed           *prometheus.Desc
}

// NewCollector returns a new Collector for printer metrics
func NewCollector(config config.Config) *Collector {
	configuration = config
	defaultLabels := []string{"printer_address", "printer_model", "printer_name", "printer_job_name", "printer_job_path"}
	return &Collector{
		printerBedTemp:            prometheus.NewDesc("prusa_bed_temperature_celsius", "Current temp of printer bed in Celsius", defaultLabels, nil),
		printerPrintSpeed:         prometheus.NewDesc("prusa_print_speed_ratio", "Current setting of printer speed in ratio (0.0-1.0)", defaultLabels, nil),
		printerFiles:              prometheus.NewDesc("prusa_files_count", "Number of files in storage", append(defaultLabels, "printer_storage"), nil),
		printerPrintTimeRemaining: prometheus.NewDesc("prusa_printing_time_remaining_seconds", "Returns time that remains for completion of current print", defaultLabels, nil),
		printerPrintProgress:      prometheus.NewDesc("prusa_printing_progress", "Returns information about completion of current print in percents", defaultLabels, nil),
		printerMaterial:           prometheus.NewDesc("prusa_material_info", "Returns information about loaded filament. Returns 0 if there is no loaded filament", append(defaultLabels, "printer_filament"), nil),
		printerPrintTime:          prometheus.NewDesc("prusa_print_time_seconds", "Returns information about current print time.", defaultLabels, nil),
		printerUp:                 prometheus.NewDesc("prusa_up", "Return information about online printers. If printer is registered as offline then returned value is 0.", []string{"printer_address", "printer_model", "printer_name"}, nil),
		printerNozzleSize:         prometheus.NewDesc("prusa_nozzle_size_meters", "Returns information about selected nozzle size.", defaultLabels, nil),
		printerStatus:             prometheus.NewDesc("prusa_status_info", "Returns information status of printer.", append(defaultLabels, "printer_state"), nil),
		printerAxis:               prometheus.NewDesc("prusa_axis", "Returns information about position of axis.", append(defaultLabels, "printer_axis"), nil),
		printerFlow:               prometheus.NewDesc("prusa_print_flow_ratio", "Returns information about of filament flow in ratio (0.0 - 1.0).", defaultLabels, nil),
		printerInfo:               prometheus.NewDesc("prusa_info", "Returns information about printer.", append(defaultLabels, "api_version", "server_version", "version_text", "prusalink_name", "printer_location", "serial_number", "printer_hostname"), nil),
		printerMMU:                prometheus.NewDesc("prusa_mmu", "Returns information if MMU is enabled.", defaultLabels, nil),
		printerFanSpeed:           prometheus.NewDesc("prusa_fan_speed_rpm", "Returns information about speed of hotend fan in rpm.", append(defaultLabels, "fan"), nil),
		printerPrintSpeedRatio:    prometheus.NewDesc("prusa_print_speed_ratio", "Current setting of printer speed in values from 0.0 - 1.0", defaultLabels, nil),
		printerLogs:               prometheus.NewDesc("prusa_logs", "Return size of logs in Prusa Link", append(defaultLabels, "log_name"), nil),
		printerLogsDate:           prometheus.NewDesc("prusa_logs_date", "Return date of logs in Prusa Link", append(defaultLabels, "log_name"), nil),
		printerFarmMode:           prometheus.NewDesc("prusa_farm_mode", "Return if printer is set to farm mode", defaultLabels, nil),
		printerCameras:            prometheus.NewDesc("prusa_cameras_info", "Return information about cameras", append(defaultLabels, "camera_id", "camera_name", "camera_resolution"), nil),
		printerCover:              prometheus.NewDesc("prusa_cover_status", "Status of the printer - 0 = open, 1 = closed", defaultLabels, nil),
		printerAmbientTemp:        prometheus.NewDesc("prusa_ambient_temperature_celsius", "Status of the printer ambient temp", defaultLabels, nil),
		printerCPUTemp:            prometheus.NewDesc("prusa_cpu_temperature_celsius", "Status of the printer cpu temp", defaultLabels, nil),
		pritnerUVTemp:             prometheus.NewDesc("prusa_uv_temperature_celsius", "Status of the printer uv temp", defaultLabels, nil),
		printerBedTempTarget:      prometheus.NewDesc("prusa_bed_target_temperature_celsius", "Target bed temp", defaultLabels, nil),
		printerBedTempOffset:      prometheus.NewDesc("prusa_bed_offset_temperature_celsius", "Offset bed temp", defaultLabels, nil),
		printerChamberTemp:        prometheus.NewDesc("prusa_chamber_temperature_celsius", "Status of the printer chamber temp", defaultLabels, nil),
		printerChamberTempTarget:  prometheus.NewDesc("prusa_chamber_target_temperature_celsius", "Target chamber temp", defaultLabels, nil),
		printerChamberTempOffset:  prometheus.NewDesc("prusa_chamber_offset_temperature_celsius", "Offset chamber temp", defaultLabels, nil),
		printerToolTemp:           prometheus.NewDesc("prusa_tool_temperature_celsius", "Status of the printer tool temp", append(defaultLabels, "tool"), nil),
		printerToolTempTarget:     prometheus.NewDesc("prusa_tool_target_temperature_celsius", "Target tool temp", append(defaultLabels, "tool"), nil),
		printerToolTempOffset:     prometheus.NewDesc("prusa_tool_offset_temperature_celsius", "Offset tool temp", append(defaultLabels, "tool"), nil),
		printerHeatedChamber:      prometheus.NewDesc("prusa_heated_chamber_info", "Status of the printer heated chamber", defaultLabels, nil),
	}
}

// Describe implements prometheus.Collector
func (collector *Collector) Describe(ch chan<- *prometheus.Desc) {
	ch <- collector.printerBedTemp
	ch <- collector.printerPrintSpeed
	ch <- collector.printerFiles
	ch <- collector.printerPrintTime
	ch <- collector.printerPrintTimeRemaining
	ch <- collector.printerPrintProgress
	ch <- collector.printerMaterial
	ch <- collector.printerUp
	ch <- collector.printerNozzleSize
	ch <- collector.printerStatus
	ch <- collector.printerAxis
	ch <- collector.printerFlow
	ch <- collector.printerInfo
	ch <- collector.printerMMU
	ch <- collector.printerCover
	ch <- collector.printerAmbientTemp
	ch <- collector.printerCPUTemp
	ch <- collector.pritnerUVTemp
	ch <- collector.printerChamberTemp
	ch <- collector.printerToolTemp
	ch <- collector.printerHeatedChamber
	ch <- collector.printerBedTempTarget
	ch <- collector.printerBedTempOffset
	ch <- collector.printerChamberTempTarget
	ch <- collector.printerChamberTempOffset
	ch <- collector.printerToolTempTarget
	ch <- collector.printerToolTempOffset
	ch <- collector.printerCameras
	ch <- collector.printerFarmMode
	ch <- collector.printerLogsDate
	ch <- collector.printerLogs
	ch <- collector.printerFanSpeed
}

// Collect implements prometheus.Collector
func (collector *Collector) Collect(ch chan<- prometheus.Metric) {

	var wg sync.WaitGroup
	for _, cfg := range configuration.Printers {
		wg.Add(1)
		go func(cfg config.Printers) {
			defer wg.Done()

			printer, err := NewPrinter(cfg)
			if err != nil {
				log.Error().Msg("Error while creating printer at " + cfg.Address + " - " + err.Error())
				ch <- prometheus.MustNewConstMetric(collector.printerUp, prometheus.GaugeValue,
					0, cfg.Address, cfg.Type, cfg.Name)
				return
			}

			log.Debug().Msg("Printer scraping at " + printer.Address())
			printerUp := prometheus.MustNewConstMetric(collector.printerUp, prometheus.GaugeValue,
				0, printer.GetBaseLabels()...)

			job, err := printer.Job()
			if err != nil {
				log.Error().Msg("Error while scraping job endpoint at " + printer.Address() + " - " + err.Error())
				ch <- printerUp
				return
			}

			printerdata, err := printer.Printer()
			if err != nil {
				log.Error().Msg("Error while scraping printer endpoint at " + printer.Address() + " - " + err.Error())
				ch <- printerUp
				return
			}

			files, err := printer.Files()
			if err != nil {
				log.Error().Msg("Error while scraping files endpoint at " + printer.Address() + " - " + err.Error())
				ch <- printerUp
				return
			}

			version, err := printer.Version()
			if err != nil {
				log.Error().Msg("Error while scraping version endpoint at " + printer.Address() + " - " + err.Error())
				ch <- printerUp
				return
			}

			// metrics specific for both buddy and einsy
			if printerBoards[cfg.Type] == PrinterBoardTypeBuddy || printerBoards[cfg.Type] == PrinterBoardTypeEinsy {

				status, err := printer.Status()

				if err != nil {
					log.Error().Msg("Error while scraping status endpoint at " + printer.Address() + " - " + err.Error())
				}

				info, err := printer.Info()

				if err != nil {
					log.Error().Msg("Error while scraping info endpoint at " + printer.Address() + " - " + err.Error())
				}

				// only einsy related metrics
				if printerBoards[cfg.Type] == PrinterBoardTypeEinsy {
					settings, err := printer.Settings()

					if err != nil {
						log.Error().Msg("Error while scraping settings endpoint at " + printer.Address() + " - " + err.Error())
					} else {

						printerFarmMode := prometheus.MustNewConstMetric(
							collector.printerFarmMode, prometheus.GaugeValue,
							BoolToFloat(settings.Printer.FarmMode),
							printer.GetMetricLabels(job)...)

						ch <- printerFarmMode

					}

					cameras, err := printer.Cameras()

					if err != nil {
						log.Error().Msg("Error while scraping cameras endpoint at " + printer.Address() + " - " + err.Error())
					} else {

						for _, v := range cameras.CameraList {
							printerCamera := prometheus.MustNewConstMetric(
								collector.printerCameras, prometheus.GaugeValue,
								BoolToFloat(v.Connected),
								printer.GetMetricLabels(job, v.CameraID, v.Config.Name, v.Config.Resolution)...)
							ch <- printerCamera
						}
					}

					for _, v := range files.Files {
						printerFiles := prometheus.MustNewConstMetric(
							collector.printerFiles, prometheus.GaugeValue,
							float64(len(v.Children)),
							printer.GetMetricLabels(job, v.Display)...)
						ch <- printerFiles
					}

				}

				printerInfo := prometheus.MustNewConstMetric(
					collector.printerInfo, prometheus.GaugeValue,
					1,
					printer.GetMetricLabels(job, version.API, version.Server, version.Text, info.Name, info.Location, info.Serial, info.Hostname)...)

				ch <- printerInfo

				printerFanHotend := prometheus.MustNewConstMetric(collector.printerFanSpeed, prometheus.GaugeValue,
					status.Printer.FanHotend, printer.GetMetricLabels(job, "hotend")...)

				ch <- printerFanHotend

				printerFanPrint := prometheus.MustNewConstMetric(collector.printerFanSpeed, prometheus.GaugeValue,
					status.Printer.FanPrint, printer.GetMetricLabels(job, "print")...)

				ch <- printerFanPrint

				printerNozzleSize := prometheus.MustNewConstMetric(collector.printerNozzleSize, prometheus.GaugeValue,
					info.NozzleDiameter/1000, printer.GetMetricLabels(job)...)

				ch <- printerNozzleSize

				printSpeed := prometheus.MustNewConstMetric(
					collector.printerPrintSpeedRatio, prometheus.GaugeValue,
					printerdata.Telemetry.PrintSpeed/100,
					append(printer.GetBaseLabels(), job.Job.File.Name, job.Job.File.Path)...)

				ch <- printSpeed

				printTime := prometheus.MustNewConstMetric(
					collector.printerPrintTime, prometheus.GaugeValue,
					job.Progress.PrintTime,
					append(printer.GetBaseLabels(), job.Job.File.Name, job.Job.File.Path)...)

				ch <- printTime

				printTimeRemaining := prometheus.MustNewConstMetric(
					collector.printerPrintTimeRemaining, prometheus.GaugeValue,
					job.Progress.PrintTimeLeft,
					append(printer.GetBaseLabels(), job.Job.File.Name, job.Job.File.Path)...)

				ch <- printTimeRemaining

				printProgress := prometheus.MustNewConstMetric(
					collector.printerPrintProgress, prometheus.GaugeValue,
					job.Progress.Completion,
					append(printer.GetBaseLabels(), job.Job.File.Name, job.Job.File.Path)...)

				ch <- printProgress

				material := prometheus.MustNewConstMetric(
					collector.printerMaterial, prometheus.GaugeValue,
					BoolToFloat(!(strings.Contains(printerdata.Telemetry.Material, "-"))),
					append(printer.GetBaseLabels(), job.Job.File.Name, job.Job.File.Path, printerdata.Telemetry.Material)...)

				ch <- material

				printerAxisX := prometheus.MustNewConstMetric(
					collector.printerAxis, prometheus.GaugeValue,
					printerdata.Telemetry.AxisX,
					printer.GetMetricLabels(job, "x")...)

				ch <- printerAxisX

				printerAxisY := prometheus.MustNewConstMetric(
					collector.printerAxis, prometheus.GaugeValue,
					printerdata.Telemetry.AxisY,
					printer.GetMetricLabels(job, "y")...)

				ch <- printerAxisY

				printerAxisZ := prometheus.MustNewConstMetric(
					collector.printerAxis, prometheus.GaugeValue,
					printerdata.Telemetry.AxisZ,
					printer.GetMetricLabels(job, "z")...)

				ch <- printerAxisZ

				printerFlow := prometheus.MustNewConstMetric(collector.printerFlow, prometheus.GaugeValue,
					status.Printer.Flow/100, printer.GetMetricLabels(job)...)

				ch <- printerFlow

				if printerBoards[cfg.Type] == PrinterBoardTypeBuddy {
					printerMMU := prometheus.MustNewConstMetric(collector.printerMMU, prometheus.GaugeValue,
						BoolToFloat(info.Mmu), printer.GetMetricLabels(job)...)
					ch <- printerMMU
				}
			}

			// only sl related metrics
			if printerBoards[cfg.Type] == PrinterBoardTypeSL {
				printerCover := prometheus.MustNewConstMetric(collector.printerCover, prometheus.GaugeValue,
					BoolToFloat(printerdata.Telemetry.CoverClosed), printer.GetMetricLabels(job)...)

				ch <- printerCover

				printerFanBlower := prometheus.MustNewConstMetric(collector.printerFanSpeed, prometheus.GaugeValue,
					printerdata.Telemetry.FanBlower, printer.GetMetricLabels(job, "blower")...)

				ch <- printerFanBlower

				printerFanRear := prometheus.MustNewConstMetric(collector.printerFanSpeed, prometheus.GaugeValue,
					printerdata.Telemetry.FanRear, printer.GetMetricLabels(job, "rear")...)

				ch <- printerFanRear

				printerFanUV := prometheus.MustNewConstMetric(collector.printerFanSpeed, prometheus.GaugeValue,
					printerdata.Telemetry.FanUvLed, printer.GetMetricLabels(job, "uv")...)

				ch <- printerFanUV

				printerAmbientTemp := prometheus.MustNewConstMetric(collector.printerAmbientTemp, prometheus.GaugeValue,
					printerdata.Telemetry.TempAmbient, printer.GetMetricLabels(job)...)

				ch <- printerAmbientTemp

				printerCPUTemp := prometheus.MustNewConstMetric(collector.printerCPUTemp, prometheus.GaugeValue,
					printerdata.Telemetry.TempCPU, printer.GetMetricLabels(job)...)

				ch <- printerCPUTemp

				pritnerUVTemp := prometheus.MustNewConstMetric(collector.pritnerUVTemp, prometheus.GaugeValue,
					printerdata.Telemetry.TempUvLed, printer.GetMetricLabels(job)...)

				ch <- pritnerUVTemp

				printerChamberTempTarget := prometheus.MustNewConstMetric(collector.printerChamberTempTarget, prometheus.GaugeValue,
					printerdata.Temperature.Chamber.Target, printer.GetMetricLabels(job)...)

				ch <- printerChamberTempTarget

				printerChamberTempOffset := prometheus.MustNewConstMetric(collector.printerChamberTempOffset, prometheus.GaugeValue,
					printerdata.Temperature.Chamber.Offset, printer.GetMetricLabels(job)...)

				ch <- printerChamberTempOffset

				printerChamberTemp := prometheus.MustNewConstMetric(collector.printerChamberTemp, prometheus.GaugeValue,
					printerdata.Temperature.Chamber.Actual, printer.GetMetricLabels(job)...)

				ch <- printerChamberTemp
			}

			printerBedTemp := prometheus.MustNewConstMetric(collector.printerBedTemp, prometheus.GaugeValue,
				printerdata.Temperature.Bed.Actual, printer.GetMetricLabels(job)...)

			ch <- printerBedTemp

			printerBedTempTarget := prometheus.MustNewConstMetric(collector.printerBedTempTarget, prometheus.GaugeValue,
				printerdata.Temperature.Bed.Target, printer.GetMetricLabels(job)...)

			ch <- printerBedTempTarget

			printerBedTempOffset := prometheus.MustNewConstMetric(collector.printerBedTempOffset, prometheus.GaugeValue,
				printerdata.Temperature.Bed.Offset, printer.GetMetricLabels(job)...)

			ch <- printerBedTempOffset

			printerStatus := prometheus.MustNewConstMetric(
				collector.printerStatus, prometheus.GaugeValue,
				getStateFlag(printerdata),
				append(printer.GetBaseLabels(), job.Job.File.Name, job.Job.File.Path, printerdata.State.Text)...)

			ch <- printerStatus

			printerToolTempTarget := prometheus.MustNewConstMetric(collector.printerToolTempTarget, prometheus.GaugeValue,
				printerdata.Temperature.Tool0.Target, printer.GetMetricLabels(job, "0")...)

			ch <- printerToolTempTarget

			printerToolTempOffset := prometheus.MustNewConstMetric(collector.printerToolTempOffset, prometheus.GaugeValue,
				printerdata.Temperature.Tool0.Offset, printer.GetMetricLabels(job, "0")...)

			ch <- printerToolTempOffset

			printerToolTemp := prometheus.MustNewConstMetric(collector.printerToolTemp, prometheus.GaugeValue,
				printerdata.Temperature.Tool0.Actual, printer.GetMetricLabels(job, "0")...)

			ch <- printerToolTemp

			printerUp = prometheus.MustNewConstMetric(collector.printerUp, prometheus.GaugeValue,
				1, printer.GetBaseLabels()...)

			ch <- printerUp

			log.Debug().Msg("Scraping done at " + printer.Address())
		}(cfg)
	}
	wg.Wait()
}
