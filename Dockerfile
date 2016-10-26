FROM scratch

EXPOSE 8080

COPY dist/dcos_sd /dcos_sd
CMD ["/dcos_sd"]