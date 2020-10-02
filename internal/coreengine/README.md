# Core Engine

Internal core-engine code to manage test execution, test status and results reporting

# Engineering Notes

Each probe is a `Test` object:

```
type Test struct {
	TestDescriptor *TestDescriptor `json:"test_descriptor,omitempty"`

	Status *TestStatus `json:"status,omitempty"`

	Results *bytes.Buffer
}
```

The `TestDescriptor` above is used to identify which available handler should be use for a specific test.

*Note: In `probr-0.1.x`, all handlers are of the type `GoDogTestTuple` but this can change in the future by modifying `coreengine.handlers` to use duck-typed handler objects. Currently that duck-typing is not necessary, and all test handlers are of the type `GoDogTestTuple`.*

Handlers will contain a `Handler` and `Data`; two objects that cooperate to run the appropriate test:

```
type GoDogTestTuple struct {
	Handler TestHandlerFunc
	Data    *GodogTest
}
```

The test tuple here contains a `TestHandlerFunc` which is a generic function that takes its value from the its caller.

```
type TestHandlerFunc func(t *GodogTest) (int, *bytes.Buffer, error)
```

*Note: The use of a generic here allows for any type of handler to be specified by the logic for a particular feature. In `probr-0.1.x` the only test handler is `probes.GodogTestHandler`.*

An example usage for `GoDogTestTuple` is as follows:

```
coreengine.GoDogTestTuple{
    Handler: probes.GodogTestHandler,
    Data: &GODOGTESTOBJ,
}
```

In the above example, a properly instantiated `GodogTest` object is the required value for `Data`.

```
type GodogTest struct {
	TestDescriptor       *TestDescriptor
	TestSuiteInitializer func(*godog.TestSuiteContext)
	ScenarioInitializer  func(*godog.ScenarioContext)
	FeaturePath          *string
}
```

Here is an example where all of the above logic is brought together to add a new handler that can be utilized by tests:

```
td := coreengine.TestDescriptor{Group: coreengine.CloudDriver,
    Category: coreengine.General, Name: "account_manager"}

fp := filepath.Join("probes", "clouddriver", "events")

coreengine.AddTestHandler(td, &coreengine.GoDogTestTuple{
    Handler: probes.GodogTestHandler,
    Data: &coreengine.GodogTest{
        TestDescriptor:       &td,
        TestSuiteInitializer: TestSuiteInitialize,
        ScenarioInitializer:  ScenarioInitialize,
        FeaturePath:          &fp,
    },
})
```

After the above logic is called, a new handler will be available at `coreengine.handlers`, and can be used in the following manner:

```
g.Handler(g.Data)
```

In the above example, passing the `Data` object to the `Handler` object will allow the `GodogTest` values that were previously defined on the handler to be passed in to the `Handler` function to run the tests that are defined in the feature files specified in the `GodogTest` object.