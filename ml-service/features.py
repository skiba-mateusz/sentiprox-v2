from urllib.parse import unquote_plus

ATTACK_CLASSES = ["sqli", "xss", "cmdi", "traversal"]
CONTENT_CLASSES = ["benign"] + ATTACK_CLASSES


def normalize(s):
    cur = s
    for _ in range(3):
        dec = unquote_plus(cur)
        if dec == cur:
            break
        cur = dec
    return cur.lower()


def build_candidate(path="", query="", body="", ua="", referer="", cookie=""):
    raw = " ".join([path, query, body, ua, referer, cookie])
    return normalize(raw)
