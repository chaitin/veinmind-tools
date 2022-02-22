import register
import click
from veinmind import *
from plugins import *


@click.command()
@click.option('--engine',default="docker", help="engine type you use, e.g. docker/containerd")
@click.option('--name', default="", help="image name e.g. ubuntu:latest")
@click.option('--output', default="stdout", help="output format e.g. stdout/json")
def cli(engine, name, output):
    with runtime(engine) as client:
        if name != "":
            image_ids = client.find_image_ids(name)
        else:
            image_ids = client.list_image_ids()
        for id in image_ids:
            image = client.open_image_by_id(id)
            for plugin_name, plugin in register.register.plugin_dict.items():
                p = plugin()
                for r in p.detect(image):
                    print(r)

def runtime(engine):
    if engine == "docker":
        return docker.Docker()
    elif engine == "containerd":
        return containerd.Containerd()
    else:
        raise Exception("engine type doesn't match")

if __name__ == '__main__':
    cli()
