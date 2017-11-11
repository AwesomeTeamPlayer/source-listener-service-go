FROM alpine:3.6

RUN mkdir /app
ADD ./app /app

RUN chmod a+x /app/app

CMD /app/app