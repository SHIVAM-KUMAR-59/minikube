# Cobra

### What is Cobra ?
- Cobra is a library providing a simple interface to create powerful modern CLI interfaces similar to git & go tools.
- Official User Guide [here](https://github.com/spf13/cobra/blob/main/site/content/user_guide.md)
- Installation command: `go get -u github.com/spf13/cobra@latest`

---

## Key Concepts
- **Commands** represent actions, **Args** are things and **Flags** are modifiers for those actions.
- The best applications read like sentences when used, and as a result, users intuitively know how to interact with them.
- The pattern to follow is `APPNAME VERB NOUN --ADJECTIVE` or `APPNAME COMMAND ARG --FLAG`
- Example: In the following example, 'server' is a command, and 'port' is a flag:
    ```
    hugo server --port=1313
    ```
- In this command we are telling Git to clone the url bare.
    ```git
    git clone URL --bare
    ```

### Commands
- Command is the central point of the application. Each interaction that the application supports will be contained in a Command. A command can have children commands and optionally run an action.
- In the example above, `server` is the command.
- Official documentation: [here](https://pkg.go.dev/github.com/spf13/cobra#Command)

### Flags
- A flag is a way to modify the behavior of a command.
- A Cobra command can define flags that persist through to children commands and flags that are only available to that command.
- In the example above, `port` is the flag.