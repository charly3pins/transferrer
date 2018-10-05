FROM golang:1.10

RUN mkdir -p /etc/transferrer

COPY . /etc/transferrer/.

WORKDIR /etc/transferrer/

RUN make deps build

CMD [ "/etc/transferrer/cmd/server/transferrer" ]