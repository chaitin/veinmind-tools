from register import register
from common import log
from common import result
from common import regex
import os
import re

@register.register("bashrc")
class bashrc:
    """
    bashrc 后门检测插件
    """
    backdoor_regex_list = [
        r'''alias\s+ssh=[\'\"]{0,1}strace''',
        r'''alias\s+sudo='''
    ]

    bashrc_dirs = [
        "/home",
        "/root"
    ]

    def detect(self, image):
        results = []

        for bashrc_dir in self.bashrc_dirs:
            for root, dirs, files in image.walk(bashrc_dir):
                for file in files:
                    if re.match(r'''^\.[\w]*shrc$''', file):
                        filepath = os.path.join(root, file)
                    else:
                        continue
                    try:
                        f = image.open(filepath, mode="r")
                        f_content = f.read()
                        for backdoor_regex in self.backdoor_regex_list:
                            if re.search(backdoor_regex, f_content):
                                r = result.Result()
                                r.image_id = image.id()
                                if len(image.reporefs()) > 0:
                                    r.image_ref = image.reporefs()[0]
                                else:
                                    r.image_ref = image.id()
                                r.filepath = filepath
                                r.description = "regex: " + backdoor_regex
                                results.append(r)
                        for backdoor_regex in regex.backdoor_regex_list:
                            if re.search(backdoor_regex, f_content):
                                r = result.Result()
                                r.image_id = image.id()
                                if len(image.reporefs()) > 0:
                                    r.image_ref = image.reporefs()[0]
                                else:
                                    r.image_ref = image.id()
                                r.filepath = filepath
                                r.description = "regex: " + backdoor_regex
                                results.append(r)
                    except FileNotFoundError:
                        continue
                    except BaseException as e:
                        log.logger.error(e)
        return results
