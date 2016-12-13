# Datadog Event Resource

Fetch or emit events to Datadog.


## Source Configuration

* `application_key`: *Required.* The application key to use when accessing Datadog.

* `api_key`: *Required.* The api key to use when accessing Datadog.


## Behavior

### `check`: Listen for events in the event stream.

Detects new events that have been published to your Datadog event stream.


### `in`: Fetch an event

Places the following files in the destination:

* `event.json`: The event fetched, as JSON. Example:

    ```json
    {
        "date_happened": 1346449298,
        "handle": null,
        "id": 1378859526682864843,
        "priority": "normal",
        "related_event_id": null,
        "tags": [
            "environment:test"
        ],
        "text": null,
        "title": "Did you hear the news today?",
        "url": "https://app.datadoghq.com/event/jump_to?event_id=1378859526682864843"
    }
    ```

* `version`: The ID of the event fetched.

#### Parameters

*None.*


### `out`: Emit an event to Datadog

Emits an event based on the static configuration defined in your parameters.

#### Parameters

* `title`: *Required.* The event title. Limited to 100 characters.

* `text`: *Required.* The body of the event. Limited to 4000 characters. The text supports markdown.

* `priority`: *Optional.* The priority of the event ('normal' or 'low').

* `host`: *Optional.* Host name to associate with the event.

* `tags`: *Optional.* A list of tags to apply to the event.

* `alert_type`: *Optional.* "error", "warning", "info" or "success".

* `aggregation_key`: *Optional.* An arbitrary string to use for aggregation, max length of 100 characters. If you specify a key, all events using that key will be grouped together in the Event Stream.

* `source_type_name`: *Optional.* The type of event being posted. Options: nagios, hudson, jenkins, user, my apps, feed, chef, puppet, git, bitbucket, fabric, capistrano


## Example Configuration

### Resource

```yaml
- name: datadog
  type: datadog
  source:
    api_key: API-KEY
    application_key: APPLICATION-KEY
```

### Plan

```yaml
- get: datadog
```

```yaml
- put: datadog
  params:
    title: Did you hear the news today?
    text: Oh boy!
    priority: normal
    tags:
    - environment:test
    alert_type: info
```
