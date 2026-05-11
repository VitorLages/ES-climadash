# TP1 — Definição do Problema e Requisitos

## 1. ODS abordado

**ODS 13 — Ação Contra a Mudança Global do Clima.**
Interface secundária com o **ODS 7 — Energia Acessível e Limpa**, dado que a transição para fontes renováveis é um vetor fundamental para a redução de emissões.

## 2. Problema

Apesar de existirem diversas fontes oficiais de dados sobre o clima (UNFCCC, Our World in Data, World Bank, IPCC), as informações chegam ao público em formatos pouco amigáveis: planilhas extensas, relatórios em PDF ou portais com interfaces complexas. Isso dificulta que **estudantes, jornalistas, educadores e cidadãos** comparem rapidamente o desempenho ambiental de diferentes países e percebam tendências históricas.

A ausência de visualizações simples e comparáveis sobre emissões de CO₂ e participação de fontes renováveis na matriz energética reduz o engajamento público com o tema e empobrece o debate sobre políticas climáticas.

## 3. Tipo de solução escolhido

**Dashboard de dados** acessível via navegador, composto por:

- **Backend em Go**, responsável por carregar o conjunto de dados, oferecer uma API REST/JSON e servir o frontend estático.
- **Frontend em HTML/CSS/JS** com biblioteca de gráficos (Chart.js), consumindo a API.

### Justificativa da escolha

| Critério                  | Justificativa                                                                                                         |
|---------------------------|-----------------------------------------------------------------------------------------------------------------------|
| Público-alvo amplo        | Um dashboard web é acessível a qualquer pessoa com navegador — não exige instalação.                                  |
| Comparação visual         | O problema é, essencialmente, de visualização e comparação de séries numéricas — caso de uso clássico de dashboard.   |
| Backend obrigatório       | O escopo do trabalho proíbe soluções apenas-frontend. O backend Go executa filtragem, agregação e expõe API JSON.     |
| Performance               | Go tem ótima performance e baixa pegada de memória, adequado a um serviço HTTP simples sem necessidade de orquestração.|
| Simplicidade operacional  | Binário único, sem dependências externas além do runtime. Facilita o build, o deploy e a avaliação.                    |

Soluções alternativas (chatbot, app mobile, sistema transacional) foram descartadas porque o problema é primariamente **informacional/exploratório**, e um dashboard responde diretamente à necessidade de visualização e comparação.

## 4. Requisitos Funcionais

| ID    | Requisito                                                                                          | Prioridade |
|-------|----------------------------------------------------------------------------------------------------|------------|
| RF01  | O sistema deve carregar dados de emissões e energia renovável a partir de um arquivo CSV.          | Alta       |
| RF02  | O sistema deve listar todos os países disponíveis no dataset, via API.                             | Alta       |
| RF03  | O sistema deve listar os anos disponíveis no dataset, via API.                                     | Alta       |
| RF04  | O sistema deve permitir consultar a série histórica de um país específico.                         | Alta       |
| RF05  | O sistema deve permitir filtrar dados por intervalo de anos (`from`, `to`).                        | Alta       |
| RF06  | O sistema deve permitir comparar até 5 países simultaneamente em um gráfico.                       | Alta       |
| RF07  | O sistema deve apresentar visão geral com totais agregados (emissões globais e renováveis médias). | Média      |
| RF08  | O sistema deve exibir gráfico de evolução temporal de emissões de CO₂.                              | Alta       |
| RF09  | O sistema deve exibir gráfico de participação de renováveis na matriz energética.                  | Alta       |
| RF10  | O sistema deve exibir os "Top N" maiores emissores no último ano disponível.                       | Média      |
| RF11  | O sistema deve atualizar os gráficos dinamicamente após o usuário aplicar filtros.                 | Alta       |

## 5. Requisitos Não Funcionais

| ID     | Tipo            | Requisito                                                                                              |
|--------|-----------------|--------------------------------------------------------------------------------------------------------|
| RNF01  | Performance     | A API deve responder em menos de 300 ms para o dataset embarcado (assumindo execução local).           |
| RNF02  | Usabilidade     | O dashboard deve ser legível e usável em telas de pelo menos 1024×768.                                 |
| RNF03  | Manutenibilidade| Código organizado em pacotes `cmd/`, `internal/handlers`, `internal/data`, `internal/models`.          |
| RNF04  | Portabilidade   | Aplicação deve compilar e executar em Linux, macOS e Windows com Go 1.22+ instalado.                   |
| RNF05  | Disponibilidade | O servidor deve iniciar e responder requisições em menos de 2 s após `go run`.                         |
| RNF06  | Confiabilidade  | Erros de leitura de dataset devem ser reportados em log e não derrubar o servidor após o boot.         |
| RNF07  | Segurança       | A API expõe apenas operações de leitura (GET). Não há upload, autenticação ou estado mutável.          |
| RNF08  | Open-source     | Repositório público, documentação em Markdown, dependências apenas da biblioteca padrão Go + Chart.js. |
