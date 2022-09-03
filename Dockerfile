## package #####################################################################

FROM ddollar/go:1.19 AS package

WORKDIR /src

COPY . .

RUN make binaries compress

## production ##################################################################

FROM ddollar/ubuntu:lts AS production

RUN apt-get update && apt-get -y --no-install-recommends install \
	curl \
	&& apt-get clean && rm -rf /var/lib/apt/lists/*

RUN curl -s https://download.docker.com/linux/static/stable/x86_64/docker-20.10.9.tgz | \
	tar -C /usr/bin --strip-components 1 -xz

ENV GOPATH=/go
ENV PATH=$PATH:/opt/bin

WORKDIR /

COPY --from=package /src/dist/keepd /opt/bin/

CMD ["/opt/bin/keepd"]
