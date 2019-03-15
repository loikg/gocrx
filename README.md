#Gocrx

Gocrx is a simple utility to download chrome extesion .crx file from a list of extension url.

##Usage

```
A tool to download chrome extension .crx files.
Can read from a file or download by extension id.

Usage:
  gocrx <file|id> [destination] [flags]

Flags:
  -c, --chrome string   Chrome version for which extension are downloaded (default "72.0")
  -h, --help            help for gocrx
  -w, --worker int      Number of parallel workers (default 4)
```


###Example of a extension.txt:
```
privacy-badger: pkehgijcmpdhfbdbbnkijodmdjhbjlgp
https-everywhere: gcbommkclmclpchllfjekcdonpmejbdp
```