from veinmind import *
import os
import re
import jsonpickle
import pytoml as toml
from unittest import result
import requests

report_list = []
instruct_set = ("FROM", "CMD", "RUN", "LABEL", "MAINTAINER", "EXPOSE", "ENV", "ADD", "COPY", "ENTRYPOINT", "VOLUME", "USER", "WORKDIR", "ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL")
API_KEY = os.environ["API_KEY"]

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


def search_file_in_vt(id, key):
    url = "https://www.virustotal.com/api/v3/files/" + id
    headers = {
        "Accept": "application/json",
        "x-apikey": key
    }
    response = requests.request("GET", url, headers=headers)
    json_result = response.json()
    if 'error' in json_result.keys():
        # print (json_result['error']['message'])
        # print (f"File {id} is a legal file")
        return False
    elif 'data' in json_result.keys():
        last_analysis_stats = json_result['data']['attributes']['last_analysis_stats']
        if last_analysis_stats['harmless'] < last_analysis_stats['malicious']:
            # print (f"File {id} is a malicious file")
            return True
    return False


class Report():
    def __init__(self):
        self.imagename = ""
        self.abnormal_history_list = []


@command.group()
@command.option("--output", default="stdout", help="output format e.g. stdout/json")
def cli(output):
    load_rules()
    pass


@cli.image_command()
def scan_images(image):
    """scan image abnormal history instruction"""
    report = Report()
    refs = image.reporefs()
    if len(refs) > 0:
        ref = refs[0]
    else:
        ref = image.id()
    report.imagename = ref
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
                        if instruct == "COPY" or instruct == "ADD":
                            if len(command_split[1].split(":")) == 2:
                                file_id = command_split[1].split(":")[1]
                                # print(file_id)
                                if search_file_in_vt(file_id, API_KEY):
                                        report.abnormal_history_list.append(created_by)
                        command_content = command_split[1]
                        for r in rules["rules"]:
                            if r["instruct"] == instruct:
                                if re.match(r["match"], command_content):
                                    report.abnormal_history_list.append(created_by)
                                    break
                    else:
                        instruct = command_split[0]
                        command_content = " ".join(command_split[1:])
                        if instruct == "COPY" or instruct == "ADD":
                            if len(command_split) > 1 and len(command_split[1].split(":")) == 2:
                                file_id = command_split[1].split(":")[1]
                                # print(file_id)
                                if search_file_in_vt(file_id, API_KEY):
                                        report.abnormal_history_list.append(created_by)
                        for r in rules["rules"]:
                            if r["instruct"] == instruct:
                                if re.match(r["match"], command_content):
                                    report.abnormal_history_list.append(created_by)
                                    break
                else:
                    command_split = created_by.split()
                    if command_split[0] in instruct_set:
                        for r in rules["rules"]:
                            if r["instruct"] == command_split[0]:
                                if re.match(r["match"], " ".join(command_split[1:])):
                                    report.abnormal_history_list.append(created_by)
                                    break
                    else:
                        for r in rules["rules"]:
                            if r["instruct"] == "RUN":
                                if re.match(r["match"], created_by):
                                    report.abnormal_history_list.append(created_by)
                                    break
    report_list.append(report)


@cli.resultcallback()
def callback(result, output):
    if output == "stdout" and len(report_list) > 0:
        print("# ================================================================================================= #")
        tab_print("Scan Image Total: " + str(len(report_list)))
        tab_print("Unsafe Image List: ")
        for r in report_list:
            if len(r.abnormal_history_list) == 0:
                continue
            print("+---------------------------------------------------------------------------------------------------+")
            tab_print("ImageName: " + r.imagename)
            tab_print("Abnormal History Total: " + str(len(r.abnormal_history_list)))
            for history in r.abnormal_history_list:
                tab_print("History: " + history)
        print("+---------------------------------------------------------------------------------------------------+")
    elif output == "json":
        with open("output.json", mode="w") as f:
            f.write(jsonpickle.dumps(report_list))


if __name__ == '__main__':
    cli.add_info_command()
    cli()