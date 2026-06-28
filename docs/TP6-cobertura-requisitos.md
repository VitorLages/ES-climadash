# TP6 — Verificação de Cobertura dos Requisitos

Este documento verifica que **todos** os requisitos funcionais (RF01–RF11) e não
funcionais (RNF01–RNF08) definidos no [TP1 — Requisitos](TP1-requisitos.md) estão
implementados e, quando aplicável, cobertos por testes
([TP4/TP5 — Plano de testes](TP4-plano-de-testes.md)).

Legenda da coluna **Situação**: ✅ Atendido · ⚠️ Atendido com ressalva.

## 1. Requisitos Funcionais

| RF | Requisito | Onde é implementado | Evidência / Teste | Situação |
|----|-----------|---------------------|-------------------|----------|
| **RF01** | Carregar dados de CSV | `internal/data/repository.go` → `LoadFromCSV` | `TestLoadFromCSV_OK`, `TestLoadFromCSV_BadHeader`, `TestLoadFromCSV_BadRow` | ✅ |
| **RF02** | Listar países via API | `GET /api/countries` → `Repository.Countries` | `TestCountriesEndpoint`, `TestCountries_Sorted` | ✅ |
| **RF03** | Listar anos via API | `GET /api/years` → `Repository.Years` | `TestYearsEndpoint`, `TestYearsAndLatest` | ✅ |
| **RF04** | Série histórica de um país | `GET /api/countries/{country}` → `Repository.SeriesFor` | `TestCountryDetail_OK` (TC02.1) | ✅ |
| **RF05** | Filtrar por intervalo de anos | parâmetros `from`/`to` → `filterByYear` | `TestCountryDetail_YearRange`, `TestSeriesFor_YearRange` (TC04.1) | ✅ |
| **RF06** | Comparar até 5 países | `GET /api/emissions?countries=` (limite 5) → `SeriesForMany` | `TestEmissions_OK`, `TestEmissions_TooMany` (TC03.1/03.2) | ✅ |
| **RF07** | Visão geral com agregados | `GET /api/summary` → `Repository.Summary` | `TestSummaryEndpoint`, `TestSummary` (TC01.1) | ✅ |
| **RF08** | Gráfico de evolução de CO₂ | frontend `drawCompare` (métrica `co2_total_mt`) | TC01.2 (manual) | ✅ |
| **RF09** | Gráfico de participação de renováveis | frontend seletor de métrica `renewable_share_pct` em `drawCompare` | TC01.2 (manual) | ✅ |
| **RF10** | Top N maiores emissores | `GET /api/top?n=N` → `Repository.Top` | `TestTop_OK`, `TestTop_InvalidNDefaultsTo5` (TC05.1/05.2) | ✅ |
| **RF11** | Atualizar gráficos após filtros | frontend `refreshCompare`/`refreshTop` nos eventos de "Aplicar"/mudança | TC05.3 (manual) | ✅ |

**Resultado:** 11/11 requisitos funcionais atendidos.

## 2. Requisitos Não Funcionais

| RNF | Tipo | Requisito | Evidência | Situação |
|-----|------|-----------|-----------|----------|
| **RNF01** | Performance | API responde < 300 ms (dataset embarcado) | Dados carregados em memória no boot; consultas são *lookups*/agregações O(n) sobre ~140 registros. Tempo de cada requisição é registrado pelo middleware `logging` em `cmd/server/main.go`. | ✅ |
| **RNF02** | Usabilidade | Usável em telas ≥ 1024×768 | Layout responsivo (`grid`/`flex`) e *media query* em `web/static/style.css`; `viewport` declarado no HTML. | ✅ |
| **RNF03** | Manutenibilidade | Código em pacotes `cmd/`, `internal/handlers`, `internal/data`, `internal/models` | Estrutura de pastas do repositório espelha exatamente a divisão exigida. | ✅ |
| **RNF04** | Portabilidade | Compila em Linux/macOS/Windows com Go 1.22+ | Apenas biblioteca padrão (sem CGO/dependências externas); `go.mod` fixa `go 1.22`. | ✅ |
| **RNF05** | Disponibilidade | Inicia e responde em < 2 s após `go run` | Carga única do CSV em memória na inicialização; sem I/O por requisição. | ✅ |
| **RNF06** | Confiabilidade | Erros de dataset reportados em log sem derrubar o servidor após o boot | Falha de carga é reportada com `log.Fatalf` **antes** do boot (RNF06 refere-se ao pós-boot); após iniciar, os dados estão em memória e nenhuma requisição relê o CSV, logo não há caminho de falha que derrube o servidor. | ⚠️ |
| **RNF07** | Segurança | API expõe apenas operações de leitura (GET) | Todos os handlers respondem `405 Method Not Allowed` a métodos diferentes de GET; sem upload/autenticação/estado mutável. | ✅ — `TestMethodNotAllowed` |
| **RNF08** | Open-source | Repositório público, docs em Markdown, dependências só de stdlib + Chart.js | Repositório público no GitHub; `docs/` em Markdown; Chart.js via CDN; backend sem dependências externas. | ✅ |

**Resultado:** 8/8 requisitos não funcionais atendidos (RNF06 atendido com a
ressalva esclarecida acima).

## 3. Conclusão

Todos os requisitos elicitados no TP1 estão implementados e rastreáveis ao
código e aos testes. A cobertura automatizada de testes é de **95,6%** em
`internal/data` e **80,0%** em `internal/handlers` (ver
`./scripts/run-tests.sh -c`). Os requisitos predominantemente visuais (RF08,
RF09, RF11) são verificados pelos casos de teste manuais do plano de testes.
