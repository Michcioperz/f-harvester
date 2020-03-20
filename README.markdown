# f-harvester

a program that extracts URLs to all apps available from a given f-droid repository

## how to compile

git clone it, enter the directory, `go build` and you should have a `f-harvester` binary handy

no dependencies outside of standard library, bless the cursed scope creep of go

## how to obtain apks from my phone

1. open f-droid, go to nearby, press "find people near me"

2. turn wi-fi visibility on, remember ip address shown next to it (for example 192.168.2.55), press "scan qr code"

3. meanwhile make sure you are connected to the same network on your computer as you are on your phone

4. check all the apps you want to share, go to the next screen

5. wait for it to finish processing

6. if it suggests you use nfc, press skip

7. when it shows you a qr code and a url (such as http://192.168.2.55:8888) run f-harvester with that url suffixed with `/fdroid/repo/` (for example `f-harvester http://192.168.2.55:8888/fdroid/repo/`) and you should get a list of urls printed

8. to download all those apks, you can pipe the output to `wget -i-` (for example `f-harvester http://192.168.2.55:8888/fdroid/repo/ | wget -i-`)

9. for extra awesomeness you can swap in aria2c instead of wget (it's much less often installed by default on people's computers though) and everything should go faster and with cool progress reports and such

## troubleshooting

you probably just forgot to finish the url with `/`, i always do
