# go
first app in go

To install package:
go get 'ur package'

Executing a go app:
go run main.go


Configuration

We need a couple small configuration files to tell Heroku how it should run our application. The first one is the Procfile, which allows us to define which processes should be run for our application. By default, Go will name the executable after the containing directory of your main package. For instance, if my web application lived in GOPATH/github.com/codegangsta/bwag/deployment, my Procfile will look like this:

web: deployment
Specifically to run Go applications, we need to also specify a .godir file to tell Heroku which dir is in fact our package directory.

deployment


Deployment

Once all these things in place, Heroku makes it easy to deploy.

Initialize the project as a Git repository:

git init
git add -A
git commit -m "Initial Commit"
Create your Heroku application (specifying the Go buildpack):

heroku create -b https://github.com/kr/heroku-buildpack-go.git
Push it to Heroku and watch your application be deployed!

git push heroku master
View your application in your browser:

heroku open
