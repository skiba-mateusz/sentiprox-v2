import csv
import os

import features as F
import joblib
from sklearn.feature_extraction.text import HashingVectorizer
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import classification_report, confusion_matrix
from sklearn.model_selection import train_test_split
from sklearn.pipeline import make_pipeline

DATASET = os.getenv("DATASET", "dataset.csv")
MODEL_PATH = os.getenv("MODEL_PATH", "model.joblib")


def load_csv(path):
    data, labels = [], []

    with open(path, newline="", encoding="utf-8") as fh:
        for row in csv.DictReader(fh):
            text = (row.get("content") or "").strip()
            label = (row.get("label") or "").strip()
            if text and label in F.CONTENT_CLASSES:
                data.append(text)
                labels.append(label)

    return data, labels

def load_training_data():
    d, l = load_csv(DATASET)
    return d, l

def main():
    data, labels = load_training_data()
    print(f"training ({len(data)} samples)", )

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

    pipe.fit(data, labels)
    joblib.dump(pipe, MODEL_PATH)
    print(f"\nsaved {MODEL_PATH} (refit on all {len(data)} samples)")

if __name__ == "__main__":
    main()
