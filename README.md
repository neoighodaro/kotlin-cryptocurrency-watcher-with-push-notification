# Crypto Alat
An app that shows notifications when crypto value exceeds your saved limit.


## Getting Started

Clone the repository. The repository contains a `backend` folder for the server. You can open this project directly in Android studio. Replace the `google-services.json` file with the one from your 
Firebase dashboard. Replace the key holders in the app with the keys from your Pusher Beams.

Open the backend folder install the following dependencies:

```
npm install body-parser
npm install @pusher/push-notifications-server express --save
npm install axios
```

and run this this command to get your server up: 

```
node index.js
```

### Prerequisites

You need the following installed:

* [Android Studio](https://developer.android.com/studio/index)
* [Node](http://nodejs.org)


## Built With

* [Kotlin](https://kotlinlang.org/) - Used to build the Android client
* [Pusher](https://pusher.com/) - APIs to enable devs building realtime features
* [Node](http://nodejs.org) - Used to build the server

## Acknowledgments
