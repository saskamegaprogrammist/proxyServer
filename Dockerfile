FROM golang:1.13 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN ls
RUN cd repeater && go build .
RUN cd proxy && go build .

FROM ubuntu:18.04 AS release

ENV PGVER 10
RUN apt -y update && apt install -y postgresql-$PGVER

USER postgres


RUN /etc/init.d/postgresql start &&\
	psql --command "CREATE USER docker WITH SUPERUSER PASSWORD 'docker';" &&\
	createdb -O docker docker &&\
	psql --command "GRANT ALL ON DATABASE docker TO docker;" &&\
    /etc/init.d/postgresql stop

ENV POSTGRES_DSN=postgres://docker:docker@localhost/docker

# Adjust PostgreSQL configuration so that remote connections to the
# database are possible.
RUN echo "host all  all    0.0.0.0/0  md5" >> /etc/postgresql/$PGVER/main/pg_hba.conf
RUN echo "listen_addresses='*'" >> /etc/postgresql/$PGVER/main/postgresql.conf

RUN echo "synchronous_commit='off'" >> /etc/postgresql/$PGVER/main/postgresql.conf
RUN echo "fsync = 'off'" >> /etc/postgresql/$PGVER/main/postgresql.conf


EXPOSE 5432

VOLUME ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]

USER root

EXPOSE 5000
EXPOSE 5001

COPY --from=build /opt/app/proxy /usr/bin/
COPY --from=build /opt/app/repeater /usr/bin/
COPY --from=build /opt/app/rootCert.cert /rootCert.cert
COPY --from=build  /opt/app/rootKey.pem /rootKey.pem

RUN apt-get update \
  && apt-get install -y curl

CMD service postgresql start && proxy && repeater

