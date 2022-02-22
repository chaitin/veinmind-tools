from register import register
from common import log
from stat import *
from common import result
import os

@register.register("sshd")
class sshd():
    """
    sshd 软连接后门检测插件，支持检测常规软连接后门
    """
    rootok_list = ("su", "chsh", "chfn", "runuser")

    def detect(self, image):
        results = []

        for root, dirs, files in image.walk("/"):
            for f in files:
                try:
                    filepath = os.path.join(root, f)
                    f_lstat = image.lstat(filepath)
                    if S_ISLNK(f_lstat.st_mode):
                        f_link = image.evalsymlink(filepath)
                        f_exename = filepath.split("/")[-1]
                        f_link_exename = f_link.split("/")[-1]
                        if f_exename in self.rootok_list and f_link_exename == "sshd":
                            r = result.Result()
                            r.image_id = image.id()
                            if len(image.reporefs()) > 0:
                                r.image_ref = image.reporefs()[0]
                            else:
                                r.image_ref = image.id()
                            r.filepath = filepath
                            r.description = "sshd symlink backdoor"
                            results.append(r)
                except FileNotFoundError:
                    continue
                except BaseException as e:
                    log.logger.error(e)
        return results
