FROM ubuntu:22.04
# FROM debian:bullseye-slim


ENV DEBIAN_FRONTEND=noninteractive

RUN   apt-get update && \ 
      apt-get install -y curl git build-essential libcap-dev libssl-dev zlib1g-dev && \
      rm -rf /var/lib/apt/lists/*;

# GCC
RUN   apt-get update && \ 
      apt-get install -y gcc g++ && \
      rm -rf /var/lib/apt/lists/*;

# GO - https://golang.org/dl
ENV GO_VERSION 1.20.3
RUN   curl -fSsL "https://go.dev/dl/go1.20.3.linux-amd64.tar.gz" -o /tmp/go-$GO_VERSION.tar.gz && \
      mkdir /usr/local/go && \
      tar -xf /tmp/go-$GO_VERSION.tar.gz -C /usr/local/go --strip-components=1 && \
      rm -rf /tmp/*;   

# Python2 - https://www.python.org/downloads
ENV PYTHON2_VERSION 2.7.18
RUN   curl -fSsL "https://www.python.org/ftp/python/$PYTHON2_VERSION/Python-$PYTHON2_VERSION.tar.xz" -o /tmp/python-$PYTHON2_VERSION.tar.xz && \
      mkdir /tmp/python-$PYTHON2_VERSION && \
      tar -xf /tmp/python-$PYTHON2_VERSION.tar.xz -C /tmp/python-$PYTHON2_VERSION --strip-components=1 && \
      rm /tmp/python-$PYTHON2_VERSION.tar.xz && \
      cd /tmp/python-$PYTHON2_VERSION && \
      ./configure \
      --prefix=/usr/local/python2 && \
      make -j$(nproc) && \
      make -j$(nproc) install && \
      rm -rf /tmp/*; 

# Python3 - https://www.python.org/downloads
ENV PYTHON3_VERSION 3.11.3
RUN   curl -fSsL "https://www.python.org/ftp/python/$PYTHON3_VERSION/Python-$PYTHON3_VERSION.tar.xz" -o /tmp/python-$PYTHON3_VERSION.tar.xz && \
      mkdir /tmp/python-$PYTHON3_VERSION && \
      tar -xf /tmp/python-$PYTHON3_VERSION.tar.xz -C /tmp/python-$PYTHON3_VERSION --strip-components=1 && \
      rm /tmp/python-$PYTHON3_VERSION.tar.xz && \
      cd /tmp/python-$PYTHON3_VERSION && \
      ./configure \
      --prefix=/usr/local/python3 && \
      make -j$(nproc) && \
      make -j$(nproc) install && \
      rm -rf /tmp/*; 

# Rust - https://www.rust-lang.org
ENV RUST_VERSION 1.69.0
RUN   curl -fSsL "https://static.rust-lang.org/dist/rust-$RUST_VERSION-x86_64-unknown-linux-gnu.tar.gz" -o /tmp/rust-$RUST_VERSION.tar.gz && \
      mkdir /tmp/rust-$RUST_VERSION && \
      tar -xf /tmp/rust-$RUST_VERSION.tar.gz -C /tmp/rust-$RUST_VERSION --strip-components=1 && \
      rm /tmp/rust-$RUST_VERSION.tar.gz && \
      cd /tmp/rust-$RUST_VERSION && \
      ./install.sh \
      --prefix=/usr/local/rust \
      --components=rustc,rust-std-x86_64-unknown-linux-gnu && \
      rm -rf /tmp/*; 

# Node - https://nodejs.org/en/download
ENV NODE_VERSION 18.16.0
RUN   curl -fSsL "https://nodejs.org/dist/v$NODE_VERSION/node-v$NODE_VERSION-linux-x64.tar.gz" -o /tmp/node-$NODE_VERSION.tar.gz && \
      mkdir /usr/local/node && \
      tar -xf /tmp/node-$NODE_VERSION.tar.gz -C /usr/local/node --strip-components=1 && \
      rm -rf /tmp/*;   

ENV TYPESCRIPT_VERSION 5.0.4      
RUN PATH="/usr/local/node/bin:${PATH}" npm install -g typescript@$TYPESCRIPT_VERSION

# JDK 20.0.1 - https://jdk.java.net
RUN   curl -fSsL "https://download.java.net/java/GA/jdk20.0.1/b4887098932d415489976708ad6d1a4b/9/GPL/openjdk-20.0.1_linux-x64_bin.tar.gz" -o /tmp/openjdk20.0.1.tar.gz && \
      mkdir /usr/local/openjdk && \
      tar -xf /tmp/openjdk20.0.1.tar.gz -C /usr/local/openjdk --strip-components=1 && \
      rm /tmp/openjdk20.0.1.tar.gz 

# CSharp (Mono)
RUN   apt-get update && \ 
      apt-get install -y mono-devel && \
      rm -rf /var/lib/apt/lists/*;

# PHP - https://www.php.net/downloads
ENV PHP_VERSION  8.2.6
RUN   apt-get update && \
      apt-get install -y --no-install-recommends bison re2c autoconf libxml2-dev libsqlite3-dev && \
      rm -rf /var/lib/apt/lists/* && \
      curl -fSsL "https://codeload.github.com/php/php-src/tar.gz/php-$PHP_VERSION" -o /tmp/php-$PHP_VERSION.tar.gz && \
      mkdir /tmp/php-$PHP_VERSION && \
      tar -xf /tmp/php-$PHP_VERSION.tar.gz -C /tmp/php-$PHP_VERSION --strip-components=1 && \
      rm /tmp/php-$PHP_VERSION.tar.gz && \
      cd /tmp/php-$PHP_VERSION && \
      ./buildconf --force && \
      ./configure \
      --prefix=/usr/local/php && \
      make -j$(nproc) && \
      make -j$(nproc) install && \
      rm -rf /tmp/*; 

# Swift - https://swift.org/download
ENV SWIFT_VERSION 5.8
RUN   apt-get update && \
      apt-get install -y --no-install-recommends libncurses5 libsqlite3-dev libc6 && \
      rm -rf /var/lib/apt/lists/* && \
      curl -fSsL "https://download.swift.org/swift-$SWIFT_VERSION-release/ubuntu2204/swift-$SWIFT_VERSION-RELEASE/swift-$SWIFT_VERSION-RELEASE-ubuntu22.04.tar.gz" -o /tmp/swift-$SWIFT_VERSION.tar.gz && \
      mkdir /usr/local/swift && \
      tar -xf /tmp/swift-$SWIFT_VERSION.tar.gz -C /usr/local/swift --strip-components=2 && \
      rm -rf /tmp/*; 

# Kotlin - https://kotlinlang.org
ENV KOTLIN_VERSION 1.8.21
RUN   apt-get update && \
      apt-get install -y unzip && \
      rm -rf /var/lib/apt/lists/* && \
      curl -fSsL "https://github.com/JetBrains/kotlin/releases/download/v$KOTLIN_VERSION/kotlin-compiler-$KOTLIN_VERSION.zip" -o /tmp/kotlin-$KOTLIN_VERSION.zip && \
      unzip -d /usr/local/kotlin /tmp/kotlin-$KOTLIN_VERSION.zip && \
      mv /usr/local/kotlin/kotlinc/* /usr/local/kotlin/ && \
      rm -rf /usr/local/kotlin/kotlinc && \
      rm -rf /tmp/*; 


# Ruby (3.2.2) - https://www.ruby-lang.org/en/downloads
ENV RUBY_VERSION 3.2.2
RUN   apt-get update && \
      apt-get install -y libyaml-dev && \
      rm -rf /var/lib/apt/lists/* && \
      curl -fSsL "https://cache.ruby-lang.org/pub/ruby/3.2/ruby-$RUBY_VERSION.tar.gz" -o /tmp/ruby-$RUBY_VERSION.tar.gz && \
      mkdir /tmp/ruby-$RUBY_VERSION && \
      tar -xf /tmp/ruby-$RUBY_VERSION.tar.gz -C /tmp/ruby-$RUBY_VERSION --strip-components=1 && \
      rm /tmp/ruby-$RUBY_VERSION.tar.gz && \
      cd /tmp/ruby-$RUBY_VERSION && \
      ./configure \
      --disable-install-doc \
      --prefix=/usr/local/ruby && \
      make -j$(nproc) && \
      make -j$(nproc) install && \
      rm -rf /tmp/*;  


# Scala (3.2.2) - https://scala-lang.org
RUN   curl -fSsL "https://github.com/lampepfl/dotty/releases/download/3.2.2/scala3-3.2.2.tar.gz" -o /tmp/scala-3.2.2.tgz && \
      mkdir /usr/local/scala && \
      tar -xf /tmp/scala-3.2.2.tgz -C /usr/local/scala --strip-components=1 && \
      rm -rf /tmp/*; 


# Elixir - https://github.com/elixir-lang/elixir/releases
ENV ELIXIR_VERSION 1.14.4
RUN   apt-get update && \
      apt-get install -y --no-install-recommends unzip locales && \
      locale-gen en_US.UTF-8 && \
      rm -rf /var/lib/apt/lists/* && \
      curl -fSsL "https://github.com/elixir-lang/elixir/releases/download/v$ELIXIR_VERSION/elixir-otp-23.zip" -o /tmp/elixir-$ELIXIR_VERSION.zip && \
      unzip -d /usr/local/elixir /tmp/elixir-$ELIXIR_VERSION.zip && \
      rm -rf /tmp/*;  


# ERLANG - https://github.com/erlang/otp/releases
ENV ERLANG_VERSION 25.3.1
RUN   apt-get update && \
      apt-get install -y libncurses5-dev && \
      rm -rf /var/lib/apt/lists/* && \
      curl -fSsL "https://github.com/erlang/otp/releases/download/OTP-$ERLANG_VERSION/otp_src_$ERLANG_VERSION.tar.gz" -o /tmp/erlang-$ERLANG_VERSION.tar.gz && \
      mkdir /tmp/erlang-$ERLANG_VERSION && \
      tar -xf /tmp/erlang-$ERLANG_VERSION.tar.gz -C /tmp/erlang-$ERLANG_VERSION --strip-components=1 && \
      rm /tmp/erlang-$ERLANG_VERSION.tar.gz && \
      cd /tmp/erlang-$ERLANG_VERSION && \
      ./otp_build autoconf && \
      ./configure \
      --prefix=/usr/local/erlang && \
      make -j$(nproc) && \
      make -j$(nproc) install && \
      rm -rf /tmp/*;     


# Haskell - https://www.haskell.org/ghc/download.html
ENV HASKELL_VERSION 9.6.1
RUN   apt-get update && \
      apt-get install -y --no-install-recommends libgmp-dev libtinfo5 && \
      rm -rf /var/lib/apt/lists/* && \
      curl -fSsL "https://downloads.haskell.org/~ghc/$HASKELL_VERSION/ghc-$HASKELL_VERSION-x86_64-deb11-linux.tar.xz" -o /tmp/ghc-$HASKELL_VERSION.tar.xz && \
      mkdir /tmp/ghc-$HASKELL_VERSION && \
      tar -xf /tmp/ghc-$HASKELL_VERSION.tar.xz -C /tmp/ghc-$HASKELL_VERSION --strip-components=1 && \
      rm /tmp/ghc-$HASKELL_VERSION.tar.xz && \
      cd /tmp/ghc-$HASKELL_VERSION && \
      ./configure \
      --prefix=/usr/local/ghc && \
      make -j$(nproc) install && \
      rm -rf /tmp/*; 


# Zig - https://ziglang.org/download/
RUN   curl -fSsL "https://ziglang.org/builds/zig-linux-x86_64-0.11.0-dev.3132+465272921.tar.xz" -o /tmp/zig.tar.gz && \
      mkdir /usr/local/zig && \
      tar -xf /tmp/zig.tar.gz -C /usr/local/zig --strip-components=1 && \
      rm -rf /tmp/*;   

WORKDIR /app
RUN git clone https://github.com/ioi/isolate.git

WORKDIR /app/isolate
RUN make isolate
RUN make install