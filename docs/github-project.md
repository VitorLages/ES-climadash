# Quadro de Tarefas — GitHub Projects (espelho)

Este documento espelha o quadro do **GitHub Projects** do repositório.
Cada item é um *card* a ser criado no Projects, organizado por coluna:

- **Project Backlog** — itens previstos, ainda não puxados para sprint.
- **TODO** — itens planejados para o **próximo sprint**.
- **In Progress** — em execução no sprint corrente.
- **Done** — concluídos.

Conforme o enunciado, ao final de cada TP a coluna **TODO** deve conter as tarefas do sprint seguinte (na entrega do TP1, TODO = tarefas do TP2; na do TP2, TODO = TP3; etc.). Esta é a **entrega final (TP1 → TP6)**, correspondente à release **v1.0.0**: todos os cards estão em **Done** e não há sprints subsequentes planejados.

> **Como usar:** crie um Project (Board) no repositório com as colunas acima e replique os cards desta tabela. Conforme avançar, altere o status diretamente no GitHub (a fonte da verdade é o Projects; este arquivo é um snapshot inicial).

---

## Sprint TP1 — Definição do Problema e Planejamento Inicial

| # | Card | Sprint | Status |
|---|------|--------|--------|
| 1 | Definir ODS abordado (ODS 13 com interface ODS 7) | TP1 | Done |
| 2 | Descrever problema a ser resolvido | TP1 | Done |
| 3 | Escolher e justificar tipo de solução (dashboard de dados) | TP1 | Done |
| 4 | Elicitar requisitos funcionais (RF01–RF11) | TP1 | Done |
| 5 | Elicitar requisitos não-funcionais (RNF01–RNF08) | TP1 | Done |
| 6 | Modelar diagrama de casos de uso (UC01–UC05) | TP1 | Done |
| 7 | Criar repositório público no GitHub | TP1 | Done |
| 8 | Configurar GitHub Projects (colunas e cards iniciais) | TP1 | Done |
| 9 | Escrever README do projeto | TP1 | Done |

## Sprint TP2 — Projeto de Software

| # | Card | Sprint | Status |
|---|------|--------|--------|
| 10 | Definir stack tecnológica (Go stdlib + Chart.js) | TP2 | Done |
| 11 | Elaborar diagrama C4 — Contexto | TP2 | Done |
| 12 | Elaborar diagrama C4 — Container | TP2 | Done |
| 13 | Elaborar diagrama C4 — Componente | TP2 | Done |
| 14 | Definir contrato da API REST (endpoints, payloads) | TP2 | Done |
| 15 | Definir estrutura de pastas do projeto | TP2 | Done |
| 16 | Escrever documentação de arquitetura (TP2-arquitetura.md) | TP2 | Done |
| 17 | Atualizar GitHub Projects com tarefas do TP3 | TP2 | Done |

## Sprint TP3 — Sprint de Desenvolvimento

| # | Card | Sprint | Status |
|---|------|--------|--------|
| 18 | Inicializar módulo Go (`go mod init`) | TP3 | Done |
| 19 | Criar dataset CSV de indicadores climáticos | TP3 | Done |
| 20 | Implementar `internal/models` (Indicator, Series, Summary, TopEmitter) | TP3 | Done |
| 21 | Implementar `internal/data` — carregamento e parsing do CSV | TP3 | Done |
| 22 | Implementar `internal/data` — queries (Countries, Years, SeriesFor, Summary, Top) | TP3 | Done |
| 23 | Implementar handler `/api/health` | TP3 | Done |
| 24 | Implementar handler `/api/countries` e `/api/countries/{country}` | TP3 | Done |
| 25 | Implementar handler `/api/years` | TP3 | Done |
| 26 | Implementar handler `/api/summary` | TP3 | Done |
| 27 | Implementar handler `/api/emissions` (até 5 países, filtros de ano) | TP3 | Done |
| 28 | Implementar handler `/api/top` | TP3 | Done |
| 29 | Implementar `cmd/server/main.go` (bootstrap, logging, graceful shutdown) | TP3 | Done |
| 30 | Implementar frontend — estrutura HTML | TP3 | Done |
| 31 | Implementar frontend — folha de estilos (CSS) | TP3 | Done |
| 32 | Implementar frontend — lógica JS (consumo da API + filtros) | TP3 | Done |
| 33 | Integrar Chart.js — gráfico comparativo (linha) | TP3 | Done |
| 34 | Integrar Chart.js — gráfico de Top N (barra horizontal) | TP3 | Done |
| 35 | Smoke test end-to-end (curl + abertura no navegador) | TP3 | Done |
| 36 | Atualizar GitHub Projects com tarefas do TP4 | TP3 | Done |

## Sprint TP4 — Sprint de Desenvolvimento + Plano de Testes

| # | Card | Sprint | Status |
|---|------|--------|--------|
| 37 | Definir formato padrão do plano de testes (template TC) | TP4 | Done |
| 38 | Escrever 3 casos de teste para UC01 — Visualizar visão geral | TP4 | Done |
| 39 | Escrever 3 casos de teste para UC02 — Consultar série histórica | TP4 | Done |
| 40 | Escrever 3 casos de teste para UC03 — Comparar países | TP4 | Done |
| 41 | Escrever 3 casos de teste para UC04 — Filtrar por intervalo de anos | TP4 | Done |
| 42 | Escrever 3 casos de teste para UC05 — Listar Top N emissores | TP4 | Done |
| 43 | Adicionar testes unitários para `internal/data` (parsing, filtros) | TP4 | Done |
| 44 | Adicionar testes de integração HTTP (`httptest`) para handlers | TP4 | Done |
| 45 | Melhorias de UX no frontend (mensagens de erro, estados vazios) | TP4 | Done |
| 46 | Página `docs/TP4-plano-de-testes.md` consolidando os TCs | TP4 | Done |
| 47 | Atualizar GitHub Projects com tarefas do TP5 | TP4 | Done |

## Sprint TP5 — Execução dos Testes

| # | Card | Sprint | Status |
|---|------|--------|--------|
| 48 | Executar todos os TCs do plano e registrar resultados | TP5 | Done |
| 49 | Corrigir defeitos descobertos durante a execução dos TCs | TP5 | Done |
| 50 | Automatizar execução de testes unitários (`go test ./...`) em script | TP5 | Done |
| 51 | Atualizar `docs/TP4-plano-de-testes.md` com colunas de resultado | TP5 | Done |
| 52 | Atualizar GitHub Projects com tarefas do TP6 | TP5 | Done |

## Sprint TP6 — Entrega Final

| # | Card | Sprint | Status |
|---|------|--------|--------|
| 53 | Revisão geral da documentação | TP6 | Done |
| 54 | Polimento final do frontend | TP6 | Done |
| 55 | Verificar cobertura de todos os RF e RNF | TP6 | Done |
| 56 | Atualizar README com instruções finais | TP6 | Done |
| 57 | Tag de release `v1.0.0` no repositório | TP6 | Done |

---

## Resumo por status (entrega final — v1.0.0)

| Status | Quantidade |
|--------|------------|
| Done (TP1–TP6) | 57 |
| Todo | 0 |
| Backlog | 0 |
| **Total** | **57** |
