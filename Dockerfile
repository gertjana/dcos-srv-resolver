FROM scratch

EXPOSE 8000

COPY dist/yp /yp
COPY dist/cmx /cmx
COPY dumb-init /dumb-init

RUN ["/cmx"]

ENTRYPOINT ["/dumb-init", "--"]
CMD ["/yp"]
