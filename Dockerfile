FROM gcr.io/distroless/static-debian11:nonroot
ENTRYPOINT ["/baton-contentful"]
COPY baton-contentful /