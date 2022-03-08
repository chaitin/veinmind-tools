from veinmind import *
import click
import os
import re
import jsonpickle
import pytoml as toml

def load_rules():
    global rules
    with open(os.path.join(os.path.abspath(os.path.dirname(__file__)), "rules.toml"), encoding="utf8") as f:
        rules = toml.load(f)


instruct_set = ("FROM", "CMD", "RUN", "LABEL", "MAINTAINER", "EXPOSE", "ENV", "ADD", "COPY", "ENTRYPOINT", "VOLUME", "USER", "WORKDIR", "ARG", "ONBUILD", "STOPSIGNAL", "HEALTHCHECK", "SHELL")


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


class Report():
    def __init__(self):
        self.imagename = ""
        self.abnormal_history_list = []


@click.command()
@click.option('--engine',default="docker", help="engine type you use, e.g. docker/containerd")
@click.option('--name', default="", help="image name e.g. ubuntu:latest")
@click.option('--output', default="stdout", help="output format e.g. stdout/json")
def cli(engine, name, output):
    load_rules()
    report_list = []
    with runtime(engine) as client:
        if name != "":
            image_ids = client.find_image_ids(name)
        else:
            image_ids = client.list_image_ids()

        for id in image_ids:
            report = Report()
            image = client.open_image_by_id(id)
            refs = image.reporefs()
            if len(refs) > 0:
                ref = refs[0]
            else:
                ref = image.id()
            report.imagename = ref
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
                                            report.abnormal_history_list.append(created_by)
                                            break
                            else:
                                instruct = command_split[0]
                                command_content = " ".join(command_split[1:])
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

def runtime(engine):
    if engine == "docker":
        return docker.Docker()
    elif engine == "containerd":
        return containerd.Containerd()
    else:
        raise Exception("engine type doesn't match")


if __name__ == '__main__':
    cli()