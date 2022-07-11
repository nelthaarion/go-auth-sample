FROM  alpine:latest 

RUN mkdir /app 

COPY ./app /app



CMD [ "./app/app" ]