# go-core

## Creating new modules

1. Create a folder on root directory of this repository with the name of your new module
2. Enter the newly created directory and execute: 
 
```sh
$ go mod init github.com/facily-tech/go-core/MODULE-NAME
```

3. Put your module code inside the subdirectory you just created
4. Create a README.md describing how to use the module

**Tags should be created only in the main branch!**

```sh
$ git tag -a MODULE-NAME/v0.1.0
$ git push --tags
```

## How to import modules

```sh
go get github.com/facily-tech/go-core/MODULE-NAME
```
