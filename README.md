# Crypto Alat
An app that shows notifications when crypto value exceeds your saved limit.

[Part One](https://pusher.com/tutorials/cryptocurrency-kotlin-go-part-1) | [Part Two](https://pusher.com/tutorials/cryptocurrency-kotlin-go-part-2)


## Getting Started

Clone the repository. The repository contains a `backend` folder for the backend server and an `android-app` directory. You can open the Android project directly in Android studio. Replace the `google-services.json` file with the one from your
Firebase dashboard. Replace the key holders in the app with the keys from your Pusher Beams.

Open the `backend` folder install the following dependencies:

```
$ go get github.com/labstack/echo
$ go get github.com/labstack/echo/middleware
$ go get github.com/pusher/push-notifications-go
```

and run this this command to get your server up:

```
$ go run main.go
```

Replace the pusher keys in `./notification/push.go` before starting the server.

### Prerequisites

You need the following installed:

- [Android Studio](https://developer.android.com/studio/index) installed on your machine (v3.x or later). Download here.
- Go version 1.10.2 or later [installed](https://golang.org/doc/install#install).
- SQLite installed on your machine.
- Basic knowledge on using the Android Studio IDE.
- Basic knowledge of Kotlin programming language. See the official docs.
- Basic knowledge of Go and the Echo framework.


## Built With
* [Kotlin](https://kotlinlang.org/) - Used to build the Android client
* [Pusher](https://pusher.com/) - APIs to enable devs building realtime features
* [Go](https://golang.org/doc/install#install) - Used to build the server
