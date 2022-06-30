# docker build -t hc-echo-example .
# docker run --name hc-echo-example --rm -p 5678:5678 hc-echo-example -text='my name is earl!'
# open http://localhost:5678
FROM hashicorp/http-echo

COPY --from=gregthebunny/go-healthcheck /bin/healthcheck /bin/healthcheck

ENV PROBE='{{ .Assert.HTTPBodyContains .HTTP.Handler "GET" "http://127.0.0.1:5678" nil "my name is earl!" }}'
HEALTHCHECK --start-period=10s --retries=3 --timeout=10s --interval=10s \
    CMD ["/bin/healthcheck", "probe", "$PROBE"]
