FROM alpine:3.12.3

COPY ./src/build/monzo-ynab /usr/local/bin/

CMD ["monzo-ynab", "run", "--port", "80"]