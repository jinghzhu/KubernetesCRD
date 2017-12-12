[![Build Status](https://travis-ci.org/jinghzhu/GoUtils.svg?branch=master)](https://travis-ci.org/jinghzhu/GoUtils)

# Introduction

This is just a personal repo to write some Go helper packages.


## Collaborating

One of the most effective ways to collaborate on GitHub is by using a forking/branching model as described in the [GitHub Guides](https://guides.github.com/).


## Doing Work

* Each time you begin doing work on a new issue, check out the master branch by doing `git checkout master`. You will only be able to do this if you don't have any changes in your local codebase.
* Pull in the latest changes from upstream's master branch - `git pull upstream master`
* Create a new [branch](https://guides.github.com/introduction/flow/), named something relevant to the issue being worked on - `git checkout -b {{branch-name}}`, replacing `{{branch-name}}` with the name of your branch.
* Push your new branch to your origin remote - `git push -u origin {{branch-name}}`
* Add your commits and push to that branch - `git push origin {{branch-name}}`
* Issue a Pull Request in to the upstream repository when the work is done. Make sure the Pull Request comment includes a [keyword for closing issues](https://help.github.com/articles/closing-issues-via-commit-messages/) for closing the issue the work is for - `Resolves #42` (with `42` being the issue number)
* Once the Pull Request is merged, delete the local and remote branch you worked on - `git branch -d {{branch-name}}` for local, `git push origin :{{branch-name}}` for remote. **Important: Never Reuse A Branch After It Has Been Merged**



## sitespeedAnalyzer

[Sitespeed.io](https://www.sitespeed.io/) is an open source tool that helps you analyze your website speed and performance based on performance best practices and timing metrics.

After running Sitespeed.io, there would generate a json file which records the socre of each performance rule.

`sitespeedAnalyzer` aims to analyse these results and gets some helpful statics information.



