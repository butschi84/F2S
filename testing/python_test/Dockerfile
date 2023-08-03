FROM romanhuesler/f2sfizzlet AS f2sfizzlet

FROM python
WORKDIR /app
COPY --from=f2sfizzlet /app/f2sfizzlet ./f2sfizzlet
COPY ./executables/ /

CMD ["./f2sfizzlet"]