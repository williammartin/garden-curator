# garden-curator
Garden Curator is a way to specify and grow your Garden containers in a declarative manner. By specifying a blueprint of your application, and using the `garden-curator grow` command, you can very easily and reproducibly bring up each component in a Garden container.

## Contributing
Curator is open to any contributions!

### Getting started
`git clone https://github.com/williammartin/garden-curator.git`

### Development guidelines
Curator is written in a BDD manner using the Ginkgo testing library for Go. It would be appreciated if any code contributions were well tested.

### Running the tests
The `cqt (curator question time)` suite requires a running Garden server to create containers against. The easiest way to stand up Garden is on bosh-lite, using the manifest from https://github.com/cloudfoundry/garden-runc-release. Once you have this deployed:

`ginkgo -r` and voila!

### Building Curator
`go build` and voila!

