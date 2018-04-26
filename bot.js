const Discord = require('discord.js');
const client = new Discord.Client();
const auth = require('./auth.json');
const config = require('./config.json');
var request = require('request');

var commandRegex = /![a-zA-Z]+ [a-zA-Z]+/g
var stockRegex = /^[a-zA-Z]{1,4}$/g
var coinRegex = /^[a-zA-Z]{1,4}$/g

client.on('ready', () => {
  console.log('I am ready!');
});

client.on('message', message => {
  if (message.content.match(commandRegex)) {
    var contents = message.content.split(' ');
    var output;
    if (contents[0].toLowerCase() === "!stock") {
      if (contents[1].match(stockRegex)) {
        var tickerUrl = config.stock_api_url + contents[1] + "/batch?types=quote";
        request.get(tickerUrl, function (err, res, body) {
          try {
            var info = JSON.parse(body).quote;
          } catch (e) {
            message.channel.send(body);
            return
          }
          output = {
            "Symbol": info.symbol,
            "Company Name": info.companyName,
            "Current": info.latestPrice,
            "High": info.high,
            "Low": info.low,
            "Open": info.open,
            "Close": info.close,
            "Percent Change (1 Day)": info.changePercent * 100 + "%",
            "Volume": info.latestVolume
          }
          message.channel.send(formatOutput(JSON.stringify(output)));
        });
      } else {
        message.channel.send('Invalid Symbol Fuck You');
      }
    } else if(contents[0].toLowerCase() === "!coin") {
      if(contents[1].match(coinRegex)) {
        var coinUrl = config.coin_api_url + contents[1].toUpperCase() + "&tsyms=USD";
        request.get(coinUrl, function (err, res, body) {
          var info = JSON.parse(body);
          if (info.Message != undefined) {
            message.channel.send(info.Message);
          } else {
            output = {
              "Symbol": contents[1].toUpperCase(),
              "Price (USD)": info.USD
            }
            message.channel.send(formatOutput(JSON.stringify(output)));
          }
        });
      } else {
        message.channel.send('Invalid Coin Symbol Buckfoy')
      }
    } else {
      message.channel.send('Invalid Command Fuck You');
    }
  }
});

var formatOutput = function (message) {
  var output = message.split(',').join('\n');
  output = output.substring(1, output.length - 1);
  return ("```" + output + "```");
}

client.login(auth.bot_token);
