from fastapi import FastAPI, HTTPException
from pydantic import BaseModel
from typing import List
from app.schemas import JobOfferInput, CVOptimizationResponse
from app.services.cv_analyzer import CVAnalyzerService

app = FastAPI(title="ChambaAI Agent Service", version="1.0.0")
analyzer_service = CVAnalyzerService()

class AnalyzeRequest(BaseModel):
    cv_markdown: str
    jobs: List[JobOfferInput]

@app.get("/health")
def health_check():
    return {"status": "ok", "service": "ai-agent-service"}

@app.post("/api/v1/ai/analyze-cv", response_model=CVOptimizationResponse)
def analyze_cv(payload: AnalyzeRequest):
    if not payload.cv_markdown.strip():
        raise HTTPException(status_code=400, detail="El contenido del CV no puede estar vacío.")
    
    if not payload.jobs:
        raise HTTPException(status_code=400, detail="Se requiere al menos una oferta para el análisis.")

    try:
        return analyzer_service.analyze_cv_against_jobs(payload.cv_markdown, payload.jobs)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error analizando el CV: {str(e)}")