import pytoml as toml
import click
import re
import os
import magic
import logging
import fnmatch
import chardet
from veinmind import *
from stat import *

# logger
formatter = logging.Formatter('%(asctime)s %(levelname)-8s %(message)s')
handler = logging.StreamHandler()
handler.setFormatter(formatter)
logger = logging.getLogger("veinmind-sensitive")
logger.setLevel(logging.INFO)
logger.addHandler(handler)

def load_rules():
    global rules
    with open(os.path.join(os.path.abspath(os.path.dirname(__file__)), "rules.toml"), encoding="utf8") as f:
        rules = toml.load(f)

@click.command()
@click.option('--engine',default="docker", help="engine type you use, e.g. docker/containerd")
@click.option('--name', default="", help="image name e.g. ubuntu:latest")
def cli(engine, name):
    load_rules()
    with runtime(engine) as client:
        if name != "":
            image_ids = client.find_image_ids(name)
        else:
            image_ids = client.list_image_ids()
        for id in image_ids:
            image = client.open_image_by_id(id)
            refs = image.reporefs()
            if len(refs) > 0:
                ref = refs[0]
            else:
                ref = image.id()
            logger.info("start scan: " + ref)
            for root, dirs, files in image.walk('/'):
                for filepath in files:
                    try:
                        filepath = os.path.join(root, filepath)

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

                        chardet_guess = chardet.detect(f_content_byte[0:64])
                        if chardet_guess["encoding"] != None:
                            try:
                                f_content = f_content_byte.decode(chardet_guess["encoding"])
                            except:
                                pass
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
                                            logger.warn("find sensitive file: " + filepath)
                                    else:
                                        if re.match(match, f_content):
                                            logger.warn("find sensitive file: " + filepath)
                    except Exception as e:
                        print(e)


def runtime(engine):
    if engine == "docker":
        return docker.Docker()
    elif engine == "containerd":
        return containerd.Containerd()
    else:
        raise Exception("engine type doesn't match")

if __name__ == '__main__':
    cli()