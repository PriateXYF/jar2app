package main

import (
	"embed"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"text/template"
)

const version string = "v1.0.5"

//go:embed assets/Info.plist
var plist embed.FS

//go:embed assets/universalJavaApplicationStub
var stub []byte

func copyFile(src, dst string) error {
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()
	dstFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dstFile.Close()
	_, err = io.Copy(dstFile, srcFile)
	return err
}

type Jar struct {
	AppName    string
	Info       string
	JarPath    string
	ClassPath  string
	IconPath   string
	IconName   string
	MainClass  string
	JavaHome   string
	Identifier string
	Version    string
	Copyright  string
}

func (jar *Jar) getJavaHome() {
	javaHome := os.Getenv("JAVA_HOME")
	if javaHome == "" {
		log.Println("⚠️", "get env $JAVA_HOME failed, auto checking...")
		cmd := exec.Command("bash", "-c", "java -XshowSettings:properties -version 2>&1 > /dev/null | grep 'java.home'")
		output, err := cmd.CombinedOutput()
		if err != nil {
			log.Println("❗️", "get env $JAVA_HOME failed, exit.")
			os.Exit(1)
		}
		javaHomeStr := strings.TrimSpace(string(output))
		if strings.HasPrefix(javaHomeStr, "java.home =") {
			javaHome = strings.TrimSpace(strings.TrimPrefix(javaHomeStr, "java.home ="))
		}
	}
	javaHomeInfo, err := os.Stat(javaHome)
	if javaHome == "" || err != nil || !javaHomeInfo.IsDir() {
		log.Println("❗️", "get env $JAVA_HOME failed, exit.")
		os.Exit(1)
	}
	jar.JavaHome = javaHome
	log.Println("ℹ️", "java home  : ", jar.JavaHome)
}

func (jar *Jar) getJarPath() {
	absJarPath, err := filepath.Abs(jar.JarPath)
	_, err2 := os.Stat(absJarPath)
	if err != nil || os.IsNotExist(err2) || filepath.Ext(absJarPath) != ".jar" {
		log.Println("❗️", "get jar file path failed, exit.")
		os.Exit(1)
	}
	jar.JarPath = absJarPath
	jar.ClassPath = fmt.Sprintf("Contents/Java/%s", filepath.Base(jar.JarPath))
	log.Println("ℹ️", "jar path   : ", jar.JarPath)
}

func (jar *Jar) getMainClass() {
	cmd := exec.Command("bash", "-c", "unzip -p "+jar.JarPath+" META-INF/MANIFEST.MF | grep -i '^Main-Class:' | awk '{print $2}'")
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Println("❗️", "get env $JAVA_HOME failed, exit.")
		os.Exit(1)
	}
	mainClass := strings.TrimSpace(string(output))
	if mainClass == "" {
		log.Println("❗️", "get main class failed, exit.")
		os.Exit(1)
	}
	jar.MainClass = mainClass
	log.Println("ℹ️", "main class : ", jar.MainClass)
}

func (jar *Jar) getAppName() {
	appName := jar.AppName
	if appName == "" {
		baseName := filepath.Base(jar.JarPath)
		extName := filepath.Ext(jar.JarPath)
		appName = strings.TrimSuffix(baseName, extName)
	}
	jar.AppName = appName
	log.Println("ℹ️", "app name   : ", jar.AppName)
}

func (jar *Jar) getIconPath() {
	absIconPath, err := filepath.Abs(jar.IconPath)
	_, err2 := os.Stat(absIconPath)
	if err != nil || os.IsNotExist(err2) || filepath.Ext(absIconPath) != ".icns" {
		log.Println("⚠️", "icon path  : ", "get icon file path failed, use default icon.")
		jar.IconPath = "default"
	} else {
		jar.IconPath = absIconPath
		baseName := filepath.Base(jar.IconPath)
		extName := filepath.Ext(jar.IconPath)
		jar.IconName = strings.TrimSuffix(baseName, extName)
		log.Println("ℹ️", "icon path  : ", jar.IconPath)
	}
}

func (jar *Jar) checkAndParse() {
	jar.getJavaHome()
	jar.getJarPath()
	jar.getMainClass()
	jar.getAppName()
	jar.getIconPath()
}

func (jar *Jar) createFolder() {
	paths := []string{
		fmt.Sprintf("%s.app/Contents/MacOS", jar.AppName),
		fmt.Sprintf("%s.app/Contents/Resources", jar.AppName),
		fmt.Sprintf("%s.app/Contents/Java", jar.AppName),
	}
	for _, p := range paths {
		err := os.MkdirAll(p, 0755)
		if err != nil {
			log.Println("❗️", "create init folder failed, exit.")
			os.Exit(1)
		}
	}
}

func (jar *Jar) outputPlist() {
	plistPath := fmt.Sprintf("%s.app/Contents/Info.plist", jar.AppName)
	plistFile, err := os.Create(plistPath)
	if err != nil {
		log.Println("❗️", "create Info.plist failed, exit.")
		os.Exit(1)
	}
	defer plistFile.Close()
	tmpl, err := template.ParseFS(plist, "assets/Info.plist")
	if err != nil {
		log.Println("❗️", "parse Info.plist failed, exit.")
		os.Exit(1)
	}
	err = tmpl.Execute(plistFile, jar)
	if err != nil {
		log.Println("❗️", "template Info.plist failed, exit.")
		os.Exit(1)
	}
}

func (jar *Jar) outputStub() {
	stubPath := fmt.Sprintf("%s.app/Contents/MacOS/universalJavaApplicationStub", jar.AppName)
	if err := os.WriteFile(stubPath, stub, 0770); err != nil {
		log.Println("❗️", "output stub file failed, exit.")
		os.Exit(1)
	}
}

func (jar *Jar) copyJar() {
	jarPath := fmt.Sprintf("%s.app/Contents/Java/%s", jar.AppName, filepath.Base(jar.JarPath))
	err := copyFile(jar.JarPath, jarPath)
	if err != nil {
		log.Println("❗️", "copy jar file failed, exit.")
		os.Exit(1)
	}
}

func (jar *Jar) copyIcon() {
	if jar.IconPath == "default" {
		return
	}
	iconPath := fmt.Sprintf("%s.app/Contents/Resources/%s", jar.AppName, filepath.Base(jar.IconPath))
	err := copyFile(jar.IconPath, iconPath)
	if err != nil {
		log.Println("❗️", "copy jar file failed, exit.")
		os.Exit(1)
	}
}

func (jar *Jar) generateFiles() {
	jar.createFolder()
	jar.outputPlist()
	jar.outputStub()
	jar.copyJar()
	jar.copyIcon()
}

func (jar *Jar) Build() {
	log.Println("📦", "jar2app", version)
	jar.checkAndParse()
	jar.generateFiles()
	log.Println("🎉", fmt.Sprintf("build app %s.app successful!", jar.AppName))
}

func parseFlag() (*Jar, error) {
	var jarPathFlag, iconPathFlag, appNameFlag, identifierFlag, copyrightFlag, verFlag string
	flag.StringVar(&jarPathFlag, "jar", "", ".jar file path")
	flag.StringVar(&iconPathFlag, "icon", "", ".icns icon file path")
	flag.StringVar(&iconPathFlag, "info", "Made by virts.", "app info")
	flag.StringVar(&appNameFlag, "name", "", "app name")
	flag.StringVar(&identifierFlag, "id", "app.virts", "app identifier")
	flag.StringVar(&copyrightFlag, "copyright", "Copyright 2025 virts", "app copyright")
	flag.StringVar(&verFlag, "v", "1.0.0", "app version")
	flag.Parse()
	return &Jar{
		JarPath:    jarPathFlag,
		IconPath:   iconPathFlag,
		AppName:    appNameFlag,
		Identifier: identifierFlag,
		Copyright:  copyrightFlag,
		Version:    verFlag,
	}, nil
}

func main() {
	log.SetOutput(os.Stdout)
	jar, err := parseFlag()
	if err != nil {
		log.Println("❗️", "parse flag error :", err.Error())
		os.Exit(1)
	}
	jar.Build()
}
