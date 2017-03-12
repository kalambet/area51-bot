# teledisq

[Telegram](https://telegram.org/) bot that reacts to [Discourse](https://discourse.org) notifications

This is work in progress but right now it serves its main purpose: it post chat messages to subscribed channels about updates from forums powered by Discourse. Right now only `post_created` and `topic_created` event types are supported. (see more: https://meta.discourse.org/t/setting-up-webhooks/49045)

## Instalaltion

All the configuration can be done only in `app.yaml` file before deploy.

- `DISCOURSE_WEBHOOK`: webhook suffix for Discoures. Discourse will call it every time event occures. (ex. https://project-id.appspot.com/hook/DISCOURSE_WEBHOOK/?type=discourse) `type` is the value introduced for general purpose but only `discourse` type right now is supported
- `TELEGRAM_WEBHOOK`: webhook suffix for Telegram. Telegram will call it on every message dedicated to this bot. Bot suppose to work in `privacy_mode` (ex. https://project-id.appspot.com/telegram/TELEGRAM_WEBHOOK/)
- `TELEGRAM_SECRET`: Telegram secret BotFather gave you
- `TELEGRAM_BOT_USERNAME`: username for you bot, used to do unsubscribtion when bot is kicked off from chat 
- `BOT_COMMAND`: Bot command is THE command you are goig to comunicate with your bot, it strats with `/` (ex. `/mybot`)
- `BOT_SUBSCRIPTION_COMMAND`: Argument that you need to pass to your command in order to subscribe it to your Discourse forum (ex. `subscribe`)

After deployment to AppEninge you need to put your https://project-id.appspot.com/hook/DISCOURSE_WEBHOOK/?type=discourse to the Discourse settins as it described in this topic: https://meta.discourse.org/t/setting-up-webhooks/49045

And set https://project-id.appspot.com/telegram/TELEGRAM_WEBHOOK/ webhook to your bot:

```
curl -X "POST" "https://api.telegram.org/TELEGRAM_SECRET/setWebhook" \
     -H "Content-Type: application/json; charset=utf-8" \
     -d $'{"url": "https://project-id.appspot.com/telegram/TELEGRAM_WEBHOOK/"}'
```

Afer that just add bot to desired channel and call `BOT_COMMAND BOT_SUBSCRIPTION_COMMAND` (ex. `/mybot subscribe`)

## Some notes

This bot is supposed to be run on Google App Engine Standard environment mostly beacsue of its free limits https://cloud.google.com/free/docs/always-free-usage-limits so it's possible to have such function for free.

There is no flawor of goroutine or eny parallel execution in code, this was done mostly because AppEngine does not support any parallel execution right now https://cloud.google.com/appengine/docs/standard/go/runtime:


> The Go runtime environment for App Engine provides full support for goroutines, but not for parallel execution: goroutines are scheduled onto a single operating system thread. This single-thread restriction may be lifted in future versions. Multiple requests may be handled concurrently by a given instance; that means that if one request is, say, waiting for a datastore API call, another request may be processed by the same instance.
