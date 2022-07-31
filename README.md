<h1 style="border-bottom: none;" align="center">Devy</h1>

<p align="center">
    <img src="./devy.svg" height="175">
</p>

<p align="center">
    <a href="https://codecov.io/gh/osbytes/devy">
        <img src="https://codecov.io/gh/osbytes/devy/branch/main/graph/badge.svg" alt="codecov" />
    </a>
    <a href="https://github.com/osbytes/devy/blob/main/LICENSE">
        <img src="https://img.shields.io/github/license/osbytes/devy.svg" alt="License" />
    </a>
</p>

A developer focused discord bot written in go

## How to Get Started

Install all go dependencies

```sh
go get ./...
```

### Run the app

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

Note: If you grab an issue that is labled TODO, please delete the TODO comment.

## How to Test on a Test Bot

## Testing

If you are adding a test please make sure to delete any of the todo comment once you push your changes

```go
// TODO Tests: GetFirstContributionYearByUsername
// labels: tests
func TestGetFirstContributionYearByUsername(t *testing.T) {

}

// TODO Tests: GetFirstContributionYearByUsername
// labels: tests
func TestGetFirstContributionYearByUsername(t *testing.T) {

}
```

### Test Naming

```go
// function to test
func (g *GithubService) GetContributionsByUsername() {
    // logic
}

// notice the naming for the main test for GetContributionsByUsername
// the struct followed but a single underscore and the receiver method name
func TestGithubService_GetContributionsByUsername(t *testing.T) {
    // test
}

// test modifiers are separated by a double underscore followed by what you are testing for
func TestGithubService_GetContributionsByUsername__MultiYear(t *testing.T) {
    // test
}

```

### Mocking

Run the command below to mock all of your interfaces

Pull mockery docker image

```sh
docker pull vektra/mockery
```

Run this in devy, replace pwd with root pwd

```sh
docker run -v "$PWD":/src -w /src vektra/mockery --all
```

For in package

```sh
docker run -v "$PWD":/src -w /src vektra/mockery --all --inpackage
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
