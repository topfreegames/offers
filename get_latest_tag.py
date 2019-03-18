#!/usr/bin/env python3
# -*- coding: utf-8 -*-
# Offers
# https://github.com/topfreegames/offers
# Licensed under the MIT license:
# http://www.opensource.org/licenses/mit-license
# Copyright © 2016 Top Free Games <backend@tfgco.com>

import datetime
import json
import urllib.request


def main():
    url = "https://registry.hub.docker.com/v2/repositories/tfgco/offers/tags?page_size=100"
    with urllib.request.urlopen(url) as response:
        res = json.loads(response.read().decode('utf-8'))
        last_tag = get_last_tag(res['results'])
        print(last_tag)


def get_tag_value(tag):
    # format should be $BUILD_NUMBER-$GIT_TAG-$GIT_COMMIT
    if len(tag['name'].split('-')) < 3:
        return {'tag': tag['name'], 'value': 0}
    d = datetime.datetime.strptime(tag['last_updated'], '%Y-%m-%dT%H:%M:%S.%fZ').timestamp()
    return {'tag': tag['name'], 'value': d}


def get_last_tag(tags):
    return max(
        [get_tag_value(t) for t in tags],
        key=lambda t: t['value']
    )['tag']


if __name__ == "__main__":
    main()
