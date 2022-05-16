FROM alpine:3.14

RUN apk add --no-cache tzdata
ADD build/ /srv/
# If need custom config, you can use `-v path:/config` mount volume to docker
ADD config/config.yml /config/config.yml

ENV env=prod
ENV TZ=Asia/Shanghai
# RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone
ENTRYPOINT ["/srv/blbconf"]