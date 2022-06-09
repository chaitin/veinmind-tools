import logging

# logger
formatter = logging.Formatter('[%(levelname)s] %(message)s')
handler = logging.StreamHandler()
handler.setFormatter(formatter)
logger = logging.getLogger("veinmind-sensitive")
logger.setLevel(logging.INFO)
logger.addHandler(handler)