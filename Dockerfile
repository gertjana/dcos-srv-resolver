FROM scratch

EXPOSE 8000

COPY dist/dcos_sd /dcos_sd
COPY dist/cmx /cmx
COPY dumb-init /dumb-init

RUN ["/cmx"]

ENTRYPOINT ["/dumb-init", "--"]
CMD ["/dcos_sd"]