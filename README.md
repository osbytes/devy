# dev-hub-bot

Discord bot for dev hub discord channel

## Run

Copy `.env.sample` to `.env` and add secrets

```sh
cp .env.sample .env
```

Install [godotenv](https://github.com/joho/godotenv)

```sh
go install github.com/joho/godotenv/cmd/godotenv@latest
```

Run the following command to run the application.

```sh
godotenv -f .env go run cmd/bot/main.go
```

## Stack

- Go v1.17
- Discord

## How to Contribute

- Fork the project
- Push changes
- Create a PR and add reveiwers

## How to Test on a Test Bot

## Testing

### Test Naming

```go
// function to test
func GetContributionsByUsername() {
    // logic
}

// notice the naming for the main test for GetContributionsByUsername
func TestGetContributionsByUsername(t *testing.T) {
    // test
}

// test modifiers have an underscore followed by what you are testing for
func TestGetContributionsByUsername_MultiYear(t *testing.T) {
    // test
}


func TestGetContributionsByUsername_DatesZeroValue(t *testing.T) {
    // test
}

```

### Mocking

Run the command below to mock all of your interfaces

```sh
mockery --all --inpackage
```

If you need to monkey patch or create pointer functions follow this convention

```go
// keep the pointer functions at the top of the file
var (
    doSomethingF = doSomething
)

// make sure in the implementation you call the pointer
func GetContributionsByUsername() {
    something, err := doSomethingF(args)
}

// now you can mock that function
func TestGetContributionsByUsername(t *testing.T) {
    doSomethingF = func(args) (something, error) { return something, err }
}

```

## Ideas

- [ ] Forces you to change nick name to real name
- [ ] Displays github data
- [ ] Gives us newest fireship videos
- [ ] Scrapes for new changes to certain lang's
- [ ] Coding challenges and scoreboards
- [ ] Maybe something with leetcode
- [ ] New Job openings for those looking for a new job

## Known Issues

Error displaying in console

```sh
YYYY/MM/DD hh:mm:ss error closing resp body
```

https://github.com/bwmarrin/discordgo/issues/1028
https://github.com/golang/go/issues/49366
