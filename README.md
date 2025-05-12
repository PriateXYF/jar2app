# jar2app

Pack Jar files elegantly into MacOS Apps, supporting icon and name settings.

## Installation

go install :

```bash
go install -v github.com/PriateXYF/jar2app@latest
```

Or you can use brew :

```bash
brew install PriateXYF/tap/jar2app
```

## How it works

* jar2app can only run on MacOS, and does not support Windows systems.
* Before using jar2app, please install the corresponding version of JDK / JRE on your Mac. Considering the file size, jar2app will not package the runtime into the App.

## Usage

```bash
Usage of jar2app:
  -copyright string
    	app copyright (default "Copyright 2025 virts")
  -icon string
    	.icns icon file path
  -id string
    	app identifier (default "app.virts")
  -info string
    	app info (default "Made by virts.")
  -jar string
    	.jar file path
  -name string
    	app name
  -v string
    	app version (default "1.0.0")
```

* Example

```bash
# Base usage
jar2app --jar file.jar --icon icon.icns --name MyApp
```

```bash
# Use a specified $JAVA_HOME
JAVA_HOME="/Library/Java/JavaVirtualMachines/jdk-1.8.jdk/Contents/Home/" jar2app --jar file.jar --icon icon.icns --name MyApp
```

```bash
# Full parameters
export JAVA_HOME=/Library/Java/JavaVirtualMachines/jdk-1.8.jdk/Contents/Home/
jar2app --jar file.jar --icon icon.icns --name MyApp --info "My App" --v "1.0.0" --id "app.virts" --copyright "Copyright 2025 virts"
```

## Reference

* [universalJavaApplicationStub](https://github.com/tofi86/universalJavaApplicationStub)