# GitHubby
GitHubby is the hubby your github has always needed. It comes in your life, and brings a lot of greens with it. Are you someone who could do with some greens?

### Who are you? Are you one of the people whose Github profile looks like this?
*<kbd>
  ![Who you are](https://i.imgur.com/rsFMYog.png)
</kbd>

This is so sad. You are so sad. You're a loser and you've always wanted not to be you.

### Now, do you wish your Github profile would look like this? 
*<kbd>
![Who you are not](https://i.imgur.com/VvDLzpy.png)
</kbd>

Worry not, it's time to bring Githubby into your life and add some spice to your Github.

## What is Githubby?
GitHubby is a git repository, duh. Within the git repository is a piece of code. Double duh. That piece of code, if set up right, kind of is a better version of you in a sense because it automatically contirbutes to it's own repository (make edits, git add files, commits, and pushes). It makes your Github look greener, something you could never manage.

It is also pretty smart. In order to make it look more natural, it randomly selects the number commits to make and decides to randomly wait from anywhere to 1s to 1hr before making the commit. And, most importantly, the commit messages are pulled from http://whatthecommit.com/index.txt (which is a great project in itself).

## Getting Started
### Installing
* Request to be a contributer on this project.
* Install Golang
* Clone this package 
* Install go dependecies by running `make install` from within this repository folder
* Build/compile the package using `make build`.

### Running
From the project directory do `make run`. The program will automatically make a change to the changeme.txt file, add the file to git, fetch a new commit message and run git commmit, and then push the changes upstream. If you have permisions to push to this repository, everything should be fine. You'll see contirbutions from yourself in the commits section of this project.

### Automate
Ideally, you shouldn't have to run this manually. You can set up a cron job in your system to invoke the program automatically. Your cron job can call the command `<path/to/repo>/githubby.out` at a time of your choice. Shhh. Now, you're pretty active on Github. You can go find your next passion.
