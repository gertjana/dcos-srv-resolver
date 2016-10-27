FROM scratch

EXPOSE 8000

COPY dist/dcos_sd /dcos_sd
CMD ["/dcos_sd"]