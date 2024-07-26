FROM golang:alpine
ADD . /work
ENV CGO_ENABLED=0
RUN go build -C /work -o app

FROM scratch
COPY --from=0 /work/app /app
ENTRYPOINT ["/app"]
