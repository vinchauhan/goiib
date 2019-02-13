# goiib

This is project is inspired by the iib-maven-plugin and is used to build IIB ( IBM Integration Bus) projects. The goal is to integrate this utility on the IIB CI/CD workflow. The challenge with the Maven Plugin build is that the BAR file created are somewhat FAT bars as the Maven class loader include everything in the Application project in the classpath and create an unecessay JAR file.

# Install

go get -v github.com/vinchauhan/goiib

# Dependency

The Application/Library need to include a build.yaml file at root just like Maven uses pom.xml

# Usage

Usage:
  goiib [command]

Available Commands:

`clean       Clean the target directory just like mvn clean`

`compile     Compile IIB Application/Project and prepare a bar file`

`package     Package the application by applying specific bar overrides`

`help        Help about any command`

Flags:
  -h, --help   help for go-iib

Use "goiib [command] --help" for more information about a command.

# License

MIT

# Contributor

Vineet Chauhan
