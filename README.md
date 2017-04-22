## Twitch interesting stream fragments detection service

This service intended for monitoring and detecting spikes in Twitch stream chat and call some endpoint for uploading stream fragment when it occured. Also, this service can store information about this spikes and uploaded fragments at some storage (now only MongoDB is available). Spikes detecion is based in average number at floating window.

### Usage

This service can be used with this [project](https://github.com/dimorinny/twitch-interesting-fragments-frontend), that allows us to monitoring twitch streams and upload this fragments to vk group.

### Params

All params is passed as env params.

#### Twitch params:

* **HOST** - Twitch IRC endpoint (Default: `irc.chat.twitch.tv`)
* **NICKNAME** - Twitch Nickname
* **OAUTH** - Oauth token (You can quickly get oauth token for your account with this [helpful page](http://twitchapps.com/tmi/))
* **CHANNEL** - Channel name

#### Uploader params:

* **UPLOADER\_HOST** - Uploader endpoint
* **UPLOADER\_PORT** - Uploader port

#### Storage params:

* **STORAGE\_TYPE** - Storage type for saving fragments metadata. Available values: no, mongo. (Default: `no`)
* **STORAGE\_HOST** - Storage (like MongoDB) endpoint.

#### Detection params:

* **MESSAGES\_BUFFER\_TIME** - Chat messages analyzed by groups. This param determine group size (Default: `25`)
* **WINDOW\_SIZE** - Size of floating window for calculating average value. (Default: `10`)
* **RECORD\_DELAY** - Delay after detection spike and executing upload endpoint. (Default: `20`)
* **SPIKE\_RATE** - Spike detection coefficient. Spike detected like this `currentValue > average * spikeRate` (Default: `4`)
* **SMOOTH\_RATE** - Сompensate micro spike coefficient (Default: `2`)

For more information about detection mechanism you can read this [code](https://github.com/dimorinny/twitch-interesting-fragments/blob/master/detection/detection.go#L21-L62).
