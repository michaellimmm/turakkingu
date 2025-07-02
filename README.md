# Turakkingu

simple tracking

## Installation

### Prerequisites

- Go 1.23 or higher
- [air](https://github.com/air-verse/air)
- [caddy](https://caddyserver.com/docs/install)
- install bun

### How to run

#### Backend

- install dependency `go mod tidy`

```bash
air
```

or

```bash
go run main.go
```

#### Sample of client landing page

- install dependecy

```bash
cd example/car-dealer
bun install
bun dev
```

#### Sample of client thank you page (cross-domain)

- install dependecy

```bash
cd example/car-form
bun install
bun dev
```

### Enable Multi-domain

#### add domain to `etc/hosts`

Map your custom domains to 127.0.0.1:

```bash
sudo nano /etc/hosts
```

Add:

```plaintext
127.0.0.1 tracker.local
127.0.0.1 cardealer.local
127.0.0.1 cardealerform.local
```

#### start caddy

```bash
caddy run --config ./Caddyfile --adapter caddyfile
```

## Javascript Code Snipped

```html
<script>
  var script = document.createElement("script");
  script.defer = true;
  script.dataset.trackingId = "685cbfb8085b1462689b2447";
  script.src = "http://localhost:8080/static/conversion.js";
  document.getElementsByTagName("head")[0].appendChild(script);
</script>
```
