import express, { Request, Response } from 'express';
import { chromium } from 'playwright-extra';
import stealthPlugin from 'puppeteer-extra-plugin-stealth';

chromium.use(stealthPlugin());

const app = express();
app.use(express.json());

interface JobOffer {
  id: string;
  platform: string;
  title: string;
  company: string;
  location: string;
  url: string;
}

app.post('/scrape', async (req: Request, res: Response) => {
  const { query, location } = req.body;

  if (!query) {
    return res.status(400).json({ error: 'Falta el parámetro query' });
  }

  console.log(`[LinkedIn Stealth] Buscando: ${query} en ${location || 'Global'}`);

  const browser = await chromium.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const context = await browser.newContext({
    userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36',
    viewport: { width: 1280, height: 800 },
    locale: 'es-ES'
  });

  const page = await context.newPage();

  try {
    const searchUrl = `https://www.linkedin.com/jobs/search?keywords=${encodeURIComponent(query)}&location=${encodeURIComponent(location || '')}`;
    
    // Navegar y esperar a que la red se estabilice
    await page.goto(searchUrl, { waitUntil: 'domcontentloaded', timeout: 35000 });

    // Scroll para cargar el contenido lazy-loaded
    await page.evaluate(() => window.scrollBy(0, 500));
    await page.waitForTimeout(2000);

    // Intentar esperar cualquier selector de lista válido
    await page.waitForSelector('ul.jobs-search__results-list, div.job-search-card, li', { timeout: 15000 }).catch(() => {
      console.log('[LinkedIn Stealth] Warning: No se detectó selector primario, intentando extracción fallback...');
    });

    const offers: JobOffer[] = await page.evaluate(() => {
      const cards = Array.from(document.querySelectorAll('ul.jobs-search__results-list > li, div.job-search-card, .base-card'));
      
      return cards.map((card) => {
        const titleElem = card.querySelector('.base-search-card__title, .job-card-list__title, h3');
        const companyElem = card.querySelector('.base-search-card__subtitle, .job-card-container__company-name, h4');
        const locElem = card.querySelector('.job-search-card__location, .job-card-container__metadata-item');
        const linkElem = card.querySelector('a.base-card__full-link, a.job-card-container__link, a') as HTMLAnchorElement;

        const url = linkElem ? linkElem.href : '';
        const idMatch = url.match(/view\/(\d+)/);

        return {
          id: idMatch ? idMatch[1] : '',
          platform: 'linkedin',
          title: titleElem ? titleElem.textContent?.trim() || '' : '',
          company: companyElem ? companyElem.textContent?.trim() || '' : '',
          location: locElem ? locElem.textContent?.trim() || '' : '',
          url: url
        };
      }).filter(o => o.title !== '' && o.url.includes('linkedin.com/jobs'));
    });

    await browser.close();
    console.log(`[LinkedIn Stealth] Éxito: ${offers.length} ofertas extraídas.`);
    return res.json({ success: true, count: offers.length, offers });

  } catch (err: any) {
    await browser.close();
    console.error('[LinkedIn Stealth] Error:', err.message);
    return res.status(500).json({ error: 'Error durante el scraping', details: err.message });
  }
});

app.post('/scrape-indeed', async (req: Request, res: Response) => {
  const { query, location } = req.body;

  if (!query) {
    return res.status(400).json({ error: 'Falta el parámetro query' });
  }

  console.log(`[Indeed Stealth] Buscando: ${query} en ${location || 'España'}`);

  const browser = await chromium.launch({
    headless: true,
    args: ['--no-sandbox', '--disable-setuid-sandbox']
  });

  const context = await browser.newContext({
    userAgent: 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36',
    locale: 'es-ES'
  });

  const page = await context.newPage();

  try {
    const searchUrl = `https://es.indeed.com/jobs?q=${encodeURIComponent(query)}&l=${encodeURIComponent(location || '')}`;
    await page.goto(searchUrl, { waitUntil: 'domcontentloaded', timeout: 30000 });

    await page.waitForSelector('.job_seen_beacon, td.resultContent, div.jobsearch-ResultsList', { timeout: 15000 }).catch(() => {
      console.log('[Indeed Stealth] Warning: No se encontró la lista de ofertas.');
    });

    const offers: JobOffer[] = await page.evaluate(() => {
      const cards = Array.from(document.querySelectorAll('div.job_seen_beacon, td.resultContent'));
      return cards.map((card) => {
        const titleElem = card.querySelector('h2.jobTitle span, a.jcs-JobTitle');
        const companyElem = card.querySelector('[data-testid="company-name"]');
        const locElem = card.querySelector('[data-testid="text-location"]');
        const linkElem = card.querySelector('a.jcs-JobTitle') as HTMLAnchorElement;

        const url = linkElem ? linkElem.href : '';
        const jobKey = linkElem ? linkElem.getAttribute('data-jk') || '' : '';

        return {
          id: jobKey,
          platform: 'indeed',
          title: titleElem ? titleElem.textContent?.trim() || '' : '',
          company: companyElem ? companyElem.textContent?.trim() || '' : '',
          location: locElem ? locElem.textContent?.trim() || '' : '',
          url: url
        };
      }).filter(o => o.title !== '' && o.url !== '');
    });

    await browser.close();
    return res.json({ success: true, count: offers.length, offers });

  } catch (err: any) {
    await browser.close();
    return res.status(500).json({ error: 'Error durante el scraping de Indeed', details: err.message });
  }
});

const PORT = process.env.PORT || 3001;
app.listen(PORT, () => {
  console.log(`🚀 Servicio LinkedIn Stealth escuchando en http://localhost:${PORT}`);
});

