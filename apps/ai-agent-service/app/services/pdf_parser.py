import io
import json
import re
import requests
from pypdf import PdfReader

class PDFParserService:
    def __init__(self, lm_studio_url: str = "http://localhost:1234/v1/chat/completions", model_name: str = "qwen/qwen3-vl-4b"):
        self.lm_studio_url = lm_studio_url
        self.model_name = model_name

    def extract_text_from_pdf(self, pdf_bytes: bytes) -> str:
        reader = PdfReader(io.BytesIO(pdf_bytes))
        text = ""
        for page in reader.pages:
            extracted = page.extract_text()
            if extracted:
                text += extracted + "\n"
        return text

    def parse_pdf_to_json(self, pdf_bytes: bytes) -> dict:
        raw_text = self.extract_text_from_pdf(pdf_bytes)

        prompt = f"""
        Extrae la información del siguiente CV en texto plano y devuélvela estrictamente en formato JSON válido.
        
        INSTRUCCIÓN IMPORTANTE:
        - MANTÉN el contenido en su IDIOMA ORIGINAL.
        - NO utilices escapes hexadecimales de C como \\xfa o \\xf3. Escribe las tildes y caracteres en UTF-8 directamente (ej. "García", "Información").

        ESQUEMA JSON REQUERIDO:
        {{
          "full_name": "Nombre completo",
          "title": "Título profesional",
          "subtitle": "Subtítulo o especialidad",
          "location": "Ubicación",
          "phone": "Teléfono",
          "email": "Email",
          "linkedin": "LinkedIn profile",
          "github": "GitHub profile",
          "technical_skills": {{
              "Lenguajes": "...",
              "Backend & AI": "...",
              "Bases de Datos": "...",
              "Frontend": "...",
              "DevOps & Cloud": "..."
          }},
          "languages": ["Español: Nativo", "Inglés: B2"],
          "summary": "Perfil profesional / Summary",
          "experience": [
            {{
              "role": "Puesto",
              "company": "Empresa | Ubicación",
              "dates": "Fechas reales (ej. Jan 2025 -- Sep 2025)",
              "bullets": ["Bullet 1", "Bullet 2"]
            }}
          ],
          "education": [
            {{
              "degree": "Título obtenido",
              "institution": "Universidad / Institución",
              "dates": "Fechas"
            }}
          ]
        }}

        TEXTO DEL PDF:
        {raw_text}
        """

        payload = {
            "model": self.model_name,
            "messages": [
                {"role": "system", "content": "You are a JSON extraction engine. Return valid JSON only. Do not use \\x hex escapes for accented characters."},
                {"role": "user", "content": prompt}
            ],
            "temperature": 0.1
        }

        try:
            response = requests.post(self.lm_studio_url, json=payload, timeout=120)
            response.raise_for_status()
            
            result = response.json()
            content = result["choices"][0]["message"]["content"]
            
            # 1. Eliminar bloques de código markdown si el modelo los añade
            cleaned_content = content.replace("```json", "").replace("```", "").strip()
            
            # 2. Corregir escapes \xHEX inválidos en JSON (ej. \xfa -> \ufa o decodificación utf-8)
            cleaned_content = re.sub(r'\\x([0-9a-fA-F]{2})', r'\\u00\1', cleaned_content)

            return json.loads(cleaned_content)
        except requests.exceptions.RequestException as e:
            raise RuntimeError(f"Error conectando con LM Studio en local: {str(e)}")
        except (json.JSONDecodeError, KeyError) as e:
            raise RuntimeError(f"LM Studio no devolvió un JSON estructurado válido: {str(e)}")