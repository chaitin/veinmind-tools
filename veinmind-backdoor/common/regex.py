backdoor_regex_list = [
    # reverse shell
    r'''(nc|ncat|netcat)\b.*(-e|--exec|-c)\b.*?\b(ba|da|z|k|c|a|tc|fi|sc)?sh\b''',
    r'''python\w*\b.*\bsocket\b.*?\bconnect\b.*?\bsubprocess\b.*?\bsend\b.*?\bstdout\b.*?\bread\b''',
    r'''python\w*\b.*\bsocket\b.*?\bconnect\b.*?\bos\.dup2\b.*?\b(call|spawn|popen)\b\s*\([^)]+?\b(ba|da|z|k|c|a|tc|fi|sc)?sh\b''',
    r'''(sh|bash|dash|zsh)\b.*-c\b.*?\becho\b.*?\bsocket\b.*?\bwhile\b.*?\bputs\b.*?\bflush\b.*?\|\s*tclsh\b''',
    r'''(sh|bash|dash|zsh)\b.*-c\b.*?\btelnet\b.*?\|&?.*?\b(ba|da|z|k|c|a|tc|fi|sc)?sh\b.*?\|&?\s*telnet\b''',
    r'''(sh|bash|dash|zsh)\b.*-c\b.*?\bcat\b.*?\|&?.*?\b(ba|da|z|k|c|a|tc|fi|sc)?sh\b.*?\|\s*(nc|ncat)\b''',
    r'''(sh|bash|dash|zsh)\b.*sh\s+(-i)?\s*>&?\s*/dev/(tcp|udp)/.*?/\d+\s+0>&\s*(1|2)''',
]