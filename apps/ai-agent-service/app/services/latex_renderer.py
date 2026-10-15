class LaTeXRendererService:
    @staticmethod
    def render_json_to_latex(cv_data: dict) -> str:
        is_english = any("Languages" in k or "Summary" in k for k in cv_data.get("technical_skills", {}).keys())
        
        sec_contact = "Contact" if is_english else "Contacto"
        sec_skills = "Technical Skills" if is_english else "Habilidades Técnicas"
        sec_langs = "Languages" if is_english else "Idiomas"
        sec_summary = "Professional Summary" if is_english else "Perfil Profesional"
        sec_exp = "Work Experience" if is_english else "Experiencia Laboral"
        sec_edu = "Education" if is_english else "Educación"

        skills_lines = [
            f"\\textbf{{{cat}}} \\\\\n{items} \\\\[8pt]"
            for cat, items in cv_data.get("technical_skills", {}).items()
        ]
        skills_block = "\n".join(skills_lines)

        langs_lines = []
        for l in cv_data.get("languages", []):
            parts = l.split(":")
            if len(parts) == 2:
                langs_lines.append(f"\\textbf{{{parts[0].strip()}:}} {parts[1].strip()}")
            else:
                langs_lines.append(l)
        langs_block = " \\\\[4pt]\n".join(langs_lines)

        exp_blocks = []
        for exp in cv_data.get("experience", []):
            bullets = "\n".join([f"  \\item {b}" for b in exp.get("bullets", [])])
            block = f"""\\textbf{{{exp.get('role', '')}}} \\hfill \\textit{{{exp.get('dates', '')}}} \\\\
\\textit{{{exp.get('company', '')}}}
\\begin{{itemize}}[leftmargin=12pt, topsep=4pt, itemsep=3pt]
  \\small
{bullets}
\\end{{itemize}}
\\vspace{{10pt}}"""
            exp_blocks.append(block)
        exp_block = "\n".join(exp_blocks)

        edu_blocks = [
            f"\\textbf{{{edu.get('degree', '')}}} \\hfill \\textit{{{edu.get('dates', '')}}} \\\\\n\\textit{{{edu.get('institution', '')}}} \\\\[6pt]"
            for edu in cv_data.get("education", [])
        ]
        edu_block = "\n".join(edu_blocks)

        babel_lang = "english" if is_english else "spanish"

        return f"""\\documentclass[10pt, a4paper]{{article}}
\\usepackage[utf8]{{inputenc}}
\\usepackage[{babel_lang}]{{babel}}
\\usepackage[empty]{{fullpage}}
\\usepackage{{geometry}}
\\usepackage{{xcolor}}
\\usepackage{{enumitem}}
\\usepackage{{hyperref}}
\\usepackage{{titlesec}}

\\definecolor{{primary}}{{HTML}}{{1A2B4C}}
\\definecolor{{secondary}}{{HTML}}{{2C3E50}}
\\definecolor{{textgray}}{{HTML}}{{4A5568}}

\\geometry{{left=0.8cm, top=0.7cm, right=0.8cm, bottom=0.7cm}}
\\pagestyle{{empty}}
\\setlength{{\\parindent}}{{0pt}}

\\titleformat{{\\section}}{{\\color{{primary}}\\large\\bfseries\\scshape}}{{}}{{0em}}{{}}[\\color{{primary}}\\titlerule\\vspace{{-2pt}}]

\\begin{{document}}

\\begin{{minipage}}[t]{{1.0\\textwidth}}
    {{\\Huge \\bfseries \\color{{primary}} {cv_data.get('full_name', '')}}} \\\\[6pt]
    {{\\Large \\bfseries \\color{{secondary}} {cv_data.get('title', '')}}} \\\\[4pt]
    {{\\small \\color{{textgray}} {cv_data.get('subtitle', '')}}}
\\end{{minipage}}

\\vspace{{16pt}}

\\begin{{minipage}}[t]{{0.33\\textwidth}}

    \\section{{{sec_contact}}}
    \\small
    \\textbf{{Ubicación:}} {cv_data.get('location', '')} \\\\[4pt]
    \\textbf{{Teléfono:}} {cv_data.get('phone', '')} \\\\[4pt]
    \\textbf{{Correo:}} \\href{{mailto:{cv_data.get('email', '')}}}{{{cv_data.get('email', '')}}} \\\\[4pt]
    \\textbf{{LinkedIn:}} {cv_data.get('linkedin', '')} \\\\[4pt]
    \\textbf{{GitHub:}} {cv_data.get('github', '')}

    \\vspace{{20pt}}

    \\section{{{sec_skills}}}
    \\small
    {skills_block}

    \\vspace{{20pt}}

    \\section{{{sec_langs}}}
    \\small
    {langs_block}

\\end{{minipage}}
\\hfill
\\begin{{minipage}}[t]{{0.63\\textwidth}}

    \\section{{{sec_summary}}}
    \\small \\color{{textgray}}
    {cv_data.get('summary', '')}

    \\vspace{{16pt}}

    \\section{{{sec_exp}}}
    {exp_block}

    \\vspace{{16pt}}

    \\section{{{sec_edu}}}
    {edu_block}

\\end{{minipage}}

\\end{{document}}"""