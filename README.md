# SimpleServer

<img src="https://travis-ci.com/dtp263/SimpleServer.svg?branch=master">


This is a simple GoLang server with 3 dummy endpoints.

* GET / -> return "Default Get"
* POST /json -> takes a json blob request, parses it into a go struct then returns valid json
* POST /donation_jsonapi -> takes a blog request with Donation data in the json:api format. It parses the Request into a go struct with type Donation, then converts it back into a string and returns the json blob in the Respone