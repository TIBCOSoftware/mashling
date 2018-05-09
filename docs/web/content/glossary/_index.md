---
title: Glossary
weight: 7000

pre: "<i class=\"fa fa-book\" aria-hidden=\"true\"></i> "
---

Mashling terms and constructs, defined here, all in one place, in a logical order vs. alphabetical.

#### Recipe
One or more trigger and handler pair with mock configurations where necessary. It is use case driven and represents a gateway once fully configured.

#### Trigger
Trigger point in a recipe execution. For example, a trigger can be a subscriber on an MQTT topic or Kafka topic. The trigger is responsible for accepting the incoming event and invoking one or more defined actions (flows). Triggers are not coupled to flows, that is, a flow can exist without a trigger.

#### Handler
A handler is used to map triggers to activities that perform gateway functions.

#### Flow 
A flow is an implementation of an action and is the primary tool to implement business logic in Mashery. A flow can consist of a number of different constructs:

* One or more activities that implement specific logic (for example write to a database, invoke a REST endpoint, etc).
* Each activity is connected via a link.
* Links can contain conditional logic to alter the path of a flow

Flows, as previously stated in the triggers section, can exist without a trigger. Thus, flows operate very similar to functions, that is, a single flow can define its own input & output parameters. Thus, enabling a flow to be reused regardless of the trigger entrypoint. All logic in the flow only operates against the following data:

* Flow input parameters
* Environment variables
* Application properties
* The output data from activities referenced in the flow

The flow cannot access trigger data directly, trigger input and output data must be mapped into the flows input and output parameters.

#### Event link
A concrete binding condition between trigger and handler.

#### Event
Inbound object into a trigger.

#### Kafka event
A record from kafka topic.

#### HTTP event
An HTTP request.

#### Service
Destination of a flow’s outbound object.

#### Gateway
A deployable unit that contains complete configuration of runtime entities.

#### Gateway model 
Description of all flows in a gateway.

#### Gateway model schema
JSON schema for a gateway model.

#### Mashling.json
File containing a gateway model.
