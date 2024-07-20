package prusalink_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/prometheus/client_golang/prometheus/testutil"
	"github.com/pstrobl96/prusa_exporter/config"
	"github.com/pstrobl96/prusa_exporter/prusalink"
)

func TestBuddyPrinter(t *testing.T) {
	server := runPrusaLinkAPIServer("api/buddy")
	defer server.Close()

	collector := prusalink.NewCollector(config.Config{
		Printers: []config.Printers{
			{
				Address: strings.TrimPrefix(server.URL, "http://"),
				Type:    "MK4",
			},
		},
	})

	output := `
# HELP prusa_axis Returns information about position of axis.
# TYPE prusa_axis gauge
prusa_axis{printer_address="PRINTER_ADDRESS",printer_axis="x",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 0
prusa_axis{printer_address="PRINTER_ADDRESS",printer_axis="y",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 0
prusa_axis{printer_address="PRINTER_ADDRESS",printer_axis="z",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 0
# HELP prusa_bed_offset_temperature_celsius Offset bed temp
# TYPE prusa_bed_offset_temperature_celsius gauge
prusa_bed_offset_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 0
# HELP prusa_bed_target_temperature_celsius Target bed temp
# TYPE prusa_bed_target_temperature_celsius gauge
prusa_bed_target_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 60
# HELP prusa_bed_temperature_celsius Current temp of printer bed in Celsius
# TYPE prusa_bed_temperature_celsius gauge
prusa_bed_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 60.1
# HELP prusa_fan_speed_rpm Returns information about speed of hotend fan in rpm.
# TYPE prusa_fan_speed_rpm gauge
prusa_fan_speed_rpm{fan="hotend",printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 0
prusa_fan_speed_rpm{fan="print",printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 0
# HELP prusa_info Returns information about printer.
# TYPE prusa_info gauge
prusa_info{api_version="2.0.0",printer_address="PRINTER_ADDRESS",printer_hostname="PrusaMK4",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_location="",printer_model="MK4",printer_name="",prusalink_name="",serial_number="10589-3742441631728135",server_version="2.1.2",version_text="PrusaLink"} 1
# HELP prusa_job_info Returns information about current job.
# TYPE prusa_job_info gauge
prusa_job_info{printer_address="PRINTER_ADDRESS",printer_job_id="109",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 1
# HELP prusa_material_info Returns information about loaded filament. Returns 0 if there is no loaded filament
# TYPE prusa_material_info gauge
prusa_material_info{printer_address="PRINTER_ADDRESS",printer_filament="PLA",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 1
# HELP prusa_mmu Returns information if MMU is enabled.
# TYPE prusa_mmu gauge
prusa_mmu{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 0
# HELP prusa_nozzle_size_meters Returns information about selected nozzle size.
# TYPE prusa_nozzle_size_meters gauge
prusa_nozzle_size_meters{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 0.0004
# HELP prusa_print_flow_ratio Returns information about of filament flow in ratio (0.0 - 1.0).
# TYPE prusa_print_flow_ratio gauge
prusa_print_flow_ratio{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 1
# HELP prusa_print_speed_ratio Current setting of printer speed in values from 0.0 - 1.0
# TYPE prusa_print_speed_ratio gauge
prusa_print_speed_ratio{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 1
# HELP prusa_print_time_seconds Returns information about current print time.
# TYPE prusa_print_time_seconds gauge
prusa_print_time_seconds{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 254
# HELP prusa_printing_progress Returns information about completion of current print in percents
# TYPE prusa_printing_progress gauge
prusa_printing_progress{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 0
# HELP prusa_printing_time_remaining_seconds Returns time that remains for completion of current print
# TYPE prusa_printing_time_remaining_seconds gauge
prusa_printing_time_remaining_seconds{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name=""} 20100
# HELP prusa_status_info Returns information status of printer.
# TYPE prusa_status_info gauge
prusa_status_info{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name="",printer_state="Printing"} 4
# HELP prusa_tool_offset_temperature_celsius Offset tool temp
# TYPE prusa_tool_offset_temperature_celsius gauge
prusa_tool_offset_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name="",tool="0"} 0
# HELP prusa_tool_target_temperature_celsius Target tool temp
# TYPE prusa_tool_target_temperature_celsius gauge
prusa_tool_target_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name="",tool="0"} 170
# HELP prusa_tool_temperature_celsius Status of the printer tool temp
# TYPE prusa_tool_temperature_celsius gauge
prusa_tool_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="multiple_grots_0.4n_0.15mm_PLA,PLA,PLA,PLA_XLIS_5h36m.bgcode",printer_job_path="/usb/MULTIP~1.BGC",printer_model="MK4",printer_name="",tool="0"} 169
# HELP prusa_up Return information about online printers. If printer is registered as offline then returned value is 0.
# TYPE prusa_up gauge
prusa_up{printer_address="PRINTER_ADDRESS",printer_model="MK4",printer_name=""} 1
`
	output = strings.ReplaceAll(output, "PRINTER_ADDRESS", strings.TrimPrefix(server.URL, "http://"))
	if err := testutil.CollectAndCompare(collector, strings.NewReader(output)); err != nil {
		t.Fatal(err)
	}
}

func TestEinsyPrinter(t *testing.T) {
	server := runPrusaLinkAPIServer("api/einsy")
	defer server.Close()

	collector := prusalink.NewCollector(config.Config{
		Printers: []config.Printers{
			{
				Address: strings.TrimPrefix(server.URL, "http://"),
				Type:    "I3MK3S",
			},
		},
	})

	output := ` 
# HELP prusa_axis Returns information about position of axis.
# TYPE prusa_axis gauge
prusa_axis{printer_address="PRINTER_ADDRESS",printer_axis="x",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0
prusa_axis{printer_address="PRINTER_ADDRESS",printer_axis="y",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0
prusa_axis{printer_address="PRINTER_ADDRESS",printer_axis="z",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0.4
# HELP prusa_bed_offset_temperature_celsius Offset bed temp
# TYPE prusa_bed_offset_temperature_celsius gauge
prusa_bed_offset_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0
# HELP prusa_bed_target_temperature_celsius Target bed temp
# TYPE prusa_bed_target_temperature_celsius gauge
prusa_bed_target_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 60
# HELP prusa_bed_temperature_celsius Current temp of printer bed in Celsius
# TYPE prusa_bed_temperature_celsius gauge
prusa_bed_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 61.7
# HELP prusa_fan_speed_rpm Returns information about speed of hotend fan in rpm.
# TYPE prusa_fan_speed_rpm gauge
prusa_fan_speed_rpm{fan="hotend",printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 4080
prusa_fan_speed_rpm{fan="print",printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0
# HELP prusa_farm_mode Return if printer is set to farm mode
# TYPE prusa_farm_mode gauge
prusa_farm_mode{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 1
# HELP prusa_files_count Number of files in storage
# TYPE prusa_files_count gauge
prusa_files_count{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name="",printer_storage="PrusaLink gcodes"} 0
prusa_files_count{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name="",printer_storage="SD Card"} 80
# HELP prusa_info Returns information about printer.
# TYPE prusa_info gauge
prusa_info{api_version="0.9.0-legacy",printer_address="PRINTER_ADDRESS",printer_hostname="connect.prusa3d.com",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_location="Elf on a shelf",printer_model="I3MK3S",printer_name="",prusalink_name="MK3S with MMU3",serial_number="CZPX5222X004XK04220",server_version="0.7.2",version_text="PrusaLink 0.7.2"} 1
# HELP prusa_job_info Returns information about current job.
# TYPE prusa_job_info gauge
prusa_job_info{printer_address="PRINTER_ADDRESS",printer_job_id="0",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 1
# HELP prusa_material_info Returns information about loaded filament. Returns 0 if there is no loaded filament
# TYPE prusa_material_info gauge
prusa_material_info{printer_address="PRINTER_ADDRESS",printer_filament=" - ",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0
# HELP prusa_nozzle_size_meters Returns information about selected nozzle size.
# TYPE prusa_nozzle_size_meters gauge
prusa_nozzle_size_meters{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0.0004
# HELP prusa_print_flow_ratio Returns information about of filament flow in ratio (0.0 - 1.0).
# TYPE prusa_print_flow_ratio gauge
prusa_print_flow_ratio{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0.95
# HELP prusa_print_speed_ratio Current setting of printer speed in values from 0.0 - 1.0
# TYPE prusa_print_speed_ratio gauge
prusa_print_speed_ratio{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 1
# HELP prusa_print_time_seconds Returns information about current print time.
# TYPE prusa_print_time_seconds gauge
prusa_print_time_seconds{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0
# HELP prusa_printing_progress Returns information about completion of current print in percents
# TYPE prusa_printing_progress gauge
prusa_printing_progress{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 0
# HELP prusa_printing_time_remaining_seconds Returns time that remains for completion of current print
# TYPE prusa_printing_time_remaining_seconds gauge
prusa_printing_time_remaining_seconds{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name=""} 26160
# HELP prusa_status_info Returns information status of printer.
# TYPE prusa_status_info gauge
prusa_status_info{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name="",printer_state="Printing"} 4
# HELP prusa_tool_offset_temperature_celsius Offset tool temp
# TYPE prusa_tool_offset_temperature_celsius gauge
prusa_tool_offset_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name="",tool="0"} 0
# HELP prusa_tool_target_temperature_celsius Target tool temp
# TYPE prusa_tool_target_temperature_celsius gauge
prusa_tool_target_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name="",tool="0"} 215
# HELP prusa_tool_temperature_celsius Status of the printer tool temp
# TYPE prusa_tool_temperature_celsius gauge
prusa_tool_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_job_path="/SD Card/fosdem_0.2mm_PLA,PLA_MK3SMMU3_7h16m.gcode",printer_model="I3MK3S",printer_name="",tool="0"} 214.6
# HELP prusa_up Return information about online printers. If printer is registered as offline then returned value is 0.
# TYPE prusa_up gauge
prusa_up{printer_address="PRINTER_ADDRESS",printer_model="I3MK3S",printer_name=""} 1
`
	output = strings.ReplaceAll(output, "PRINTER_ADDRESS", strings.TrimPrefix(server.URL, "http://"))
	if err := testutil.CollectAndCompare(collector, strings.NewReader(output)); err != nil {
		t.Fatal(err)
	}
}

func TestSLPrinter(t *testing.T) {
	server := runPrusaLinkAPIServer("api/sl")
	defer server.Close()

	collector := prusalink.NewCollector(config.Config{
		Printers: []config.Printers{
			{
				Address: strings.TrimPrefix(server.URL, "http://"),
				Type:    "SL1",
			},
		},
	})

	output := `
# HELP prusa_ambient_temperature_celsius Status of the printer ambient temp
# TYPE prusa_ambient_temperature_celsius gauge
prusa_ambient_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 24.2
# HELP prusa_bed_offset_temperature_celsius Offset bed temp
# TYPE prusa_bed_offset_temperature_celsius gauge
prusa_bed_offset_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 0
# HELP prusa_bed_target_temperature_celsius Target bed temp
# TYPE prusa_bed_target_temperature_celsius gauge
prusa_bed_target_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 0
# HELP prusa_bed_temperature_celsius Current temp of printer bed in Celsius
# TYPE prusa_bed_temperature_celsius gauge
prusa_bed_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 51.1
# HELP prusa_chamber_offset_temperature_celsius Offset chamber temp
# TYPE prusa_chamber_offset_temperature_celsius gauge
prusa_chamber_offset_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 0
# HELP prusa_chamber_target_temperature_celsius Target chamber temp
# TYPE prusa_chamber_target_temperature_celsius gauge
prusa_chamber_target_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 0
# HELP prusa_chamber_temperature_celsius Status of the printer chamber temp
# TYPE prusa_chamber_temperature_celsius gauge
prusa_chamber_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 24.2
# HELP prusa_cover_status Status of the printer - 0 = open, 1 = closed
# TYPE prusa_cover_status gauge
prusa_cover_status{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 1
# HELP prusa_cpu_temperature_celsius Status of the printer cpu temp
# TYPE prusa_cpu_temperature_celsius gauge
prusa_cpu_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 51.1
# HELP prusa_fan_speed_rpm Returns information about speed of hotend fan in rpm.
# TYPE prusa_fan_speed_rpm gauge
prusa_fan_speed_rpm{fan="blower",printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 0
prusa_fan_speed_rpm{fan="rear",printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 0
prusa_fan_speed_rpm{fan="uv",printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 0
# HELP prusa_job_info Returns information about current job.
# TYPE prusa_job_info gauge
prusa_job_info{printer_address="PRINTER_ADDRESS",printer_job_id="0",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 1
# HELP prusa_status_info Returns information status of printer.
# TYPE prusa_status_info gauge
prusa_status_info{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name="",printer_state="Ready"} 1
# HELP prusa_tool_offset_temperature_celsius Offset tool temp
# TYPE prusa_tool_offset_temperature_celsius gauge
prusa_tool_offset_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name="",tool="0"} 0
# HELP prusa_tool_target_temperature_celsius Target tool temp
# TYPE prusa_tool_target_temperature_celsius gauge
prusa_tool_target_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name="",tool="0"} 0
# HELP prusa_tool_temperature_celsius Status of the printer tool temp
# TYPE prusa_tool_temperature_celsius gauge
prusa_tool_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name="",tool="0"} 26.5
# HELP prusa_up Return information about online printers. If printer is registered as offline then returned value is 0.
# TYPE prusa_up gauge
prusa_up{printer_address="PRINTER_ADDRESS",printer_model="SL1",printer_name=""} 1
# HELP prusa_uv_temperature_celsius Status of the printer uv temp
# TYPE prusa_uv_temperature_celsius gauge
prusa_uv_temperature_celsius{printer_address="PRINTER_ADDRESS",printer_job_name="",printer_job_path="",printer_model="SL1",printer_name=""} 26.5
`
	output = strings.ReplaceAll(output, "PRINTER_ADDRESS", strings.TrimPrefix(server.URL, "http://"))
	if err := testutil.CollectAndCompare(collector, strings.NewReader(output)); err != nil {
		t.Fatal(err)
	}
}

func runPrusaLinkAPIServer(baseFolder string) *httptest.Server {
	endPointFileMappings := map[string]string{
		"/api/v1/status":  "v1/status.json",
		"/api/v1/job":     "v1/job.json",
		"/api/v1/storage": "v1/storage.json",
		"/api/v1/info":    "v1/info.json",

		"/api/files":   "files.json",
		"/api/job":     "job.json",
		"/api/printer": "printer.json",
		"/api/version": "version.json",

		"/api/v1/cameras": "v1/cameras.json",
		"/api/settings":   "settings.json",
	}

	mux := http.NewServeMux()
	for endPoint, file := range endPointFileMappings {
		endPoint := endPoint
		file := file

		mux.HandleFunc(endPoint, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, baseFolder+"/"+file)
		})
	}

	return httptest.NewServer(mux)
}
