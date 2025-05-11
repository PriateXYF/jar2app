# jar2app

Pack Jar files elegantly into MacOS Apps, supporting icon and name settings.

## Installation

```bash
go install -v github.com/PriateXYF/jar2app@latest
```

## Usage

```bash
./jar2app -h
Usage of ./jar2app:
  -copyright string
    	app copyright (default "Copyright 2025 virts")
  -icon string
    	.icns icon file path
  -id string
    	app identifier (default "app.virts")
  -jar string
    	.jar file path
  -name string
    	app name
  -v string
    	app version (default "1.0.0")
```

* Example

```bash
jar2app --jar file.jar --icon icon.icns --name MyApp
```