# Turakkingu

simple tracking

## Installation

### Prerequisites

- Go 1.23 or higher
- [air](https://github.com/air-verse/air)
- install dependency `go mod tidy`

### How to run

```bash
air
```

or

```bash
go run main.go
```

## Javascript Code Snipped

```html
<script>
  var script = document.createElement("script");
  script.defer = true;
  script.dataset.trackingId = "685cbfb8085b1462689b2447"
  script.src = "http://localhost:8080/static/conversion.js";
  document.getElementsByTagName("head")[0].appendChild(script);
</script>
```
