# TP4/TP5 — Plano de Testes e Resultados de Execução

Este documento consolida o **plano de testes** do ClimaDash (TP4) e o **registro
de execução** dos casos de teste (TP5). Cada caso de teste (TC) está vinculado a
um caso de uso (UC01–UC05, ver [TP1 — Casos de uso](TP1-casos-de-uso.md)) e a um
requisito quando aplicável.

## 1. Objetivo e escopo

Verificar que as funcionalidades descritas nos casos de uso estão corretamente
implementadas no backend (API REST em Go) e no frontend (dashboard em
HTML/CSS/JS). O escopo cobre:

- **Testes unitários** da camada de dados (`internal/data`) — parsing do CSV,
  filtros, agregações.
- **Testes de integração HTTP** dos handlers (`internal/handlers`) via
  `net/http/httptest`.
- **Testes manuais de interface** (UI) para os fluxos visuais e estados de erro.

## 2. Ambiente de testes

| Item | Valor |
|------|-------|
| Sistema operacional | Linux |
| Linguagem | Go 1.22 |
| Dataset | `data/emissions.csv` (14 países, anos 2013–2022) |
| Execução automatizada | `./scripts/run-tests.sh` (ou `go test ./...`) |
| Execução manual | Navegador em <http://localhost:8080> (`go run ./cmd/server`) |

Os testes automatizados não dependem do dataset real: cada teste monta um CSV
de *fixture* determinístico em diretório temporário, garantindo resultados
estáveis e independentes de alterações no dataset de produção.

## 3. Template padrão de caso de teste

Todo caso de teste segue o formato abaixo (campo **Tipo**: Automatizado ou
Manual):

| Campo | Descrição |
|-------|-----------|
| **ID** | Identificador único (`TC<UC>.<n>`) |
| **Caso de uso** | UC relacionado |
| **Tipo** | Automatizado (`go test`) ou Manual (UI) |
| **Pré-condições** | Estado necessário antes da execução |
| **Passos** | Sequência de ações |
| **Resultado esperado** | Comportamento correto |
| **Resultado obtido** | O que ocorreu na execução (preenchido no TP5) |
| **Status** | ✅ Passou / ❌ Falhou |

## 4. Estratégia e rastreabilidade

| Camada | Arquivo | Cobertura |
|--------|---------|-----------|
| Unitário — dados | `internal/data/repository_test.go` | parsing, cabeçalho/linha inválidos, `Countries`, `Years`, `LatestYear`, `Summary`, `Top`, `SeriesFor`, filtros, repositório vazio |
| Integração — HTTP | `internal/handlers/handlers_test.go` | `/api/health`, `/api/countries`, `/api/countries/{c}`, `/api/years`, `/api/summary`, `/api/emissions`, `/api/top`, método não permitido |
| Manual — UI | este documento (§5) | renderização de cartões/gráficos, estados vazios e mensagens de erro |

---

## 5. Casos de teste

> Coluna **Resultado obtido** e **Status** preenchidas na execução do TP5
> (ver §6 para o resumo e §7 para os defeitos encontrados).

### UC01 — Visualizar visão geral

| ID | Tipo | Passos | Resultado esperado | Resultado obtido | Status |
|----|------|--------|--------------------|------------------|--------|
| **TC01.1** | Automatizado (`TestSummaryEndpoint`) | `GET /api/summary` | `200`; `latest_year=2022`, `total_countries=14`, `global_co2_mt=25539.99`, `avg_renewable_share_pct=15.44` | Idêntico ao esperado | ✅ Passou |
| **TC01.2** | Manual | Abrir <http://localhost:8080> | Os 4 cartões (ano, países, emissão global, média renováveis) são preenchidos e os gráficos iniciais renderizam | Cartões e gráficos renderizados | ✅ Passou |
| **TC01.3** | Manual | Iniciar com `DATA_PATH` apontando para CSV só com cabeçalho e recarregar a página | Banner "Sem dados disponíveis" exibido; controles desabilitados (UC01 A1) | Falhou na 1ª execução (mostrava "0"/"—"); após correção exibe o banner | ✅ Passou (após DEF-01) |

### UC02 — Consultar série histórica de um país

| ID | Tipo | Passos | Resultado esperado | Resultado obtido | Status |
|----|------|--------|--------------------|------------------|--------|
| **TC02.1** | Automatizado (`TestCountryDetail_OK`) | `GET /api/countries/Brazil` | `200`; série do país com pontos ordenados por ano | Série ordenada retornada | ✅ Passou |
| **TC02.2** | Automatizado + Manual (`TestCountryDetail_NotFound`) | `GET /api/countries/Narnia` | `404`; frontend exibe "País não encontrado" (UC02 A1) | `404` retornado; mensagem exibida na UI | ✅ Passou |
| **TC02.3** | Automatizado (`TestCountryDetail_YearRange`) | `GET /api/countries/Brazil?from=2021&to=2022` | `200`; apenas pontos no intervalo | 2 pontos (2021, 2022) | ✅ Passou |

### UC03 — Comparar países

| ID | Tipo | Passos | Resultado esperado | Resultado obtido | Status |
|----|------|--------|--------------------|------------------|--------|
| **TC03.1** | Automatizado (`TestEmissions_OK`) | `GET /api/emissions?countries=Brazil,Atlantis` | `200`; uma série por país solicitado | 2 séries retornadas | ✅ Passou |
| **TC03.2** | Automatizado + Manual (`TestEmissions_TooMany`) | `GET /api/emissions?countries=a,b,c,d,e,f` | `400` "máximo de 5 países"; frontend impede 6ª marcação e avisa (UC03 A1) | `400` retornado; aviso exibido na UI | ✅ Passou |
| **TC03.3** | Automatizado (`TestEmissions_IgnoresUnknown`) | `GET /api/emissions?countries=Brazil,Narnia` | `200`; país inexistente é ignorado | Apenas Brazil retornado | ✅ Passou |

### UC04 — Filtrar por intervalo de anos

| ID | Tipo | Passos | Resultado esperado | Resultado obtido | Status |
|----|------|--------|--------------------|------------------|--------|
| **TC04.1** | Automatizado (`TestSeriesFor_YearRange`) | Filtrar série por `from`/`to` | Apenas pontos dentro do intervalo; `from` ou `to` isolados também funcionam | Filtros corretos | ✅ Passou |
| **TC04.2** | Automatizado (`TestCountryDetail_InvalidRange`) | `GET /api/countries/Brazil?from=2022&to=2020` | `400` "intervalo inválido: from > to" | `400` retornado | ✅ Passou |
| **TC04.3** | Manual | Na UI, definir "De" > "Até" e aplicar | Avisa intervalo inválido e reverte para o último intervalo válido (UC04 A1) | Falhou na 1ª execução (só exibia erro, sem reverter); após correção reverte | ✅ Passou (após DEF-02) |

### UC05 — Listar Top N emissores

| ID | Tipo | Passos | Resultado esperado | Resultado obtido | Status |
|----|------|--------|--------------------|--------------------|--------|
| **TC05.1** | Automatizado (`TestTop_OK`) | `GET /api/top?n=3` | `200`; 3 maiores emissores do último ano em ordem decrescente | No dataset real: China, United States, India | ✅ Passou |
| **TC05.2** | Automatizado (`TestTop_InvalidNDefaultsTo5`) | `GET /api/top?n=abc` | `200`; assume `N=5` (UC05 A1) | Default aplicado | ✅ Passou |
| **TC05.3** | Manual | Selecionar N=5 e N=10 na UI | Gráfico de barras horizontais atualiza com o Top N e o título reflete o N e o ano | Gráfico atualizado corretamente | ✅ Passou |

---

## 6. Resumo da execução (TP5)

| Métrica | Valor |
|---------|-------|
| Total de casos de teste | 15 |
| Automatizados | 10 |
| Manuais (UI) | 5 |
| Passaram | 15 |
| Falharam na 1ª rodada | 2 (TC01.3, TC04.3) |
| Defeitos encontrados | 2 |
| Defeitos corrigidos | 2 |
| Passaram após correção | 15 / 15 ✅ |

Os testes automatizados são executados com:

```bash
./scripts/run-tests.sh        # ou: go test ./...
```

Saída esperada (resumo):

```
ok  climadash/internal/data
ok  climadash/internal/handlers
```

## 7. Defeitos encontrados e corrigidos (TP5 — issue #49)

| ID | TC | Descrição | Correção | Arquivo |
|----|----|-----------|----------|---------|
| **DEF-01** | TC01.3 | Com dataset vazio, o frontend exibia "0"/"—" em vez do estado vazio previsto em UC01 A1. | Banner "Sem dados disponíveis" e desabilitação dos controles quando não há países/anos. | `web/static/app.js`, `index.html`, `style.css` |
| **DEF-02** | TC04.3 | Intervalo inválido (`De > Até`) apenas exibia mensagem de erro, sem reverter ao último intervalo válido (UC04 A1). | Armazenamento do último intervalo válido e reversão automática via `applyRange`. | `web/static/app.js` |

Ambos os defeitos foram corrigidos e os respectivos casos de teste reexecutados
com sucesso. As melhorias de UX correspondentes (mensagens de erro e estados
vazios) atendem também à issue **#45**.
