from fastapi import APIRouter, HTTPException, UploadFile, File
from app.schemas.cv import ParsePDFResponse, AnalyzeRequest
from app.services.pdf_parser import PDFParserService
from app.services.latex_renderer import LaTeXRendererService
from app.services.cv_analyzer import CVAnalyzerService

router = APIRouter(prefix="/api/v1/ai", tags=["cv"])

parser_service = PDFParserService()
renderer_service = LaTeXRendererService()
analyzer_service = CVAnalyzerService()

@router.post("/parse-pdf", response_model=ParsePDFResponse)
async def parse_pdf(file: UploadFile = File(...)):
    if not file.filename.endswith(".pdf"):
        raise HTTPException(status_code=400, detail="El archivo debe ser un PDF")

    try:
        contents = await file.read()
        json_data = parser_service.parse_pdf_to_json(contents)
        initial_latex = renderer_service.render_json_to_latex(json_data)

        return ParsePDFResponse(
            success=True,
            parsed_json=json_data,
            base_latex=initial_latex
        )
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error al procesar el archivo PDF: {str(e)}")

@router.post("/analyze-cv")
def analyze_cv(payload: AnalyzeRequest):
    if not payload.cv_markdown.strip():
        raise HTTPException(status_code=400, detail="El contenido del CV no puede estar vacío.")
    
    if not payload.jobs:
        raise HTTPException(status_code=400, detail="Se requiere al menos una oferta para el análisis.")

    try:
        return analyzer_service.analyze_cv_against_jobs(payload.cv_markdown, payload.jobs)
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Error analizando el CV: {str(e)}")