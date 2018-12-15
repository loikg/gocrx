# Gocrx

Gocrx is a simple utility to download chrome extesion .crx file from a list of extension url.

## Usage

By default gocrx for a `extension.txt` in the directory it is executed in and download the crx in that same directory

### Option

* `--file`: path to the file containing the list of extension
* `--output`: path to the directory where crx file will be downloaded
* `--version`: version of chrome the extension will be used with (default: 70.0)

### Example of a extensio.txt:
```
https://chrome.google.com/webstore/detail/privacy-badger/pkehgijcmpdhfbdbbnkijodmdjhbjlgp
https://chrome.google.com/webstore/detail/https-everywhere/gcbommkclmclpchllfjekcdonpmejbdp
```
