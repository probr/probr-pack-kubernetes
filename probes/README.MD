# Probes Engineering Notes

Each "probe" is part of a Probr "event" which provides topical organization. Events are sorted within topical subdirectories within `probes/`.

## Features == Events // true

The entry point for each **event** is a behavior-driven development (BDD) _feature file_ that contains human-readable instructions for one or more probes.  Each of these feature files represents a separate Probr **event**.

The term "feature" here harkens to the original intent of BDD and _cucumber_, the technology we currently use for BDD. Cucumber and BDD logic traditionally revolves around writing feature tests, which simulates user interactions with an application.

If it serves you to simply read _feature_ instead of **event**, you would not be inaccurate for making that mental substitution when reading the Probr code, documentation, and logs. However, we have chosen to use the term "events" because we believe it more clearly demonstrates that this is part of the core Probr functionality- and not simply a feature test meant to ensure system reliability.

In the audit logs you will see the term "Events" as a top level object, which will contain each individual **event**, which will then subsequently contain every individual probe.

For more reading on how a _feature file_ is written, please review the [Cucumber documentation](https://cucumber.io/docs/gherkin/reference/)

## Probes == Scenarios // true

Probes are the heart of Probr, for obvious reasons. Each individual **probe** will attempt a behavior, and return a success or failure based on whether the result matches the expected result.

These **probes** are defined in _feature files_, and in Gherkin they would be considered _scenario outlines_.

Similar to the section above, we believe that the most accurately descriptive terminology for our use-case is to substitute in the term **probe** for "scenario outline."

## Event/Feature Logic

Because each probe is a _scenario outline_, it will have a corresponding golang function for each _step_.
 The primary logic for the steps is found in a `.go` file that is named according to the event. Comments surrounding the logic will help distinguish which function belongs to which probe.

This is not the end of the story, however. The functions within the event go files serves only to connect internal functionality with the BDD step. The complete, object-oriented, auditable logic for each probe will live in an appropriately named directory within `internal/`.
