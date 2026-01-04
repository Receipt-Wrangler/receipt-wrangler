FROM golang:1.24-trixie

# First install curl if not already installed
RUN apt-get update
RUN apt-get install -y curl

# Setup NodeSource repository for the LTS version
RUN curl -fsSL https://deb.nodesource.com/setup_lts.x | bash -

# Install Node.js
RUN apt-get install -y nodejs

# Install Java
RUN apt-get install -y default-jre

# Start in app dir
WORKDIR /app

# Copy repository into app directory
COPY . /app

# Setup Desktop
WORKDIR /app/desktop
RUN npm install -g @angular/cli
RUN npm install -g @openapitools/openapi-generator-cli

# Build desktop
RUN npm install
RUN npm run build

# Set up docs
WORKDIR /docs

# Clone docs source
RUN git clone https://github.com/Receipt-Wrangler/receipt-wrangler-doc.git
WORKDIR /docs/receipt-wrangler-doc

RUN npm install

# Setup API
WORKDIR /app/api

# Add local bin to path for python dependencies
ENV PATH="~/.local/bin:${PATH}"

# Set env
ENV ENV="dev"

# Set base path
ENV BASE_PATH="/app/api"

# Install tesseract dependencies
RUN bash set-up-dependencies.sh

# Install python dependencies for imap-client
RUN pip3 install -r imap-client/requirements.txt --break-system-packages

# Build api
RUN go build

# Set up data volume
RUN mkdir data
VOLUME /app/api/data

# Set up temp directory
RUN mkdir temp

# Set up sqlite volume
VOLUME /app/api/sqlite

# Add logs volume
RUN mkdir logs
VOLUME /app/api/logs

# Set pythonpath
ENV PYTHONPATH="/app/api/wranglervenv/lib/python3.11/site-packages"

# Expose http port
EXPOSE 80

# Setup SSH
RUN apt-get update
RUN apt-get install -y curl openssh-server

# Set root password for SSH access (change 'your_password' to something secure)
RUN echo 'root:development' | chpasswd
# Configure SSH to allow root login
RUN sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config
# Configure SSH to allow password authentication
RUN sed -i 's/#PasswordAuthentication yes/PasswordAuthentication yes/' /etc/ssh/sshd_config

# Create startup script to start SSH
RUN echo '#!/bin/bash\n\
/usr/sbin/sshd\n\
tail -f /dev/null' > /app/startup.sh && \
chmod +x /app/startup.sh

# Expose SSH
EXPOSE 22

# Add Go environment variables to root's bash profile
RUN echo 'export PATH=$PATH:/usr/local/go/bin:/go/bin' >> /root/.bashrc
RUN echo 'export GOPATH=/go' >> /root/.bashrc

WORKDIR /app

# Keep er runnin
CMD ["/app/startup.sh"]
