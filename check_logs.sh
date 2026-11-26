#!/bin/bash
echo "--- 1. flood de commandes (200 look d'un coup)"
(printf 'CONNECT flood\n'; sleep 1; for i in $(seq 1 200); do echo look; done; sleep 1) | timeout 6 nc localhost 8080 > /dev/null

echo "--- 2. spam de chat (5 messages d'un coup)"
(printf 'CONNECT spam\n'; sleep 1; for i in 1 2 3 4 5; do echo "chat global hi"; done; sleep 1) | timeout 6 nc localhost 8080 > /dev/null

echo "--- 3. connexions rapides (15 en même temps)"
for i in $(seq 1 15); do (printf 'CONNECT rapid%s\n' $i; sleep 2) | timeout 5 nc localhost 8080 > /dev/null & done
wait

echo "--- 4. ligne trop longue (3000 caractères)"
(printf 'CONNECT big\n'; sleep 1; head -c 3000 /dev/zero | tr '\0' a; echo; sleep 1) | timeout 6 nc localhost 8080 > /dev/null

echo
echo "=== Résultat dans server.log ==="
for m in "abuse: command flooding" "abuse: chat spam" "abuse: rapid connections" "abuse: line too long"; do
    printf '%-28s %s\n' "$m" "$(grep -c "\"msg\":\"$m\"" server.log)"
done
