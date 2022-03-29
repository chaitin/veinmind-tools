#!/usr/bin/env python3
import pytoml as toml
import click
import re
import os
import magic
import fnmatch
import chardet
import time
import jsonpickle
from veinmind import *
from stat import *

report_list = []

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


class Report():
    def __init__(self):
        self.scan_counts = 0
        self.imagename = ""
        self.spend_time = 0
        self.sensitive_filepath_lists = []
        self.sensitive_env_lists = []


@command.group()
@command.option("--output", default="stdout", help="output format e.g. stdout/json")
def cli(output):
    load_rules()
    pass


@cli.image_command()
def scan_images(image):
    """scan image sensitive file"""
    report = Report()
    start = time.time()
    refs = image.reporefs()
    if len(refs) > 0:
        ref = refs[0]
    else:
        ref = image.id()
    report.imagename = ref
    log.info("start scan: " + ref)

    # detect env
    ocispec = image.ocispec_v1()
    if 'config' in ocispec.keys() and 'Env' in ocispec['config'].keys():
        env_list = image.ocispec_v1()['config']['Env']
        for env in env_list:
            env_split = env.split("=")
            if len(env_split) >= 2:
                for r in rules["rules"]:
                    if "env" in r.keys():
                        env_regex = r["env"]
                        if re.match(env_regex, env, re.IGNORECASE):
                            report.sensitive_env_lists.append(env)
                            log.warning("find sensitive env: " + env)
                            break

    for root, dirs, files in image.walk('/'):
        report.scan_counts = report.scan_counts + 1
        for dir in dirs:
            try:
                dirpath = os.path.join(root, dir)

                # detect filepath or filename
                for r in rules["rules"]:
                    if "filepath" in r.keys():
                        filepath_match_regex = r["filepath"]
                        if re.match(filepath_match_regex, dirpath):
                            report.sensitive_filepath_lists.append(dirpath)
                            log.warning("find sensitive filepath: " + dirpath)
                            break
            except Exception as e:
                print(e)

        for filename in files:
            try:
                filepath = os.path.join(root, filename)

                # skip whitelist
                whitelist = rules["whitelist"]
                white_match = False
                white_paths = whitelist["paths"]
                for wp in white_paths:
                    if fnmatch.filter([filepath], wp):
                        white_match = True
                        break
                if white_match:
                    continue

                try:
                    # skip not regular file
                    f_stat = image.stat(filepath)
                    if not S_ISREG(f_stat.st_mode):
                        continue

                    f = image.open(filepath, mode="rb")
                    f_content_byte = f.read()
                except FileNotFoundError as e:
                    continue
                except BaseException as e:
                    print(e)

                # detect filepath or filename
                match = False
                for r in rules["rules"]:
                    if "filepath" in r.keys():
                        filepath_match_regex = r["filepath"]
                        if re.match(filepath_match_regex, filepath):
                            report.sensitive_filepath_lists.append(filepath)
                            log.warning("find sensitive filepath: " + filepath)
                            match = True
                            break
                if match:
                    continue

                chardet_guess = chardet.detect(f_content_byte[0:64])
                if chardet_guess["encoding"] != None:
                    try:
                        f_content = f_content_byte.decode(chardet_guess["encoding"])
                    except:
                        continue
                else:
                    f_content = str(f_content_byte)
                mime_guess = magic.from_buffer(f_content_byte, mime=True)
                for r in rules["rules"]:
                    # mime
                    mime_find = False
                    if "mime" in r.keys():
                        if r["mime"] == mime_guess:
                            mime_find = True
                    else:
                        if mime_guess.startswith("text/"):
                            mime_find = True
                    if mime_find:
                        if "match" in r.keys():
                            match = r["match"]
                            if match.startswith("$contains:"):
                                keyword = match.lstrip("$contains:")
                                if keyword in f_content:
                                    log.warning("find sensitive file: " + filepath)
                            else:
                                if re.match(match, f_content):
                                    report.sensitive_filepath_lists.append(filepath)
                                    log.warning("find sensitive file: " + filepath)
            except Exception as e:
                print(e)
    spend_time = time.time() - start
    report.spend_time = spend_time
    report_list.append(report)

@cli.resultcallback()
def callback(result, output):
    if output == "stdout" and len(report_list) > 0:
        print("# ================================================================================================= #")
        tab_print("Scan Image Total: " + str(len(report_list)))
        spend_time_total = 0
        sensitive_file_total = 0
        for r in report_list:
            spend_time_total = spend_time_total + r.spend_time
            sensitive_file_total = sensitive_file_total + len(r.sensitive_filepath_lists)
        tab_print("Spend Time Total: " + spend_time_total.__str__() + "s")
        tab_print("Sensitive File Total: " + str(sensitive_file_total))
        tab_print("Unsafe Image List: ")
        for r in report_list:
            if len(r.sensitive_filepath_lists) == 0:
                continue
            print(
                "+---------------------------------------------------------------------------------------------------+")
            tab_print("ImageName: " + r.imagename)
            tab_print("Scan Total: " + str(r.scan_counts))
            tab_print("Spend Time: " + r.spend_time.__str__() + "s")
            tab_print("Sensitive File Total: " + str(len(r.sensitive_filepath_lists)))
            for fp in r.sensitive_filepath_lists:
                tab_print("Sensitive File: " + fp)
            for env in r.sensitive_env_lists:
                tab_print("Sensitive Env: " + env)
        print("+---------------------------------------------------------------------------------------------------+")
    elif output == "json":
        with open("output.json", mode="w") as f:
            f.write(jsonpickle.dumps(report_list))


if __name__ == '__main__':
    cli.add_info_command()
    cli()