FROM openshift/base-centos7 

RUN yum install -y golang && \
    yum clean all

ENV GOLANG_VERSION=1.9 \
    GOPATH=/go

ENV PATH=$GOPATH/bin:/usr/local/go/bin:$PATH

COPY go-wrapper /usr/local/bin/

RUN mkdir -p /go/src/github.com/dedgarsites/s3-browser
WORKDIR /go/src/github.com/dedgarsites/s3-browser

COPY . /go/src/github.com/dedgarsites/s3-browser
RUN go-wrapper download && go-wrapper install

EXPOSE 8080

USER 1001

CMD ["go-wrapper", "run"]
