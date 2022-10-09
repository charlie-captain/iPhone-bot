FROM alpine

WORKDIR /iPhone
COPY .. .

RUN chmod 777 ./iphoneBot

CMD [ "./iphoneBot" ]
