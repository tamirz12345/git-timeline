# GitTimeline
A timeline for social network posts with editing capabilities based on VCS data storage.
Allowing to edit and view past versions of a post.

This code example was written to show usage of Git as a data storage for the article "Git Not Just For Your Code".

## Storages
This are the current implimented storages for the server.
* GitStorage that running with local Git Repostiory using GoGit lib to store posts content and edit versions of it
* MysqlStorage that running with Mysql DB to store metadata of the posts.

The server works with interfaces of TimelinePostStorage and PostMetadataStorage and the two storages mentioned above are the implemented onces but it is possible to impliment them with any DB or VCS like product that fits the interface.

Git Storage works with interface for RepositoryCommunicator and the current implimented way is with LocalGitRepositoryCommunicator which is for working with repository on local directory.

## Local Run
Steps to run with local directory for git and mysql container running locally for mysql

1. Create Folder for the repository
```
mkdir /tmp/gitPlayground
```
2. Run mysql container
```
docker run --name my_mysql_container -p 3306:3306 -e MYSQL_ROOT_PASSWORD=pass -e MYSQL_DATABASE=git_timeline -d mysql:latest
```
3. Define environment variables .env example
```
DB_HOST=127.0.0.1:3306
DB_PASS=pass
DB_USER=root
REPOSITORIES_ROOT_PATH=/tmp/gitPlayground
SERVER_PORT=8080
```
4. Run with favorite idea

