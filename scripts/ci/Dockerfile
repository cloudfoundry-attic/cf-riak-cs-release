FROM ubuntu:trusty
MAINTAINER https://github.com/cloudfoundry/cf-riak-cs-release

RUN \
      apt-get update && \
      apt-get -y install --fix-missing \
            build-essential \
            curl \
            git \
            libreadline6 \
            libreadline6-dev \
            wget \
            unzip \
      && \
      apt-get clean

# Install ruby-install
RUN curl https://codeload.github.com/postmodern/ruby-install/tar.gz/v0.5.0 | tar xvz -C /tmp/ && \
          cd /tmp/ruby-install-0.5.0 && \
          make install

ENV DOCKER_RUBY_VERSION 2.2.3

# Install Ruby
RUN ruby-install ruby ${DOCKER_RUBY_VERSION} -- --disable-install-rdoc

# Add ruby to PATH
ENV PATH $PATH:/home/root/.gem/ruby/${DOCKER_RUBY_VERSION}/bin:/opt/rubies/ruby-${DOCKER_RUBY_VERSION}/lib/ruby/gems/${DOCKER_RUBY_VERSION}/bin:/opt/rubies/ruby-${DOCKER_RUBY_VERSION}/bin

# Install gems
RUN /opt/rubies/ruby-${DOCKER_RUBY_VERSION}/bin/gem install bosh_cli --no-rdoc --no-ri && \
          /opt/rubies/ruby-${DOCKER_RUBY_VERSION}/bin/gem install bundler --no-rdoc --no-ri

# Install cf cli
RUN curl http://go-cli.s3-website-us-east-1.amazonaws.com/releases/v6.10.0/cf-linux-amd64.tgz | tar xvz -C /bin

# Install go
RUN curl https://storage.googleapis.com/golang/go1.4.2.linux-amd64.tar.gz | tar xvz -C /usr/local/

# Add golang environment variables
ENV GOPATH /usr/local/go
ENV PATH $PATH:/usr/local/go/bin

# Install spiff
RUN wget -O /tmp/spiff.zip https://github.com/cloudfoundry-incubator/spiff/releases/download/v1.0.7/spiff_linux_amd64.zip
RUN unzip /tmp/spiff.zip -d /tmp/ && mv /tmp/spiff $GOPATH/bin/

# Install jq 1.5
RUN wget -O /usr/local/bin/jq https://github.com/stedolan/jq/releases/download/jq-1.5/jq-linux64 && \
  chmod +x /usr/local/bin/jq

# Add workspace directory
RUN mkdir -p /root/workspace && ln -s /root/workspace /workspace

