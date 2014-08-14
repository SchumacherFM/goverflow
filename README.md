# Goverflow: a twitter bot to power @goverflow

Searches the StackExchange API and tweets new questions.

You need a twitter account and you need to create a new App [https://apps.twitter.com/](https://apps.twitter.com/).

You do NOT need a StackExchange account. The search API can be used up to 300 request per 24h. That means to check every
288 seconds for new questions.

### What can I search on the StackExchange API?

Please see [https://api.stackexchange.com/docs/search](https://api.stackexchange.com/docs/search)

### config.json file

```json
{
	"host": "https://api.stackexchange.com",
	"apiVersion": "2.2",
	"searchParams":	"order=desc&sort=creation&tagged=go&site=stackoverflow",
	"twitterConfigFile": "config.twitter.json",
	"tweetTplFile": "tweet.tpl.txt"
}
```

Leave `host` and `apiVersion` unchanged or you know what you are doing.

`searchParams` is the string which you can generate [here](https://api.stackexchange.com/docs/search)

`fromDate` not listed in the config.json file but this parameter will be set automatically from goverflow.

`twitterConfigFile` path to your twitter config file where all the API keys are stored.

`tweetTplFile` path to your twitter template file. All useful variable are already in the template. There are more. Have a
look in `seapi/resources.go`. Make sure the template is within a tweet length. URLs are automatically converted to a 
t.co URLs which are at the moment 20 character long.

### config.twitter.json

```json
{
	"consumerKey": "12lljBydHmwOoObZuvRUfh9AP",
	"consumerSecret": "twH61R4fDGBP32PL0XyUQnwGbvZKH9euH3en0TRDJpRMNt6FOT",
	"accessToken": "2255586299-AgR4i7hE3D0SVMZkw3YQVl54XZ7I8g238LCJTYv",
	"accessTokenSecret": "Fz3O1HJEttYHwxSs8PnoN84TJKehudvXn7iLw31rjIVc4"
}
```

Your personal tokens can be generated at [https://dev.twitter.com/](https://dev.twitter.com/)

**ProTip**: You have to provide your mobile phone number for write access to the API. In some countries (like in Australia)
it is not possible to enter your phone number in the web interface and get a verification. So turn on your iOS, BB or Android
device, add your new/existing twitter account and also add there your mobile phone number. You will receive a SMS
verification code and your done. Switch to the dev.twitter.com website and generate your new tokens.
Enter those tokens into this file.

Of course those API keys above are my original ones ;-)

### tweet.tpl.txt

```
{{.Title}}
{{.Link}}
#golang
```

You can replace the hash tag with any other wordings. Stick to the max length of a tweet like described in section `tweetTplFile`.

For advanced templating please a look here: [http://golang.org/pkg/text/template/](http://golang.org/pkg/text/template/)

## Build

Using [https://github.com/laher/goxc](https://github.com/laher/goxc)

Setup go/src for darwin, linux and windows [http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1](http://dave.cheney.net/2013/07/09/an-introduction-to-cross-compilation-with-go-1-1)

Run `make build`. If you are interested in pre-compiled binaries, ping me.

# Contributing

If you see something say something or better send me a PR or open an issue :-)

As my Golang coding skills are below Junior level I highly appreciate short comments, fixes, etc.

# License

General Public License

[http://www.gnu.org/copyleft/gpl.html](http://www.gnu.org/copyleft/gpl.html)

# Author

[Cyrill Schumacher](https://github.com/SchumacherFM) - [My pgp public key](http://www.schumacher.fm/cyrill.asc)

[@SchumacherFM](https://twitter.com/SchumacherFM)

Made in Sydney, Australia :-)

If you consider a donation please contribute to: [http://www.seashepherd.org/](http://www.seashepherd.org/)
