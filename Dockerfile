#
# Copyright (c) 2017 TFG Co <backend@tfgco.com>
# Author: TFG Co <backend@tfgco.com>
#
# Permission is hereby granted, free of charge, to any person obtaining a copy of
# this software and associated documentation files (the "Software"), to deal in
# the Software without restriction, including without limitation the rights to
# use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of
# the Software, and to permit persons to whom the Software is furnished to do so,
# subject to the following conditions:
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.
#
# THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
# IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS
# FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR
# COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER
# IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN
# CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
#

FROM golang:1.8.0-alpine

MAINTAINER TFG Co <backend@tfgco.com>

RUN apk update
RUN apk add make git g++ bash python wget

RUN wget https://github.com/Masterminds/glide/releases/download/v0.12.3/glide-v0.12.3-linux-amd64.tar.gz
RUN tar -zxvf glide-v0.12.3-linux-amd64.tar.gz
RUN chmod +x linux-amd64/glide && mv linux-amd64/glide /usr/local/bin/glide

RUN mkdir -p /go/src/github.com/topfreegames/offers
WORKDIR /go/src/github.com/topfreegames/offers

ADD glide.yaml /go/src/github.com/topfreegames/offers/glide.yaml
ADD glide.lock /go/src/github.com/topfreegames/offers/glide.lock
RUN glide install

ADD . /go/src/github.com/topfreegames/offers

RUN mkdir /app
RUN mv /go/src/github.com/topfreegames/offers/bin/offers /app/offers
RUN mv /go/src/github.com/topfreegames/offers/config /app/config
RUN rm -r /go/src/github.com/topfreegames/offers

WORKDIR /app

EXPOSE 8888
VOLUME /app/config

CMD /app/offers start-api -c /app/config/local.yaml
