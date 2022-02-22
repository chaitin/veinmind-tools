from register import register
from common import log
from common import result
from common import regex
import re
from stat import *
import os

@register.register("crontab")
class crontab:
    """
    crontab 后门检测插件
    """
    cron_list = ["/etc/crontab", "/etc/cron.hourly", "/etc/cron.daily" ,"/etc/cron.weekly", "/etc/cron.monthly", "/etc/cron.d"]
    environment_regex = '''[a-zA-Z90-9]+\s*=\s*[^\s]+$'''
    cron_regex = '''((\d{1,2}|\*)\s+){5}[a-zA-Z0-9]+\s+(.*)'''
    backdoor_regex_list = [
        # download
        r'''^(wget|curl)\b'''
        
        # mrig
        r'''^([\w0-9]*mrig[\w0-9]*)\b'''
    ]

    def detect(self, image):
        results = []
        for cron in self.cron_list:
            try:
                # filetype
                cron_stat = image.stat(cron)
                if S_ISDIR(cron_stat.st_mode):
                    for root, dirs, files in image.walk(cron):
                        for file in files:
                            filepath = os.path.join(root, file)
                            with image.open(filepath) as f:
                                result_dict = self.detect_crontab_content(f)
                                if len(result_dict) > 0:
                                    for regex, cmdline in result_dict.items():
                                        r = result.Result()
                                        r.image_id = image.id()
                                        r.image_ref = image.reporefs()[0]
                                        r.filepath = cron
                                        r.description = "regex: {0}, cmdline: {1}".format(regex, cmdline)
                                        results.append(r)
                elif S_ISREG(cron_stat.st_mode):
                    with image.open(cron) as f:
                        result_dict = self.detect_crontab_content(f)
                        if len(result_dict) > 0:
                            for regex, cmdline in result_dict.items():
                                r = result.Result()
                                r.image_id = image.id()
                                if len(image.reporefs()) > 0:
                                    r.image_ref = image.reporefs()[0]
                                else:
                                    r.image_ref = image.id()
                                r.filepath = cron
                                r.description = "regex: {0}, cmdline: {1}".format(regex, cmdline)
                                results.append(r)

            except FileNotFoundError:
                continue
        return results

    def detect_crontab_content(self, cron_f):
        result_dict = {}
        for line in cron_f.readlines():
            # preprocess
            line = line.strip()
            line = line.replace("\n", "")

            # environment
            if re.match(self.environment_regex, line):
                continue
            m = re.match(self.cron_regex, line)
            if m:
                if len(m.groups()) == 3:
                    cmdline = m.group(3)
                    for backdoor_regex in self.backdoor_regex_list:
                        if re.search(backdoor_regex, cmdline):
                            result_dict[backdoor_regex] = cmdline
                    for backdoor_regex in regex.backdoor_regex_list:
                        if re.search(backdoor_regex, cmdline):
                            result_dict[backdoor_regex] = cmdline
                else:
                    continue
        return result_dict