import enum
import time
import jsonpickle
import json
import os, stat
from veinmind import service, log
from typing import List

_namespace = "github.com/chaitin/veinmind-tools/veinmind-common/go/service/report"

# Normalize timezone and format into RFC3339 format.
_timezone = time.strftime('%z')
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
    path: str
    perm: int
    size: int
    gid: int
    uid: int
    ctim: int
    mtim: int
    atim: int

    def __init__(self, path: str, perm: int, size: int, gid: int, uid: int, ctim: int, mtim: int, atim: int) -> None:
        self.path = path
        self.perm = perm
        self.size = size
        self.gid = gid
        self.uid = uid
        self.ctim = ctim
        self.mtim = mtim
        self.atim = atim

    @classmethod
    def from_stat(cls, path: str, file_stat: os.stat_result):
        return cls(path=path, perm=stat.S_IMODE(file_stat.st_mode), size=file_stat.st_size, gid=file_stat.st_gid,
                   uid=file_stat.st_uid, ctim=int(file_stat.st_ctime), mtim=int(file_stat.st_mtime),
                   atim=int(file_stat.st_atime))


class HistoryDetail():
    instruction: str
    content: str
    description: str

    def __init__(self, instruction: str, content: str, description: str):
        self.instruction = instruction
        self.content = content
        self.description = description


class SensitiveFileDetail(FileDetail):
    description: str

    def __init__(self, description: str, file_detail: FileDetail):
        self.description = description
        super().__init__(file_detail.path, file_detail.perm, file_detail.size, file_detail.gid, file_detail.uid,
                         file_detail.ctim, file_detail.mtim, file_detail.atim)


class SensitiveEnvDetail():
    key: str
    value: str
    description: str

    def __init__(self, key: str, value: str, description: str):
        self.key = key
        self.value = value
        self.description = description


class BackdoorDetail(FileDetail):
    description: str

    def __init__(self, description: str, file_detail: FileDetail):
        self.description = description
        super().__init__(file_detail.path, file_detail.perm, file_detail.size, file_detail.gid, file_detail.uid,
                         file_detail.ctim, file_detail.mtim, file_detail.atim)


class AlertDetail:
    backdoor_detail: BackdoorDetail
    sensitive_file_detail: SensitiveFileDetail
    sensitive_env_detail: SensitiveEnvDetail
    history_detail: HistoryDetail

    def __init__(self, backdoor_detail=None, sensitve_file_detail=None,
                 sensitive_env_detail=None, history_detail=None):
        self.backdoor_detail = backdoor_detail
        self.sensitive_file_detail = sensitve_file_detail
        self.sensitive_env_detail = sensitive_env_detail
        self.history_detail = history_detail

    @classmethod
    def backdoor(cls, backdoor_detail: BackdoorDetail):
        return cls(backdoor_detail=backdoor_detail)

    @classmethod
    def sensitive_file(cls, sensitve_file_detail: SensitiveFileDetail):
        return cls(sensitve_file_detail=sensitve_file_detail)

    @classmethod
    def sensitive_env(cls, sensitive_env_detail: SensitiveEnvDetail):
        return cls(sensitive_env_detail=sensitive_env_detail)

    @classmethod
    def history(cls, history_detail: HistoryDetail):
        return cls(history_detail=history_detail)


class ReportEvent():
    id: str
    time: str
    level: int
    detect_type: int
    event_type: int
    alert_type: int
    alert_details: List[AlertDetail]

    def __init__(self, id: str, level: int, detect_type: int, event_type: int, alert_type: int,
                 alert_details: List[AlertDetail], t: str = time.strftime(_format)) -> None:
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


def report(evt: ReportEvent, *args, **kwargs):
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

    def report(self, evt: ReportEvent):
        report(evt)
