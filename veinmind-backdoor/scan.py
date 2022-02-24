#!/usr/bin/env python3
import register
import click
import jsonpickle
import time
from common import log
from common import tools
from veinmind import *
from plugins import *


@click.command()
@click.option('--engine',default="docker", help="engine type you use, e.g. docker/containerd")
@click.option('--name', default="", help="image name e.g. ubuntu:latest")
@click.option('--output', default="stdout", help="output format e.g. stdout/json")
def cli(engine, name, output):
    results = []

    start = time.time()
    with runtime(engine) as client:
        if name != "":
            image_ids = client.find_image_ids(name)
        else:
            image_ids = client.list_image_ids()
        for id in image_ids:
            image = client.open_image_by_id(id)
            if len(image.reporefs()) > 0:
                log.logger.info("start scan: " + image.reporefs()[0])
            else:
                log.logger.info("start scan: " + image.id())
            for plugin_name, plugin in register.register.plugin_dict.items():
                p = plugin()
                for r in p.detect(image):
                    results.append(r)
    spend_time = time.time() - start

    if output == "stdout" and len(results) > 0:
        print("# ================================================================================================= #")
        image_id_dict = {}
        for r in results:
            image_id_dict[r.image_id] = 0
        tools.tab_print("Scan Image Total: " + str(len(image_id_dict)))
        tools.tab_print("Spend Time: " + spend_time.__str__() + "s")
        tools.tab_print("Backdoor Total: " + str(len(results)))
        for r in results:
            print("+---------------------------------------------------------------------------------------------------+")
            tools.tab_print("ImageName: " + r.image_ref)
            tools.tab_print("Backdoor File Path: " + r.filepath)
            tools.tab_print("Description: " + r.description)
        print("+---------------------------------------------------------------------------------------------------+")
        print("# ================================================================================================= #")

    elif output == "json":
        with open("output.json", mode="w") as f:
            f.write(jsonpickle.dumps(results))

def runtime(engine):
    if engine == "docker":
        return docker.Docker()
    elif engine == "containerd":
        return containerd.Containerd()
    else:
        raise Exception("engine type doesn't match")

if __name__ == '__main__':
    cli()
