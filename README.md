# ChristmasList

## to use:
1. `cd /opt && sudo git clone https://github.com/squishd/ChristmasList`
2. `cd ChristmasList && go build`
3. `./christmaslist`

this will get you to the basic running version. It uses a sqlite3 database on the backend which it will create. you do not need sqlite3 installed on the machine (but will need to download some dependencies for dealing with the database anyhow - see the go.mod file and read error output during `go build` step to see which packages you're missing. install them with `go get -u <packagename>`)

## to make a service:
1. change the christmaslist.service file
    * under \[Service\] change User to the desired username
    * under \[Service\] change ExecStart and WorkingDirectory to point to the correct directory, if you chose an alternative install location.
2. `sudo cp christmaslist.service /etc/systemd/system/`
3. `sudo systemctl daemon-reload`
4. `sudo systemctl start christmaslist.service`
5. optional: enable (automatically restart at reboot):
    * `sudo systemctl enable christmaslist.service`

## insert a user:
1. from the index page (http://localhost:8080/) select "Add Gift", and type their name into the text box.
2. (optional) - add information about a gift for their wishlist.

## update items purchased
1. from the index page (http://localhost:8080/) or the general gifts page (http://localhost:8080/giftlist) select the desired user, find the item in question and click "toggle \<item name\>"

## delete an item from a list
### NOTE: AT THIS POINT DELETE OPTIONS ARE FINAL, THERE IS NO WAY TO UNDO THIS ACTION.
1. from the index page (http://localhost:8080/) or the general gifts page (http://localhost:8080/giftlist) select the desired user, find the item in question and click "delete \<item name\>"

## delete a recipient
### NOTE: AT THIS POINT DELETE OPTIONS ARE FINAL, THERE IS NO WAY TO UNDO THIS ACTION.
1. just can't make it work? from the index page (http://localhost:8080/) or the general gifts page (http://localhost:8080/giftlist) select the desired user, and click "Delete \<user name\>"

## export your list
1. from the index page (http://localhost:8080/) right-click "export plain text" and select "save target link", then save it to the desired location.