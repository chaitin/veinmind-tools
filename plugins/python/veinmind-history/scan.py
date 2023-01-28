#!/usr/bin/env python3
from veinmind import *
import re
import pytoml as toml
from veinmind_common.service.report import *

report_list = []
instruct_set = (
    "FROM", "CMD", "RUN", "LABEL", "MAINTAINER", "EXPOSE", "ENV", "ADD", "COPY", "ENTRYPOINT", "VOLUME", "USER",
    "WORKDIR",
    "ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL")


def load_rules():
    global rules
    with open(os.path.join(os.path.abspath(os.path.dirname(__file__)), "rules.toml"), encoding="utf8") as f:
        rules = toml.load(f)


def tab_print(printstr: str):
    if len(printstr) < 95:
        print(("| " + printstr + "\t|").expandtabs(100))
    else:
        char_count = 0
        printstr_temp = ""
        for char in printstr:
            char_count = char_count + 1
            printstr_temp = printstr_temp + char
            if char_count == 95:
                char_count = 0
                print(("| " + printstr_temp + "\t|").expandtabs(100))
                printstr_temp = ""
        print(("| " + printstr_temp + "\t|").expandtabs(100))


@command.group()
@command.option("--output", default="stdout", help="output format e.g. stdout/json")
def cli(output):
    load_rules()
    pass


@cli.image_command()
def scan_images(image):
    """scan image abnormal history instruction"""
    image_report = None
    refs = image.reporefs()
    if len(refs) > 0:
        ref = refs[0]
    else:
        ref = image.id()
    log.info("start scan: " + ref)

    ocispec = image.ocispec_v1()
    if 'history' in ocispec.keys() and len(ocispec['history']) > 0:
        for history in ocispec['history']:
            if 'created_by' in history.keys():
                created_by = history['created_by']
                created_by_split = created_by.split("#(nop)")
                if len(created_by_split) > 1:
                    command = "#(nop)".join(created_by_split[1:])
                    command = command.lstrip()
                    command_split = command.split()
                    if len(command_split) == 2:
                        instruct = command_split[0]
                        command_content = command_split[1]
                        for r in rules["rules"]:
                            if r["instruct"] == instruct:
                                if re.match(r["match"], command_content):
                                    detail = AlertDetail()
                                    detail.history_detail = HistoryDetail(
                                                              instruction=instruct, content=command_content,
                                                              description=r["match"]
                                                          )
                                    image_report = ReportEvent(id=image.id(),
                                                          level=Level.High.value, detect_type=DetectType.Image.value,
                                                          event_type=EventType.Risk.value,
                                                          alert_type=AlertType.AbnormalHistory.value,
                                                          alert_details=[detail], native_object=image)
                                    report(image_report)
                                    break
                    else:
                        instruct = command_split[0]
                        command_content = " ".join(command_split[1:])
                        for r in rules["rules"]:
                            if r["instruct"] == instruct:
                                if re.match(r["match"], command_content):
                                    detail = AlertDetail()
                                    detail.history_detail = HistoryDetail(
                                                              instruction=instruct, content=command_content,
                                                              description=r["match"]
                                                          )
                                    image_report = ReportEvent(id=image.id(),
                                                          level=Level.High.value, detect_type=DetectType.Image.value,
                                                          event_type=EventType.Risk.value,
                                                          alert_type=AlertType.AbnormalHistory.value,
                                                          alert_details=[detail], native_object=image)
                                    report(image_report)
                                    break
                else:
                    command_split = created_by.split()
                    if command_split[0] in instruct_set:
                        for r in rules["rules"]:
                            if r["instruct"] == command_split[0]:
                                if re.match(r["match"], " ".join(command_split[1:])):
                                    detail = AlertDetail()
                                    detail.history_detail = HistoryDetail(
                                                              instruction=command_split[0],
                                                              content=" ".join(command_split[1:]),
                                                              description=r["match"]
                                                          )
                                    image_report = ReportEvent(id=image.id(),
                                                          level=Level.High.value, detect_type=DetectType.Image.value,
                                                          event_type=EventType.Risk.value,
                                                          alert_type=AlertType.AbnormalHistory.value,
                                                          alert_details=[detail], native_object=image)
                                    report(image_report)
                                    break
                    else:
                        for r in rules["rules"]:
                            if r["instruct"] == "RUN":
                                if re.match(r["match"], created_by):
                                    detail = AlertDetail()
                                    detail.history_detail = HistoryDetail(
                                                              instruction="RUN", content=created_by,
                                                              description=r["match"]
                                                          )
                                    image_report = ReportEvent(id=image.id(),
                                                          level=Level.High.value, detect_type=DetectType.Image.value,
                                                          event_type=EventType.Risk.value,
                                                          alert_type=AlertType.AbnormalHistory.value,
                                                          alert_details=[detail], native_object=image)
                                    report(image_report)
                                    break

    if image_report != None:
        report_list.append(image_report)


@cli.resultcallback()
def callback(result, output):
    if output == "stdout" and len(report_list) > 0:
        print("# ================================================================================================= #")
        tab_print("Scan Image Total: " + str(len(report_list)))
        tab_print("Unsafe Image List: ")
        for r in report_list:
            if len(r.alert_details) == 0:
                continue
            print(
                "+---------------------------------------------------------------------------------------------------+")
            tab_print("ImageName: " + r.id)
            tab_print("Abnormal History Total: " + str(len(r.alert_details)))
            for detail in r.alert_details:
                if detail.history_detail:
                    tab_print("History: " + detail.history_detail.content)
        print("+---------------------------------------------------------------------------------------------------+")
    elif output == "json":
        with open("output.json", mode="w") as f:
            f.write(jsonpickle.dumps(report_list))


if __name__ == '__main__':
    cli.add_info_command(manifest=command.Manifest(name="veinmind-history", author="veinmind-team", description="veinmind-history scan image abnormal history"))
    cli()