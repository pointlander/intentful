FROM scratch
COPY dbnary.db /dbnary.db
COPY wikipedia.db /wikipedia.db
COPY intentful /intentful
CMD ["/intentful"]
