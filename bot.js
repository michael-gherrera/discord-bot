const Discord = require('discord.js');
const client = new Discord.Client();
const auth = require('./auth.json');
const utils = require('./utils.json');
var request = require('request');

client.on('ready', () => {
  console.log('I am ready!');
});

client.on('message', message => {
  if (message.content.indexOf('!') === 0 && message.content.split(' ').length === 2) {
    var contents = message.content.split(' ');
    if (contents[0] == "!stock") {
      var tickerUrl = utils.prices_url + "function=TIME_SERIES_DAILY&symbol=" + contents[1] + "&apikey=" + auth.api_key;
      request.get(tickerUrl, function (err, res, body) {
        var jsonBody = JSON.parse(body);
        var currentDate = jsonBody["Meta Data"]["3. Last Refreshed"];
        message.channel.send(formatOutput(JSON.stringify(jsonBody["Time Series (Daily)"][currentDate])));
      });
    } else if (contents[0] == "!coin") {
      var coinUrl = utils.prices_url + "function=DIGITAL_CURRENCY_DAILY&symbol=" + contents[1] + "&market=USD&apikey=" + auth.api_key;
      request.get(coinUrl, function (err, res, body) {
        var jsonBody = JSON.parse(body);
        var currentDate = (jsonBody["Meta Data"]["6. Last Refreshed"].split(' '))[0];
        var output = formatOutput(JSON.stringify(jsonBody["Time Series (Digital Currency Daily)"][currentDate]));
        
        message.channel.send(output)
      });
    }
  }
});

var formatOutput = function (json) {
  var output = json.split(',').join('\n');
  output = output.substring(1, output.length - 1);
  return ("```" + output + "```");
}

/*
 * Current date in YYYY-MM-DD, if it's saturday or sunday, we'll need the current date to change to friday?
 */
var getCurrentDate = function () {
  var date = new Date();
//  if (type === "stock") {
//    //Sunday
//    if (date.getDay() === 0) {
//      
//    }
//  }
  return (date.toISOString().substring(0, 10));
}

client.login(auth.bot_token);
