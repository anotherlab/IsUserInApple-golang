# IsUserInApple - Go version
## Overly long overview
I'm the Team Agent for our company's Apple App Store Development account. Our company is large enough that I do not have any interactions with most of the people with Apple accounts. After running the account for a few years, I needed a way to remove employees who had left the company from the App Store Development account account.

Unfortunately Apple does not provide any means of doing a federated login. When an employee leaves our company, their email account is deactivated. But they can still login into AppConnect.com. For security reasons, it makes sense to remove their access to that account when they leave the company.

When an employee leaves the company, our HR department contacts our IT department as part of the off-boarding process. I'm not in IT or HR, I don't have any exposure to when an employee leaves the company.

But I can make it easier for our IT people by giving them a way to see if an employee needs to be removed.  They can run "IsUserInApple some.email@company.com" and it will tell them if that employee has in our Apple App Store Development account (or not).

## How does it work?
Apple has a REST API for App Store Connect. You call the correct endpoint and pass in a <a href="https://developer.apple.com/documentation/appstoreconnectapi/generating_tokens_for_api_requests" target="_blank">signed JWT</a> and Apple will send back a JSON result set. For this tool, we use the <a href="https://developer.apple.com/documentation/appstoreconnectapi/list_users" target="_blank">List Users web service endpoint</a> to get the list of users in the account.

## The fun part
By default, Apple limits the number of objects returned via the App Connect API to 100 objects. You can increase that to 200 by appending "&limit=200" to the API call. What they don't appear to document anywhere is that there is a simple way of getting all of the records. In the JSON result set returned by the API, there is a "links" object.  It will a "self" field that contains the URL that was used to make the call. It can have an optional "next" field will contain a URL that will return the next set of objects. When you call the API, you will need to check the "next" field and call that URL until you no longer receive another "next" field in the JSON result set.

## About this version of IsUserInApple
This is a command line app written in <a href="https://golang.org/" target="_blank">Go</a>. It was written and tested on Windows, but it should run on MacOS and Linux. There are also <a href="https://github.com/anotherlab/IsUserinApple-dotnet" target="_blank">C#</a> and <a href="https://github.com/anotherlab/IsUserInApple-python" target="_blank">Python</a> versions of the code. It requires a JSON file named IsUserinApple.json located in the same folder as the executable. This JSON file should look like this:

    {
        "PrivateKeyFile": "path/to.your/privatekey.p8",
        "KeyID": "ABCDEF1234",
        "IssuerID": "d88b7c23-4c26-48fb-9d62-5649f27a25a2"
    }

The values are

| Field          | Value                                    |
|----------------|------------------------------------------|
| PrivateKeyFile | The full path and filename to the private key file |
| KeyID | Your private key ID from App Store Connect |
| IssuerID | Your issuer ID from the API Keys page in App Store Connect |

You'll need to create an API key to sign the the JWT and authorize the API requests. To create a new key, follow these steps:

1. Log in to <a href="https://appstoreconnect.apple.com/" target="_blank">App Store Connect</a>.
2. Select Users and Access, and then select the Keys tab.  Copy the Issuer ID, you'll need it later.
3. Click the "+" button to add a new key. Fill in the Name field and set the Access to "Admin".
4. Download the private key file. You'll want to make a secure copy of the private key file, this will be the only time that you can download it.

After you create a new API key, the KeyID value will be in the column labeled "KEY ID" for the key that you just created.

## Referenced packages
This app uses the <a href="https://github.com/golang-jwt/jwt" target="_blank">golang-jwt/jwt library</a> to create and sign the JWT. The reference is in <a href="https://github.com/anotherlab/IsUserInApple-golang/blob/main/go.mod" target="_blank">go.mod</a> and should be installed when the app is built. If it doesn't get installed, you can install it with

`go get github.com/golang-jwt/jwt`

## How to build it
Assuming that you have Go installed, you can build it from the command line with

`go build .`

For the smallest executable size, build it with the following options

`go build -ldflags="-s -w" .`

Then compress the exe file with <a href="https://upx.github.io/" target="_blank">upx</a>.

## How to run it
You would run IsUserInApple with a required command line parameter, the email to look for. If you leave it out, you will get an error message telling you what is needed. Remember that the IsUserInApple.json file needs to be in the same folder as the executable. A second optional command line parameter will let you specify the name of the config file.

If you are running it with the go command, you would run it as

`go run . -username "some.email@company.com"`

or

`go run . -username "some.email@company.com" -config "IsUserInApple.json"`

Otherwise you would run it as

`IsUserInApple -username "some.email@company.com"`

You can also specify a list of user names with the `-userlist` parameter

`go run . -userlist somefile.txt`

If you specifify both `-userlist` and `-username`, the value passed to `-userlist` will be used.

You can also pass in "-h" or "--help" as command line parameters
