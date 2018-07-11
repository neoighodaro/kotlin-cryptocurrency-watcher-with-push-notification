
let express = require('express');
    const axios = require('axios');
    let bodyParser = require('body-parser');
    let https = require('https');
    const PushNotifications = require('@pusher/push-notifications-server');
    
    // initialize express and pusher
    let app = express();
    let minBTC = 0, maxBTC=0;
    let minETH = 0, maxETH=0;
    let pushNotifications = new PushNotifications(require('./config.js'))
    
    // Middlewares
    app.use(bodyParser.json());
    app.use(bodyParser.urlencoded({ extended: false }));
    app.get('/fetch-values', (req, res) => {
        axios({
            method:'get',
            url:'https://min-api.cryptocompare.com/data/pricemulti?fsyms=BTC,ETH&tsyms=USD'
        })
        .then(function (response) {
        
            let currentBTC = response.data['BTC'].USD;
            let currentETH = response.data['ETH'].USD;

            var mutate = req.query.mutate;
            console.log(mutate);
            if (mutate=="true"){
                currentBTC = currentBTC + 10000;
                currentETH = currentETH + 10000;
            }
            console.log(currentBTC);
            if(maxBTC!=0 && minBTC!=0){
                if (currentBTC>maxBTC || currentBTC<minBTC){
                    pushNotifications.publish(
                        ['crypto'],{
                        fcm: {
                          notification: {
                            title: 'Bitcoin price update',
                            body: 'The new price for Bitcoin is: '+ currentBTC.toString()
                          }
                        }
                      }).then((publishResponse) => {
                        console.log('Just published:', publishResponse.publishId);
                      });
    
                }
            }
            if(maxETH!=0 && minETH!=0){
                if (currentETH>maxETH || currentETH<minETH){
                    pushNotifications.publish(
                        ['crypto'],{
                        fcm: {
                          notification: {
                            title: 'Etherum price update',
                            body: 'The new price for Etherum is: '+ currentETH.toString()
                          }
                        }
                      }).then((publishResponse) => {
                        console.log('Just published:', publishResponse.publishId);
                      });
                }
            }
            
            res.send(response.data);
    
        })
        .catch(function (error) {
            console.log(error);
        });
    });

    app.post('/btc-pref', (req, res, next) => {
        minBTC = req.body.minBTC;
        maxBTC = req.body.maxBTC;
  
        res.json({success: 200})
      });
      
    app.post('/eth-pref', (req, res, next) => {
        minETH = req.body.minETH;
        maxETH = req.body.maxETH;
  
        res.json({success: 200})
    });

    // Index
    app.get('/', (req, res) => res.json("It works!"));
    
    // Serve app
    app.listen(4000, _ => console.log('App listening on port 4000!'));
    