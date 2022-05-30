# Image Builder
FROM telkomindonesia/debian:go-1.16 AS go-builder

LABEL maintainer="patricksangian@gmail.com"

# Set Working Directory
WORKDIR /usr/src/app

# Copy Source Code
COPY . ./

# Dependencies installation and binary file builder
RUN make install \
  && make build


# Final Image
# ---------------------------------------------------
FROM dimaskiddo/debian:base

# Set Working Directory
WORKDIR /usr/src/app

# Copy Anything The Application Needs
COPY --from=go-builder /tmp/app ./
COPY --from=go-builder /usr/src/app/template ./template

# Install dependencies binary
RUN apt update &&\
    apt install -y wkhtmltopdf &&\
    apt install -y xvfb &&\
    echo 'alias wkhtmltopdf="xvfb-run wkhtmltopdf"' >> ~/.bashrc

# Expose Application Port
EXPOSE 9000

# Run The Application
CMD ["./app"]