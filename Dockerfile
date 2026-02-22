FROM golang:1.26 AS compile
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /bin/app /app/.


FROM scratch AS execute
WORKDIR /bin
COPY ./.env /bin/.env
COPY --from=compile /bin/app /bin/app
ENTRYPOINT [ "/bin/app" ]
