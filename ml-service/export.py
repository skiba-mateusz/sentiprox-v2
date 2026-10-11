import base64
import csv
import os
import sys
from datetime import datetime

import features as F
import pg8000.native

MANIFEST = sys.argv[1] if len(sys.argv) > 1 else "labels_manifest.csv"
OUT = sys.argv[2] if len(sys.argv) > 2 else "dataset.csv"

def connect():
    return pg8000.native.Connection(
        host=os.getenv("PGHOST", "localhost"),
        port=int(os.getenv("PGPORT", "5432")),
        user=os.getenv("PGUSER", "admin"),
        password=os.getenv("PGPASSWORD", "password"),
        database=os.getenv("PGDATABASE", "sentiprox"),
    )

def read_manifest(path):
    phases = []
    with open(path, newline="", encoding="utf-8") as fh:
        for row in csv.DictReader(fh):
            label = row["label"].strip()
            if label not in F.CONTENT_CLASSES:
                continue
            phases.append((
                label,
                datetime.fromisoformat(row["start"]),
                datetime.fromisoformat(row["end"]),
            ))
    return phases

def main():
    conn = connect()
    phases = read_manifest(MANIFEST)
    seen = set()
    rows = []

    for label, start, end in phases:
        result = conn.run(
            """
            SELECT features
            FROM requests_log
            WHERE features IS NOT NULL
              AND features != ''
              AND timestamp >= :s
              AND timestamp <= :e
            """,
            s=start, e=end,
        )

        if result == None:
            return

        for (content,) in result:
            key = (content, label)
            if key in seen:
                continue
            seen.add(key)
            rows.append(key)

    with open(OUT, "w", newline="", encoding="utf-8") as fh:
        writer = csv.writer(fh)
        writer.writerow(["label", "content"])

        for (content, label) in rows:
            if b"content-type: image/" in content.lower():
                content_str = base64.b64encode(content).decode('ascii')
            else:
                content_str = content.decode('utf-8', errors="replace")

            writer.writerow([label, content_str])

    counts = {}
    for _, label in rows:
        counts[label] = counts.get(label, 0) + 1
    print(f"wrote {len(rows)} to {OUT}")

    for label in sorted(counts):
        print(f"{label} - {counts[label]}")

if __name__ == "__main__":
    main()
