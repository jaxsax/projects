FROM ubuntu:19.10

RUN set -x \
        && addgroup --system --gid 101 nginx \
        && adduser --system --disabled-login --ingroup nginx --no-create-home --home /nonexistent --gecos "nginx user" --shell /bin/false --uid 101 nginx \
        && apt-get update \
        && apt-get install --no-install-recommends --no-install-suggests -y \
            ca-certificates \
            nginx

RUN  set -x \
        && ln -sf /dev/stdout /var/log/nginx/access.log \
        && ln -sf /dev/stderr /var/log/nginx/error.log \
        && mkdir /app && chown -R nginx:nginx /app

EXPOSE 80

STOPSIGNAL SIGTERM

VOLUME /app

CMD ["nginx", "-g", "daemon off;"]
