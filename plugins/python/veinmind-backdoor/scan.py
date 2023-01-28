#!/usr/bin/env python3
import register
import click
from common import log
from common import tools
from veinmind import *
from veinmind_common.service.report import *

results = []
start = 0
image_ids = []

@command.group()
@click.option('--format', default="stdout", help="output format e.g. stdout/json")
@click.option('--output', default='.', help="output path e.g. /tmp")
def cli(format, output):
    global start
    start = timep.time()
    pass


@cli.image_command()
def scan_images(image):
    """scan image backdoor"""
    global image_ids
    image_ids.append(image.id())
    if len(image.reporefs()) > 0:
        log.info("start scan: " + image.reporefs()[0])
    else:
        log.info("start scan: " + image.id())
    for plugin_name, plugin in register.register.plugin_dict.items():
        p = plugin()
        for r in p.detect(image):
            results.append(r)
            file_stat = image.stat(r.filepath)
            detail = AlertDetail.backdoor(backdoor_detail=BackdoorDetail(r.description, FileDetail.from_stat(r.filepath, file_stat)))
            report_event = ReportEvent(id=image.id(), level=Level.High.value,
                                       detect_type=DetectType.Image.value,
                                       event_type=EventType.Risk.value,
                                       alert_type=AlertType.Backdoor.value,
                                       alert_details=[detail], native_object=image)
            report(report_event)

@cli.resultcallback()
def callback(result, format, output):
    spend_time = timep.time() - start

    if format == "stdout":
        print("# ================================================================================================= #")
        tools.tab_print("Scan Image Total: " + str(len(image_ids)))
        tools.tab_print("Spend Time: " + spend_time.__str__() + "s")
        tools.tab_print("Backdoor Total: " + str(len(results)))
        for r in results:
            print("+---------------------------------------------------------------------------------------------------+")
            tools.tab_print("ImageName: " + r.image_ref)
            tools.tab_print("Backdoor File Path: " + r.filepath)
            tools.tab_print("Description: " + r.description)
        print("+---------------------------------------------------------------------------------------------------+")
        print("# ================================================================================================= #")

    elif format == "json":
        fpath = os.path.join(output, "report.json")
        with open(fpath, mode="w") as f:
            f.write(jsonpickle.dumps(results))


if __name__ == '__main__':
    cli.add_info_command(manifest=command.Manifest(name="veinmind-backdoor", author="veinmind-team", description="veinmind-backdoor scan image backdoor file"))
    cli()
