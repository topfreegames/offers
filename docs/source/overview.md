Overview
========

Offers makes the implementation and management of offers in our games an predictable and scalable task. 

Used to sell special packs of the game's items. Offers helps doing A/B testing, remote configurations of offers, scheduling and monitoring.

## Features

* **Create Offer Template** - Create a template with information about a game's  offer: what it gives to the player, when it is enabled, how many times a player can see and buy it;
* **Get Available Offers** - Given a player and a game, returns the available offers for each placement in the UI;
* **New Relic Support** - Natively support new relic with segments in each API route for easy detection of bottlenecks;

## Architecture


Offers is composed of two mains models:
  - Models that define properties of game, offer and offerTemplate.
  - An API responsible for creation of games and templates, and retrive of offers.

## The Stack

Our code is in Golang, with: 

* Database - Postgres >= 9.5;

## Who's Using it

Well, right now, only us at TFG Co, are using it, but it would be great to get a community around the project. Hope to hear from you guys soon!

## How To Contribute?

Just the usual: Fork, Hack, Pull Request. Rinse and Repeat. Also don't forget to include tests and docs (we are very fond of both).
