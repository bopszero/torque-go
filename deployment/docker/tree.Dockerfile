FROM golang:1-alpine

WORKDIR /opt/app/

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

COPY . /opt/app/

#RUN sh ./bin/build.sh
RUN python ./bin/build.py

EXPOSE 8080

CMD ["./bin/torque-go.bin", "serve_tree", "--host", "0.0.0.0", "--port", "8080"]
