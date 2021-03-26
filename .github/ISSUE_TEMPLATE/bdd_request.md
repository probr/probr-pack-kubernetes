---
name: BDD Scenario request
about: Suggest new BDD scenario specification
title: "[Scenario]"
labels: feature request, scenario
assignees: ''
---

**Please write a rough BDD specification**
_[Replace sample below with actual BDD specification. See Gherkin syntax for reference: https://cucumber.io/docs/gherkin/reference]_
```
Feature: Guess the word

  # The first example has two steps
  Scenario: Maker starts a game
    When the Maker starts a game
    Then the Maker waits for a Breaker to join
```

**Who is the SME validating this scenario?**
_[Enter name of SME(s) here]_

**Please describe a proposed implementation for the above scenario**
_[Replace sample content with actual steps]_
| Scenario Step | Implementation Plan |
|---|---|
|When the Maker starts a game|Call api endpoint and start game service|
|Then the Maker waits for a Breaker to join|Call api end point and check status is Waiting|