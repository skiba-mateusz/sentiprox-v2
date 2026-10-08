from contextlib import asynccontextmanager

import joblib
import uvicorn
from fastapi import FastAPI
from pydantic import BaseModel

model = None

@asynccontextmanager
async def lifespan(app: FastAPI):
    global model
    model = joblib.load("model.joblib")
    yield

app = FastAPI(lifespan=lifespan)

class ClassifyRequest(BaseModel):
    content: str

class ClassifyResponse(BaseModel):
    category: str
    score: float

@app.get("/health")
async def health():
    return {"status": "ok", "model_loaded": model is not None}

@app.post("/classify", response_model=ClassifyResponse)
async def classify(req: ClassifyRequest):
    if model is None:
        return {"category": "benign", "score": 0.0}

    proba = model.predict_proba([req.content])[0]
    classes = model.classes_
    idx = int(proba.argmax())
    return {"category": str(classes[idx]), "score": float(proba[idx])}

if __name__ == "__main__":
    uvicorn.run(app, host="0.0.0.0", port=8000)
