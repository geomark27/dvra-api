#!/usr/bin/env bash
# Fusiona en develop el PR de la rama actual, siempre con merge commit (conserva los
# commits convencionales de la rama). Necesita gh autenticado.
#   make merge [yes=1]
set -euo pipefail

remote=origin
base=develop

fail() { echo "✗ $*" >&2; exit 1; }
ok() { echo "✓ $*"; }

command -v gh >/dev/null && gh auth status >/dev/null 2>&1 || fail "hace falta gh autenticado: gh auth login"
branch=$(git branch --show-current)
case "$branch" in main|"$base"|"") fail "no hay PR que fusionar desde '$branch': trabaja en una rama (make start)";; esac

git fetch --quiet "$remote" "$branch" 2>/dev/null || fail "'$branch' no está en $remote: corre make push"
[ "$(git rev-list --count "$remote/$branch..$branch")" -eq 0 ] || fail "'$branch' tiene commits sin subir: corre make push"

pr=$(gh pr list --head "$branch" --state open --json number --jq '.[0].number // empty')
[ -n "$pr" ] || fail "no hay un PR abierto de '$branch': corre make pr"
target=$(gh pr view "$pr" --json baseRefName --jq .baseRefName)
[ "$target" = "$base" ] || fail "el PR #$pr va hacia '$target', no hacia $base"
title=$(gh pr view "$pr" --json title --jq .title)

if [ "${MERGE_YES:-}" != 1 ]; then
	read -r -p "¿Fusionar el PR #$pr «$title» en $base con merge commit? [s/N]: " answer
	[[ "$answer" =~ ^[sSyY] ]] || fail "sin fusionar: el PR #$pr queda abierto"
fi

# Sin --delete-branch: gh dejaría la copia local en la rama por defecto;
# make finish vuelve a develop y borra la rama ya fusionada.
gh pr merge "$pr" --merge
ok "PR #$pr fusionado en $base con merge commit"
