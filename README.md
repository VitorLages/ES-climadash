# ClimaDash — Dashboard de Ação Climática

Dashboard de dados desenvolvido em Go que apresenta indicadores de emissões de gases de efeito estufa e adoção de energias renováveis por país, contribuindo para o **ODS 13 — Ação Contra a Mudança Global do Clima** (com interface com o **ODS 7 — Energia Acessível e Limpa**).

Projeto desenvolvido como trabalho prático da disciplina de **Engenharia de Software**.

## Sumário

- [Sobre o projeto](#sobre-o-projeto)
- [ODS abordado](#ods-abordado)
- [Stack](#stack)
- [Como executar](#como-executar)
- [Estrutura do repositório](#estrutura-do-repositório)
- [Documentação](#documentação)

## Sobre o projeto

Cidadãos, jornalistas, professores e estudantes têm dificuldade em acessar, de forma consolidada e comparável, dados sobre emissões de CO₂ e participação de energias renováveis por país. Ferramentas oficiais costumam exigir conhecimento técnico ou apresentar dados em formatos pouco amigáveis.

O **ClimaDash** centraliza um conjunto de dados públicos sobre o tema e fornece uma interface gráfica interativa, com filtros por país e intervalo de anos, permitindo que o usuário compreenda tendências e compare países de forma simples.

## ODS abordado

- **ODS 13 — Ação Contra a Mudança Global do Clima.** O dashboard fornece visibilidade sobre emissões de CO₂ totais e per capita.
- **ODS 7 — Energia Acessível e Limpa** (interface). Exibe a participação de fontes renováveis na matriz energética dos países.

## Stack

- **Backend:** Go 1.22 (biblioteca padrão `net/http`, `encoding/csv`, `encoding/json`).
- **Frontend:** HTML5, CSS3, JavaScript (vanilla) e [Chart.js](https://www.chartjs.org/) (via CDN).
- **Dados:** arquivo CSV embarcado no repositório, com indicadores compilados a partir de fontes públicas (Our World in Data / World Bank).
- **Build:** `go build`. Sem dependências externas.

## Como executar

Pré-requisitos: Go 1.22+ instalado.

```bash
git clone <url-do-repositorio>
cd ES
go run ./cmd/server
```

A aplicação ficará disponível em <http://localhost:8080>.

Variáveis de ambiente opcionais:

| Variável     | Padrão              | Descrição                                  |
|--------------|---------------------|--------------------------------------------|
| `PORT`       | `8080`              | Porta de escuta do servidor HTTP.          |
| `DATA_PATH`  | `data/emissions.csv`| Caminho para o arquivo CSV de indicadores. |

## Estrutura do repositório

```
.
├── cmd/server/            # ponto de entrada (main.go)
├── internal/
│   ├── data/              # carregamento e consulta dos dados CSV
│   ├── handlers/          # handlers HTTP / API REST
│   └── models/            # tipos de domínio
├── web/static/            # frontend (HTML, CSS, JS)
├── data/                  # dataset CSV
├── docs/                  # documentação (TPs, diagramas, board)
└── README.md
```

## Documentação

- [TP1 — Requisitos](docs/TP1-requisitos.md)
- [TP1 — Casos de uso](docs/TP1-casos-de-uso.md)
- [TP2 — Arquitetura (C4)](docs/TP2-arquitetura.md)
- [Quadro de tarefas (GitHub Projects)](docs/github-project.md)

## Licença

Uso acadêmico. Dados originais sob licença CC-BY 4.0 dos provedores citados.
