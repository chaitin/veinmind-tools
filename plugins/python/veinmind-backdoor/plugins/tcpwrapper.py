from register import register
from common import log
from common import result
from common import regex
import re

@register.register("tcpwrapper")
class tcpwrapper:
    """
    tcpwrapper 后门检测插件
    """
    wrapper_config_file_list = ['/etc/hosts.allow', '/etc/hosts.deny']

    def detect(self, image):
        results = []
        for config_filepath in self.wrapper_config_file_list:
            try:
                with image.open(config_filepath, mode="r") as f:
                    f_content = f.read()
                    for backdoor_regex in regex.backdoor_regex_list:
                        if re.search(backdoor_regex, f_content):
                            r = result.Result()
                            r.image_id = image.id()
                            if len(image.reporefs()) > 0:
                                r.image_ref = image.reporefs()[0]
                            else:
                                r.image_ref = image.id()
                            r.filepath = config_filepath
                            r.description = "regex: " + backdoor_regex
                            results.append(r)
            except FileNotFoundError:
                continue
            except BaseException as e:
                log.logger.error(e)
        return results
