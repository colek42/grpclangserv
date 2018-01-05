FROM scratch

ENV LS_GRPC_PORT 4534
EXPOSE $LS_GRPC_PORT

COPY grpclangserve /
CMD ["/grpclangserve"]
