FROM caddy

COPY frontend/Caddyfile /etc/caddy/Caddyfile
COPY frontend/ /srv