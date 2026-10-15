from fastapi import FastAPI
from app.routers import cv

app = FastAPI(title="ChambaAI Agent Service", version="1.0.0")

# Incluir las rutas modularizadas
app.include_router(cv.router)

@app.get("/health")
def health_check():
    return {"status": "ok", "service": "ai-agent-service"}