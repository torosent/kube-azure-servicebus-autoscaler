FROM alpine:3.4

RUN  apk add --no-cache --update ca-certificates

COPY kube-azure-servicebus-autoscaler /

CMD ["/kube-azure-servicebus-autoscaler"]
