FROM alpine:3.12 AS rie

RUN apk add --no-cache curl && \
  cd / && \
  curl -sSL -O https://github.com/aws/aws-lambda-runtime-interface-emulator/releases/download/v1.0/aws-lambda-rie && \
  chmod +x /aws-lambda-rie

FROM lambda-s3-pgn

COPY --from=rie /aws-lambda-rie /usr/local/bin/

ENTRYPOINT ["aws-lambda-rie", "/app/main"]
