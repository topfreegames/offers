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

FROM golang:1.8.1-alpine

MAINTAINER TFG Co <backend@tfgco.com>

RUN mkdir /app
ADD ./bin/offers-linux-x86_64 /app/offers
ADD ./config /app/config

WORKDIR /app
EXPOSE 8888

ENV OFFERS_NEWRELIC_APP offers
ENV OFFERS_POSTGRES_DBNAME offers
ENV OFFERS_POSTGRES_HOST localhost
ENV OFFERS_POSTGRES_PASSWORD ""
ENV OFFERS_POSTGRES_PORT 8585
ENV OFFERS_POSTGRES_USER offers

CMD /app/offers start -v2 -c /app/config/local.yaml
