from pydantic import BaseModel, Field
from typing import List, Optional

class JobOfferInput(BaseModel):
    id: Optional[str] = None
    platform: str
    title: str
    company: str
    location: str
    url: str

class JobAnalysisResult(BaseModel):
    title: str
    company: str
    required_skills: List[str] = Field(description="Lista de habilidades técnicas requeridas")
    soft_skills: List[str] = Field(description="Lista de habilidades blandas mencionadas")
    seniority_level: str = Field(description="Nivel de experiencia estimado")
    match_score: int = Field(description="Compatibilidad con el CV de 0 a 100")
    missing_skills: List[str] = Field(description="Habilidades del empleo que faltan en el CV")
    key_matches: List[str] = Field(description="Puntos coincidentes entre el CV y la oferta")

class CVOptimizationResponse(BaseModel):
    overall_match_score: int = Field(description="Compatibilidad media de las ofertas")
    cv_improvements: List[str] = Field(description="Sugerencias para mejorar el CV")
    tailored_summary: str = Field(description="Resumen profesional sugerido para el CV")
    analyzed_jobs: List[JobAnalysisResult]