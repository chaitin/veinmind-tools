import enum
import time as timep
import jsonpickle
import json
import os, stat
from veinmind import service, log

_namespace = "github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"

# Normalize timezone and format into RFC3339 format.
_timezone = timep.strftime('%z')
assert len(_timezone) == 5
if _timezone == "+0000" or _timezone == "-0000":
    _timezone = "Z"
else:
    _timezone = _timezone[0:3] + ':' + _timezone[3:5]
_format = "%Y-%m-%dT%H:%M:%S" + _timezone


@enum.unique
class Level(enum.Enum):
    Low = 0
    Medium = 1
    High = 2
    Critical = 3


@enum.unique
class DetectType(enum.Enum):
    Image = 0
    Container = 1


@enum.unique
class EventType(enum.Enum):
    Risk = 0
    Invasion = 1


@enum.unique
class AlertType(enum.Enum):
    Vulnerability = 0
    MaliciousFile = 1
    Backdoor = 2
    Sensitive = 3
    AbnormalHistory = 4
    Weakpass = 5


class FileDetail():
    path = ""
    perm = 0
    size = 0
    gid = 0
    uid = 0
    ctim = 0
    mtim = 0
    atim = 0

    def __init__(self, path, perm, size, gid, uid, ctim, mtim, atim) -> None:
        self.path = path
        self.perm = perm
        self.size = size
        self.gid = gid
        self.uid = uid
        self.ctim = ctim
        self.mtim = mtim
        self.atim = atim

    @classmethod
    def from_stat(cls, path, file_stat):
        return cls(path=path, perm=stat.S_IMODE(file_stat.st_mode), size=file_stat.st_size, gid=file_stat.st_gid,
                   uid=file_stat.st_uid, ctim=int(file_stat.st_ctime), mtim=int(file_stat.st_mtime),
                   atim=int(file_stat.st_atime))


class HistoryDetail():
    instruction = ""
    content = ""
    description = ""

    def __init__(self, instruction, content, description):
        self.instruction = instruction
        self.content = content
        self.description = description


class SensitiveFileDetail(FileDetail):
    description = ""

    def __init__(self, description, file_detail):
        self.description = description
        super().__init__(file_detail.path, file_detail.perm, file_detail.size, file_detail.gid, file_detail.uid,
                         file_detail.ctim, file_detail.mtim, file_detail.atim)


class SensitiveEnvDetail():
    key = ""
    value = ""
    description = ""

    def __init__(self, key, value, description):
        self.key = key
        self.value = value
        self.description = description


class BackdoorDetail(FileDetail):
    description = ""

    def __init__(self, description, file_detail):
        self.description = description
        super().__init__(file_detail.path, file_detail.perm, file_detail.size, file_detail.gid, file_detail.uid,
                         file_detail.ctim, file_detail.mtim, file_detail.atim)


class AlertDetail:
    backdoor_detail = None
    sensitive_file_detail = None
    sensitive_env_detail = None
    history_detail = None

    def __init__(self, backdoor_detail=None, sensitve_file_detail=None,
                 sensitive_env_detail=None, history_detail=None):
        self.backdoor_detail = backdoor_detail
        self.sensitive_file_detail = sensitve_file_detail
        self.sensitive_env_detail = sensitive_env_detail
        self.history_detail = history_detail

    @classmethod
    def backdoor(cls, backdoor_detail):
        return cls(backdoor_detail=backdoor_detail)

    @classmethod
    def sensitive_file(cls, sensitve_file_detail):
        return cls(sensitve_file_detail=sensitve_file_detail)

    @classmethod
    def sensitive_env(cls, sensitive_env_detail):
        return cls(sensitive_env_detail=sensitive_env_detail)

    @classmethod
    def history(cls, history_detail):
        return cls(history_detail=history_detail)


class ReportEvent():
    id = ""
    time = ""
    level = 0
    detect_type = 0
    event_type = 0
    alert_type = 0
    alert_details = []

    def __init__(self, id, level, detect_type, event_type, alert_type,
                 alert_details, t = timep.strftime(_format)):
        self.id = id
        self.time = t
        self.level = level
        self.detect_type = detect_type
        self.event_type = event_type
        self.alert_type = alert_type
        self.alert_details = alert_details


@service.service(_namespace, "report")
def _report(evt):
    pass


def report(evt, *args, **kwargs):
    if service.is_hosted():
        try:
            evt_dict = json.loads(jsonpickle.encode(evt))
            _report(evt_dict)
        except RuntimeError as e:
            log.error(e)
    else:
        log.warn(jsonpickle.encode(evt, indent=4))


class Entry:
    def __init__(self, **kwargs):
        self.fields = kwargs.copy()

    def report(self, evt):
        report(evt)
