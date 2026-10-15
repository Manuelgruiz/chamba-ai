from pydantic import BaseModel
from typing import List, Dict, Any

class ParsePDFResponse(BaseModel):
    success: bool
    parsed_json: Dict[str, Any]
    base_latex: str

class AnalyzeRequest(BaseModel):
    cv_markdown: str
    jobs: List[Dict[str, Any]]