from veinmind import *


@command.group()
def cli():
    pass


@cli.image_command()
def scan_image(image):
    """write your plugin scan action here"""
    print(image.id())


if __name__ == '__main__':
    cli.add_info_command(
        manifest=command.Manifest(name="veinmind-example", author="veinmind-team", description="veinmind-example"))
    cli()
