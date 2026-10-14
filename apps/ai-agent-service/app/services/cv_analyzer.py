import os
import instructor
from openai import OpenAI
from typing import List
from app.schemas import JobOfferInput, CVOptimizationResponse

class CVAnalyzerService:
    def __init__(self):
        self.client = instructor.from_openai(
            OpenAI(
                base_url="https://api.groq.com/openai/v1",
                api_key=os.getenv("GROQ_API_KEY")
            )
        )

    def analyze_cv_against_jobs(self, cv_content: str, jobs: List[JobOfferInput]) -> CVOptimizationResponse:
        jobs_text = "\n".join([f"- {j.title} en {j.company} ({j.location}): {j.url}" for j in jobs[:5]])

        prompt = f"""
        Eres un experto reclutador tech y consultor de carrera.
        
        Contenido del CV del candidato:
        ---
        {cv_content}
        ---

        Ofertas de trabajo objetivo encontradas:
        ---
        {jobs_text}
        ---

        Analiza la compatibilidad del CV con las ofertas e indica sugerencias de optimización concretas.
        """

        return self.client.chat.completions.create(
            model="llama-3.3-70b-versatile",
            response_model=CVOptimizationResponse,
            messages=[
                {"role": "system", "content": "Eres un asistente especializado en optimización de CVs para ingenieros de software."},
                {"role": "user", "content": prompt}
            ],
            temperature=0.2
        )