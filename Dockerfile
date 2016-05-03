FROM alpine
ADD opentracing-example /home/opentracing-example
ENTRYPOINT /home/opentracing-example
