import json
import requests
from typing import List
from app.schemas.job import JobOfferInput, CVOptimizationResponse

class CVAnalyzerService:
    def __init__(self, lm_studio_url: str = "http://localhost:1234/v1/chat/completions", model_name: str = "qwen/qwen3-vl-4b"):
        self.lm_studio_url = lm_studio_url
        self.model_name = model_name

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
        
        Devuelve la respuesta estrictamente en un objeto JSON con el siguiente esquema:
        {{
            "overall_match_score": 85,
            "strengths": ["Fortaleza 1", "Fortaleza 2"],
            "missing_skills": ["Habilidad faltante 1", "Habilidad faltante 2"],
            "suggestions": ["Sugerencia de mejora 1", "Sugerencia de mejora 2"]
        }}
        """

        payload = {
            "model": self.model_name,
            "messages": [
                {"role": "system", "content": "Eres un asistente especializado en optimización de CVs. Devuelve ÚNICAMENTE JSON válido."},
                {"role": "user", "content": prompt}
            ],
            "temperature": 0.2,
            "response_format": {"type": "json_object"}
        }

        try:
            response = requests.post(self.lm_studio_url, json=payload, timeout=120)
            response.raise_for_status()
            
            result = response.json()
            content = result["choices"][0]["message"]["content"]
            parsed_data = json.loads(content)
            
            # Convierte la respuesta en el modelo Pydantic esperado
            return CVOptimizationResponse(**parsed_data)
        except requests.exceptions.RequestException as e:
            raise RuntimeError(f"Error conectando con LM Studio durante el análisis: {str(e)}")
        except (json.JSONDecodeError, KeyError, Exception) as e:
            raise RuntimeError(f"No se pudo parsear el análisis del CV a JSON: {str(e)}")