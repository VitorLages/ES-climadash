#!/usr/bin/env bash
#
# run-tests.sh — automatiza a execução dos testes do ClimaDash (TP5, issue #50).
#
# Uso:
#   ./scripts/run-tests.sh            # roda go vet + go test ./...
#   ./scripts/run-tests.sh -v         # modo verboso (-v no go test)
#   ./scripts/run-tests.sh -c         # gera relatório de cobertura (coverage.out)
#
set -euo pipefail

# Posiciona na raiz do projeto, independentemente de onde o script é chamado.
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
cd "$SCRIPT_DIR/.."

VERBOSE=""
COVER=""
for arg in "$@"; do
  case "$arg" in
    -v|--verbose) VERBOSE="-v" ;;
    -c|--cover)   COVER="1" ;;
    -h|--help)
      grep '^#' "$0" | sed 's/^# \{0,1\}//'
      exit 0
      ;;
    *) echo "argumento desconhecido: $arg" >&2; exit 2 ;;
  esac
done

echo "==> go vet ./..."
go vet ./...

echo "==> go test ./..."
if [[ -n "$COVER" ]]; then
  go test $VERBOSE -coverprofile=coverage.out ./...
  echo "==> cobertura por pacote:"
  go tool cover -func=coverage.out | tail -n 1
  echo "    (relatório completo em coverage.out — veja com: go tool cover -html=coverage.out)"
else
  go test $VERBOSE ./...
fi

echo "==> todos os testes passaram ✅"
