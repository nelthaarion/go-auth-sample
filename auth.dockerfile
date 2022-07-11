FROM  alpine:latest 
RUN mkdir app 
COPY ./app app/
WORKDIR /app
CMD [ "./app" ]