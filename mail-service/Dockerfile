FROM alpine:latest

RUN mkdir /app

COPY mailApp /app
COPY templates /templates

EXPOSE 5005

CMD ["/app/mailApp"]