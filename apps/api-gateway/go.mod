module github.com/Manuelgruiz/chamba-ai/apps/api-gateway

go 1.25.0

// Indica a Go que resuelva el paquete localmente usando una ruta relativa
replace github.com/Manuelgruiz/chamba-ai/packages/shared-scrapers => ../../packages/shared-scrapers

require (
	github.com/Manuelgruiz/chamba-ai/packages/shared-scrapers v0.0.0-00010101000000-000000000000
	github.com/gofiber/fiber/v2 v2.52.15
)

require (
	github.com/PuerkitoBio/goquery v1.13.0 // indirect
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/andybalholm/cascadia v1.3.4 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/net v0.58.0 // indirect
	golang.org/x/sys v0.47.0 // indirect
)
